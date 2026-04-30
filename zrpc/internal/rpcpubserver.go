package internal

import (
	"os"
	"strings"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/netx"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

type etcdPublisher interface {
	KeepAlive() error
	Pause()
	Resume()
}

// NewRpcPubServer returns a Server.
func NewRpcPubServer(etcd discov.EtcdConf, listenOn string,
	opts ...ServerOption) (Server, error) {
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
	server := keepAliveServer{
		publisher: pubClient,
		Server:    NewRpcServer(listenOn, opts...),
	}

	return server, nil
}

type keepAliveServer struct {
	publisher etcdPublisher
	Server
}

func (s keepAliveServer) Start(fn RegisterFn) error {
	if err := s.publisher.KeepAlive(); err != nil {
		return err
	}

	return s.Server.Start(fn)
}

// PauseEtcdRegister pauses the etcd lease renewal for this rpc server.
func (s keepAliveServer) PauseEtcdRegister() {
	if s.publisher != nil {
		s.publisher.Pause()
	}
}

// ResumeEtcdRegister resumes the etcd lease renewal for this rpc server.
func (s keepAliveServer) ResumeEtcdRegister() {
	if s.publisher != nil {
		s.publisher.Resume()
	}
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
