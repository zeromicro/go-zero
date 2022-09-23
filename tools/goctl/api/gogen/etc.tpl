Consul:
  Host: 127.0.0.1:8500 # consul endpoint
  ListenOn: {{.host}}:{{.port}}
  Key: {{.serviceName}}.api
  Meta:
    Protocol: grpc
  Tag:
    - {{.serviceName}}
    - api