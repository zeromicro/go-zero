package rest

import (
	"net"
	"os"
	"strconv"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/netx"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

// HasEtcd checks if there is etcd settings in config.
func (sc RestConf) HasEtcd() bool {
	return len(sc.Etcd.Hosts) > 0 && len(sc.Etcd.Key) > 0
}

// figureOutListenAddr figures out the listen address for etcd registration.
func figureOutListenAddr(host string, port int) string {
	if host != allEths && host != "" {
		return net.JoinHostPort(host, strconv.Itoa(port))
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		if host == "" {
			host = allEths
		}
		return net.JoinHostPort(host, strconv.Itoa(port))
	}

	return net.JoinHostPort(ip, strconv.Itoa(port))
}

// registerToEtcd registers the REST service to etcd.
func (s *Server) registerToEtcd() error {
	if !s.conf.HasEtcd() {
		return nil
	}

	listenAddr := figureOutListenAddr(s.conf.Host, s.conf.Port)
	var pubOpts []discov.PubOption
	if s.conf.Etcd.HasAccount() {
		pubOpts = append(pubOpts, discov.WithPubEtcdAccount(s.conf.Etcd.User, s.conf.Etcd.Pass))
	}
	if s.conf.Etcd.HasTLS() {
		pubOpts = append(pubOpts, discov.WithPubEtcdTLS(s.conf.Etcd.CertFile, s.conf.Etcd.CertKeyFile,
			s.conf.Etcd.CACertFile, s.conf.Etcd.InsecureSkipVerify))
	}
	if s.conf.Etcd.HasID() {
		pubOpts = append(pubOpts, discov.WithId(s.conf.Etcd.ID))
	}

	pubClient := discov.NewPublisher(s.conf.Etcd.Hosts, s.conf.Etcd.Key, listenAddr, pubOpts...)
	return pubClient.KeepAlive()
}
