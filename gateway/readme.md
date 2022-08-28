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
    ProtoSets:
      - hello.pb
    # Mappings can also be written in proto options
    Mappings:
      - Method: get
        Path: /pingHello/:ping
        RpcPath: hello.Hello/Ping
  - Grpc:
      Endpoints:
        - localhost:8081
    # reflection mode, no ProtoSet settings
    Mappings:
      - Method: post
        Path: /pingWorld
        RpcPath: world.World/Ping
```

## Generate ProtoSet files

- example command without external imports

```shell
protoc --descriptor_set_out=hello.pb hello.proto
```

- example command with external imports

```shell
protoc --include_imports --proto_path=. --descriptor_set_out=hello.pb hello.proto
```
