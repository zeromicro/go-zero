package internal

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/internal/health"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

var errNotReady = errors.New("service is not ready for a limited time")

// NewRpcPubServer returns a Server.
func NewRpcPubServer(etcd discov.EtcdConf, listenOn string,
	opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubListenOn := figureOutListenOn(listenOn)
		var pubOpts []discov.PubOption
		if etcd.HasAccount() {
			pubOpts = append(pubOpts, discov.WithPubEtcdAccount(etcd.User, etcd.Pass))
		}
		if etcd.HasTLS() {
			pubOpts = append(pubOpts, discov.WithPubEtcdTLS(etcd.CertFile, etcd.CertKeyFile,
				etcd.CACertFile, etcd.InsecureSkipVerify))
		}
		if etcd.HasID() {
			pubOpts = append(pubOpts, discov.WithId(etcd.ID))
		}
		pubClient := discov.NewPublisher(etcd.Hosts, etcd.Key, pubListenOn, pubOpts...)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd:    registerEtcd,
		Server:          NewRpcServer(listenOn, opts...),
		registerTimeout: time.Duration(etcd.RegisterTimeout) * time.Second,
	}

	return server, nil
}

type keepAliveServer struct {
	registerEtcd func() error
	Server
	registerTimeout time.Duration
}

func (s keepAliveServer) Start(fn RegisterFn) error {
	errCh := make(chan error)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go func() {
		defer close(errCh)
		select {
		case errCh <- s.Server.Start(fn):
		case <-stopCh:
			// prevent goroutine leak
		}
	}()

	// Wait for the service to start successfully, otherwise the registration service will fail.
	ctx, cancel := context.WithTimeout(context.Background(), s.registerTimeout)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

l:
	for {
		select {
		case <-ticker.C:
			if health.IsReady() {
				err := s.registerEtcd()
				if err != nil {
					return err
				}
				// break for loop
				break l
			}
		case <-ctx.Done():
			return errNotReady
		case err := <-errCh:
			return err
		}
	}
	ticker.Stop()

	return <-errCh
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
