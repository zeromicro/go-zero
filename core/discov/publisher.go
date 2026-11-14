package discov

import (
	"time"

	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/threading"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	// PubOption defines the method to customize a Publisher.
	PubOption func(client *Publisher)

	// A Publisher can be used to publish the value to an etcd cluster on the given key.
	Publisher struct {
		endpoints  []string
		key        string
		fullKey    string
		id         int64
		value      string
		lease      clientv3.LeaseID
		quit       *syncx.DoneChan
		pauseChan  chan lang.PlaceholderType
		resumeChan chan lang.PlaceholderType
	}
)

// NewPublisher returns a Publisher.
// endpoints is the hosts of the etcd cluster.
// key:value are a pair to be published.
// opts are used to customize the Publisher.
func NewPublisher(endpoints []string, key, value string, opts ...PubOption) *Publisher {
	publisher := &Publisher{
		endpoints:  endpoints,
		key:        key,
		value:      value,
		quit:       syncx.NewDoneChan(),
		pauseChan:  make(chan lang.PlaceholderType),
		resumeChan: make(chan lang.PlaceholderType),
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

// KeepAlive keeps key:value alive.
func (p *Publisher) KeepAlive() error {
	cli, err := p.doRegister()
	if err != nil {
		return err
	}

	proc.AddWrapUpListener(func() {
		p.Stop()
	})

	return p.keepAliveAsync(cli)
}

// Pause pauses the renewing of key:value.
func (p *Publisher) Pause() {
	p.pauseChan <- lang.Placeholder
}

// Resume resumes the renewing of key:value.
func (p *Publisher) Resume() {
	p.resumeChan <- lang.Placeholder
}

// Stop stops the renewing and revokes the registration.
func (p *Publisher) Stop() {
	p.quit.Close()
}

func (p *Publisher) doKeepAlive() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-p.quit.Done():
			return nil
		default:
			cli, err := p.doRegister()
			if err != nil {
				logc.Errorf(cli.Ctx(), "etcd publisher doRegister: %v", err)
				break
			}

			if err := p.keepAliveAsync(cli); err != nil {
				logc.Errorf(cli.Ctx(), "etcd publisher keepAliveAsync: %v", err)
				break
			}

			return nil
		}
	}

	return nil
}

func (p *Publisher) doRegister() (internal.EtcdClient, error) {
	cli, err := internal.GetRegistry().GetConn(p.endpoints)
	if err != nil {
		return nil, err
	}

	p.lease, err = p.register(cli)
	return cli, err
}

func (p *Publisher) keepAliveAsync(cli internal.EtcdClient) error {
	ch, err := cli.KeepAlive(cli.Ctx(), p.lease)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		wch := cli.Watch(cli.Ctx(), p.fullKey, clientv3.WithFilterPut())

		for {
			select {
			case _, ok := <-ch:
				if !ok {
					p.revoke(cli)
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				}

			case c := <-wch:
				if c.Err() != nil {
					logc.Errorf(cli.Ctx(), "etcd publisher watch: %v", c.Err())
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				}

				for _, evt := range c.Events {
					if evt.Type == clientv3.EventTypeDelete {
						logc.Infof(cli.Ctx(), "etcd publisher watch: %s, event: %v",
							evt.Kv.Key, evt.Type)

						// Keep the fast path for manually deleted keys, but revalidate both
						// before and after the re-put so an already-expired lease can't leave
						// a permanent orphaned key behind.
						if !p.leaseAlive(cli) {
							if err := p.restartKeepAlive(cli); err != nil {
								logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
							}
							return
						}

						_, err := cli.Put(cli.Ctx(), p.fullKey, p.value, clientv3.WithLease(p.lease))
						if err != nil {
							logc.Errorf(cli.Ctx(), "etcd publisher re-put key: %v", err)
							continue
						}

						if !p.keyBoundToLease(cli) {
							if err := p.restartKeepAlive(cli); err != nil {
								logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
							}
							return
						}

						logc.Infof(cli.Ctx(), "etcd publisher re-put key: %s, value: %s",
							p.fullKey, p.value)
					}
				}
			case <-p.pauseChan:
				logc.Infof(cli.Ctx(), "paused etcd renew, key: %s, value: %s", p.key, p.value)
				p.revoke(cli)
				select {
				case <-p.resumeChan:
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				case <-p.quit.Done():
					return
				}
			case <-p.quit.Done():
				p.revoke(cli)
				return
			}
		}
	})

	return nil
}

func (p *Publisher) register(client internal.EtcdClient) (clientv3.LeaseID, error) {
	resp, err := client.Grant(client.Ctx(), TimeToLive)
	if err != nil {
		return clientv3.NoLease, err
	}

	lease := resp.ID
	if p.id > 0 {
		p.fullKey = makeEtcdKey(p.key, p.id)
	} else {
		p.fullKey = makeEtcdKey(p.key, int64(lease))
	}
	_, err = client.Put(client.Ctx(), p.fullKey, p.value, clientv3.WithLease(lease))

	return lease, err
}

func (p *Publisher) revoke(cli internal.EtcdClient) {
	if _, err := cli.Revoke(cli.Ctx(), p.lease); err != nil {
		logc.Errorf(cli.Ctx(), "etcd publisher revoke: %v", err)
	}
}

func (p *Publisher) keyBoundToLease(cli internal.EtcdClient) bool {
	resp, err := cli.Get(cli.Ctx(), p.fullKey)
	if err != nil {
		logc.Errorf(cli.Ctx(), "etcd publisher verify re-put lease: %v", err)
		return false
	}

	if len(resp.Kvs) == 0 {
		logc.Errorf(cli.Ctx(), "etcd publisher verify re-put lease: key missing after put, key=%s", p.fullKey)
		return false
	}

	if resp.Kvs[0].Lease != int64(p.lease) {
		logc.Errorf(cli.Ctx(),
			"etcd publisher verify re-put lease: unexpected lease, key=%s, want=%d, got=%d",
			p.fullKey, p.lease, resp.Kvs[0].Lease)
		return false
	}

	return true
}

func (p *Publisher) leaseAlive(cli internal.EtcdClient) bool {
	resp, err := cli.TimeToLive(cli.Ctx(), p.lease)
	if err != nil || resp == nil || resp.TTL <= 0 {
		logc.Errorf(cli.Ctx(),
			"etcd publisher lease expired, skip re-put and restart keepalive: leaseID=%d, err=%v",
			p.lease, err)
		return false
	}

	return true
}

func (p *Publisher) restartKeepAlive(cli internal.EtcdClient) error {
	p.revoke(cli)
	return p.doKeepAlive()
}

// WithId customizes a Publisher with the id.
func WithId(id int64) PubOption {
	return func(publisher *Publisher) {
		publisher.id = id
	}
}

// WithPubEtcdAccount provides the etcd username/password.
func WithPubEtcdAccount(user, pass string) PubOption {
	return func(pub *Publisher) {
		RegisterAccount(pub.endpoints, user, pass)
	}
}

// WithPubEtcdTLS provides the etcd CertFile/CertKeyFile/CACertFile.
func WithPubEtcdTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) PubOption {
	return func(pub *Publisher) {
		logx.Must(RegisterTLS(pub.endpoints, certFile, certKeyFile, caFile, insecureSkipVerify))
	}
}
