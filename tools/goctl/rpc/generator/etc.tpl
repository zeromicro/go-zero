Consul:
  Host: 127.0.0.1:8500 # consul endpoint
  ListenOn: 127.0.0.1:8081
  Key: {{.serviceName}}.rpc
  Meta:
    Protocol: grpc
  Tag:
    - {{.serviceName}}
    - rpc