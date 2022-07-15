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
  - Target: etcd://localhost:2379/hello.rpc
    Mapping:
      - Path: /pingHello
        Method: hello.Hello/Ping
  - Target: etcd://localhost:2379/world.rpc
    Mapping:
      - Path: /pingWorld
        Method: world.World/Ping
```

