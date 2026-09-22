package logx

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// DefaultZapLogger must be initialize before zrpc.MustNewServer
	DefaultZapLogger *clientv3.Config = nil
)
