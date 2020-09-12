package discov

import (
	"github.com/tal-tech/go-zero/core/discov/internal"
	"github.com/tal-tech/go-zero/core/logx"
)

type (
	Facade struct {
		endpoints []string
		registry  *internal.Registry
	}

	FacadeListener interface {
		OnAdd(key, val string)
		OnDelete(key string)
	}
)

func NewFacade(endpoints []string) Facade {
	return Facade{
		endpoints: endpoints,
		registry:  internal.GetRegistry(),
	}
}

func (f Facade) Client() internal.EtcdClient {
	conn, err := f.registry.GetConn(f.endpoints)
	logx.Must(err)
	return conn
}

func (f Facade) Monitor(key string, l FacadeListener) {
	f.registry.Monitor(f.endpoints, key, listenerAdapter{l})
}

type listenerAdapter struct {
	l FacadeListener
}

func (la listenerAdapter) OnAdd(kv internal.KV) {
	la.l.OnAdd(kv.Key, kv.Val)
}

func (la listenerAdapter) OnDelete(kv internal.KV) {
	la.l.OnDelete(kv.Key)
}
