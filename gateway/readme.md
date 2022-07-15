# Gateway

## Usage

- main.go

```go
var configFile = flag.String("f", "config.yaml", "config file")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c)
	defer gw.Stop()
	gw.Start()
}
```

- config.yaml

```yaml
Name: demo-gateway
Host: localhost
Port: 8888
Upstreams:
  - Grpc:
      Etcd:
        Hosts:
        - localhost:2379
        Key: hello.rpc
    ProtoSet: hello.pb
    Mapping:
      - Path: /pingHello
        Method: hello.Hello/Ping
  - Grpc:
      Endpoints:
        - localhost:8081
    Mapping:
      - Path: /pingWorld
        Method: world.World/Ping
```

