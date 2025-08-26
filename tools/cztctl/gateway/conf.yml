Name: gateway-example # gateway name
Host: localhost # gateway host
Port: 8888 # gateway port
Upstreams: # upstreams
  - Grpc: # grpc upstream
      Target: 0.0.0.0:8080 # grpc target,the direct grpc server address,for only one node
#      Endpoints: [0.0.0.0:8080,192.168.120.1:8080] # grpc endpoints, the grpc server address list, for multiple nodes
#      Etcd: # etcd config, if you want to use etcd to discover the grpc server address
#        Hosts: [127.0.0.1:2378,127.0.0.1:2379] # etcd hosts
#        Key: greet.grpc # the discovery key
    # protoset mode
    ProtoSets:
      - hello.pb
    # Mappings can also be written in proto options
#    Mappings: # routes mapping
#      - Method: get
#        Path: /ping
#        RpcPath: hello.Hello/Ping
