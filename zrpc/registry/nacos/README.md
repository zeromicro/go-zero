### Quick Start

Prerequisites:

Download the module:

```console
go get -u github.com/zeromicro/zero-contrib/zrpc/registry/nacos
```

For example:

## Service

- main.go

```go
import _ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {

	})
	// register service to nacos
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.100.15", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	opts := nacos.NewNacosConfig("nacos.rpc", c.ListenOn, sc, cc)
	_ = nacos.RegisterService(opts)
	server.Start()
}
```

## Client

- main.go

```go
import _ "github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
```

- etc/\*.yaml

```yaml
# nacos://[user:passwd]@host/service?param=value'
Target: nacos://192.168.100.15:8848/nacos.rpc?namespaceid=public&timeout=5000s
```
