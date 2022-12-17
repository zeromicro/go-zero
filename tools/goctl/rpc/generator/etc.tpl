Name: {{.serviceName}}.rpc
ListenOn: 0.0.0.0:{{.port}}
{{if .isEnt}}
DatabaseConf:
  Type: mysql
  Host: 127.0.0.1
  Port: 3306
  DBName: simple_admin
  Username: # set your username
  Password: # set your password
  MaxOpenConn: 100
  SSLMode: disable
  CacheTime: 5

RedisConf:
  Host: 127.0.0.1:6379
  Type: node{{end}}

Log:
  ServiceName: {{.serviceName}}RpcLogger
  Mode: file
  Path: /home/data/logs/{{.serviceName}}/rpc
  Encoding: json
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100

Prometheus:
  Host: 0.0.0.0
  Port: 4001
  Path: /metrics
