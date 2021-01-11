package internal

import "github.com/tal-tech/go-zero/core/discov"

func NewRpcPubServer(etcdEndpoints []string, etcdKey, user, pass, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, user, pass, listenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}

type keepAliveServer struct {
	registerEtcd func() error
	Server
}

func (ags keepAliveServer) Start(fn RegisterFn) error {
	if err := ags.registerEtcd(); err != nil {
		return err
	}

	return ags.Server.Start(fn)
}
