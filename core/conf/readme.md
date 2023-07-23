## How to use

1. Define a config structure, like below:

```go
type RestfulConf struct {
  ServiceName  string        `json:",env=SERVICE_NAME"`  // read from env automatically
	Host         string        `json:",default=0.0.0.0"`
	Port         int
	LogMode      string        `json:",options=[file,console]"`
	Verbose      bool          `json:",optional"`
	MaxConns     int           `json:",default=10000"`
	MaxBytes     int64         `json:",default=1048576"`
	Timeout      time.Duration `json:",default=3s"`
	CpuThreshold int64         `json:",default=900,range=[0:1000]"`
}
```

2. Write the yaml, toml or json config file:

- yaml example

```yaml
# most fields are optional or have default values
port: 8080
logMode: console
# you can use env settings
maxBytes: ${MAX_BYTES}
```

- toml example

```toml
# most fields are optional or have default values
port = 8_080
logMode = "console"
# you can use env settings
maxBytes = "${MAX_BYTES}"
```

3. Load the config from a file:

```go
// exit on error
var config RestfulConf
conf.MustLoad(configFile, &config)

// or handle the error on your own
var config RestfulConf
if err := conf.Load(configFile, &config); err != nil {
  log.Fatal(err)
}

// enable reading from environments
var config RestfulConf
conf.MustLoad(configFile, &config, conf.UseEnv())
```

