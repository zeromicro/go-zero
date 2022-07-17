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
    # protoset mode
    ProtoSet: hello.pb
    Mapping:
      - Method: get
        Path: /pingHello/:ping
        RpcPath: hello.Hello/Ping
  - Grpc:
      Endpoints:
        - localhost:8081
    # reflection mode, no ProtoSet settings
    Mapping:
      - Method: post
        Path: /pingWorld
        RpcPath: world.World/Ping
```

## Generate ProtoSet files

- example command

```shell
protoc --descriptor_set_out=hello.pb hello.proto
```

