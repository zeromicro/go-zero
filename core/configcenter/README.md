# Config Center

Dynamic configuration management with hot reload support for go-zero applications.

## Features

- **Hot Reload**: Automatically reload configuration without service restart
- **Type Safety**: Generic `Configurator[T]` ensures compile-time type checking
- **Multiple Formats**: Support for JSON, YAML, and TOML
- **Multiple Backends**: etcd, Apollo, and extensible for more
- **Thread-Safe**: Concurrent access to configuration snapshots
- **Change Listeners**: Register callbacks for configuration updates

## Supported Config Centers

### etcd

etcd-based dynamic configuration with watch capability.

```go
import (
    "github.com/zeromicro/go-zero/core/configcenter"
    "github.com/zeromicro/go-zero/core/configcenter/subscriber"
)

type AppConfig struct {
    Timeout  int64  `json:"timeout"`
    MaxConns int    `json:"maxConns"`
}

// Create etcd subscriber
sub := subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
    Hosts: []string{"localhost:2379"},
    Key:   "config/app",
})

// Create config center
cc := configcenter.MustNewConfigCenter[AppConfig](configcenter.Config{
    Type: "json",
}, sub)

// Get config
config, _ := cc.GetConfig()

// Listen for changes
cc.AddListener(func() {
    newConfig, _ := cc.GetConfig()
    // Apply new configuration
})
```

### Apollo

Apollo config center integration using [agollo](https://github.com/apolloconfig/agollo).

```go
import (
    "github.com/zeromicro/go-zero/core/configcenter"
    "github.com/zeromicro/go-zero/core/configcenter/subscriber"
)

// Create Apollo subscriber
sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "go-zero-app",
    Cluster:       "default",
    NamespaceName: "application.json",
    MetaAddr:      "http://localhost:8080",
    Secret:        "",  // Optional
    Format:        "json",
})

// Create config center
cc := configcenter.MustNewConfigCenter[AppConfig](configcenter.Config{
    Type: "json",
}, sub)
```

#### Apollo Configuration Options

```go
type ApolloConf struct {
    AppID          string   // Apollo App ID (required)
    Cluster        string   // Cluster name, default: "default"
    NamespaceName  string   // Namespace name, default: "application"
    IP             string   // Client IP for grayscale release
    MetaAddr       string   // Apollo meta server address (required)
    Secret         string   // Secret key for authentication
    IsBackupConfig bool     // Enable local backup, default: true
    BackupPath     string   // Backup file path
    MustStart      bool     // Panic if connect fails, default: true
    Format         string   // Config format: json/yaml/properties, default: "json"
    Key            string   // Specific key in namespace (optional, empty = all content)
}
```

#### Apollo Usage Patterns

**Pattern 1: Entire Namespace (All Keys)**

```go
sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "my-app",
    NamespaceName: "application.json",
    MetaAddr:      "http://apollo:8080",
    Format:        "json",  // Return all keys as JSON object
})

cc := configcenter.MustNewConfigCenter[AppConfig](
    configcenter.Config{Type: "json"},
    sub,
)
```

**Pattern 2: Specific Key**

```go
sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "my-app",
    NamespaceName: "application",
    MetaAddr:      "http://apollo:8080",
    Key:           "database.url",  // Only watch this key
})

cc := configcenter.MustNewConfigCenter[string](
    configcenter.Config{Type: "json"},
    sub,
)

dbUrl, _ := cc.GetConfig()
```

**Pattern 3: Properties Format**

```go
sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "my-app",
    NamespaceName: "application.properties",
    MetaAddr:      "http://apollo:8080",
    Format:        "properties",  // key=value format
})
```

**Pattern 4: With Authentication**

```go
sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:    "my-app",
    MetaAddr: "http://apollo:8080",
    Secret:   "your-secret-key",  // For secure access
})
```

## Usage Examples

### Basic Usage

```go
type Config struct {
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    Database struct {
        DSN         string `json:"dsn"`
        MaxOpenConn int    `json:"maxOpenConn"`
    } `json:"database"`
}

// Initialize
cc := configcenter.MustNewConfigCenter[Config](
    configcenter.Config{Type: "json"},
    subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
        Hosts: []string{"localhost:2379"},
        Key:   "myapp/config",
    }),
)

// Get current config
config, err := cc.GetConfig()
if err != nil {
    log.Fatal(err)
}

// Use config
fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
```

### Hot Reload

```go
// Register change listener
cc.AddListener(func() {
    newConfig, err := cc.GetConfig()
    if err != nil {
        logx.Errorf("Failed to reload config: %v", err)
        return
    }

    // Apply new configuration
    updateDatabasePool(newConfig.Database)
    updateServerSettings(newConfig.Server)

    logx.Info("Configuration reloaded successfully")
})
```

### Integration with go-zero Services

```go
type ServiceContext struct {
    Config       config.Config
    Configurator configcenter.Configurator[DynamicConfig]
}

func NewServiceContext(c config.Config) *ServiceContext {
    // Static config from file
    svcCtx := &ServiceContext{
        Config: c,
    }

    // Dynamic config from Apollo
    if c.Apollo.Enabled {
        sub := subscriber.MustNewApolloSubscriber(c.Apollo)
        cc := configcenter.MustNewConfigCenter[DynamicConfig](
            configcenter.Config{Type: "json"},
            sub,
        )

        svcCtx.Configurator = cc

        // Listen for changes
        cc.AddListener(func() {
            dynamicConf, _ := cc.GetConfig()
            svcCtx.applyDynamicConfig(dynamicConf)
        })
    }

    return svcCtx
}

func (svc *ServiceContext) applyDynamicConfig(cfg DynamicConfig) {
    // Update feature flags
    // Update rate limits
    // Update circuit breaker thresholds
    // etc.
}
```

### Multiple Namespaces

```go
// Database config from one namespace
dbSub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "my-app",
    NamespaceName: "database.json",
    MetaAddr:      apolloAddr,
})
dbCC := configcenter.MustNewConfigCenter[DatabaseConfig](
    configcenter.Config{Type: "json"},
    dbSub,
)

// Feature flags from another namespace
featureSub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
    AppID:         "my-app",
    NamespaceName: "features.json",
    MetaAddr:      apolloAddr,
})
featureCC := configcenter.MustNewConfigCenter[FeatureFlags](
    configcenter.Config{Type: "json"},
    featureSub,
)
```

## Architecture

```
┌─────────────────────────────────────────────┐
│         Application Code                    │
│  ┌─────────────────────────────────────┐   │
│  │  Configurator[T]                    │   │
│  │  - GetConfig() -> T                 │   │
│  │  - AddListener(func())              │   │
│  └──────────────┬──────────────────────┘   │
│                 │                           │
│  ┌──────────────▼──────────────────────┐   │
│  │  ConfigCenter Implementation        │   │
│  │  - Type-safe generic                │   │
│  │  - Thread-safe snapshot             │   │
│  │  - Format unmarshaling              │   │
│  └──────────────┬──────────────────────┘   │
│                 │                           │
│  ┌──────────────▼──────────────────────┐   │
│  │  Subscriber Interface               │   │
│  │  - AddListener(func()) error        │   │
│  │  - Value() (string, error)          │   │
│  └──────────────┬──────────────────────┘   │
└─────────────────┼───────────────────────────┘
                  │
        ┌─────────┴──────────┐
        │                    │
┌───────▼──────┐    ┌────────▼────────┐
│ etcd         │    │ Apollo          │
│ Subscriber   │    │ Subscriber      │
└──────┬───────┘    └────────┬────────┘
       │                     │
┌──────▼───────┐    ┌────────▼────────┐
│ etcd Cluster │    │ Apollo Server   │
└──────────────┘    └─────────────────┘
```

## Best Practices

1. **Use Type Safety**: Define strong types for your configuration
2. **Validate Config**: Add validation in change listeners
3. **Graceful Degradation**: Handle config reload failures gracefully
4. **Separate Static/Dynamic**: Use file-based config for static values, config center for dynamic values
5. **Monitor Changes**: Log configuration changes for debugging
6. **Test Locally**: Use local etcd/Apollo for development
7. **Backup Config**: Enable Apollo backup for offline capability

## Extending with Custom Subscribers

Implement the `Subscriber` interface:

```go
type Subscriber interface {
    AddListener(listener func()) error
    Value() (string, error)
}
```

Example custom subscriber for Consul, Nacos, or other config centers.

## Examples

See [examples/apollo](./examples/apollo) for complete working examples.

## Dependencies

- etcd: `go.etcd.io/etcd/client/v3`
- Apollo: `github.com/apolloconfig/agollo/v4`
