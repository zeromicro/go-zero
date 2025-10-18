package main

import (
	"fmt"
	"time"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

type AppConfig struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Timeout  int64  `json:"timeout"`
	MaxConns int    `json:"maxConns"`
	Features struct {
		EnableCache bool `json:"enableCache"`
		EnableTrace bool `json:"enableTrace"`
	} `json:"features"`
}

func main() {
	// Create Apollo subscriber
	sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
		AppID:          "go-zero-demo",          // Your Apollo AppID
		Cluster:        "default",               // Cluster name
		NamespaceName:  "application.json",      // Namespace name with format suffix
		MetaAddr:       "http://localhost:8080", // Apollo meta server address
		Secret:         "",                      // Optional: Apollo secret key
		IsBackupConfig: true,                    // Enable backup to local file
		BackupPath:     "/tmp/apollo-backup",    // Backup directory
		Format:         "json",                  // Config format: json, yaml, or properties
	})

	// Create config center with type-safe config
	cc := configurator.MustNewConfigCenter[AppConfig](configurator.Config{
		Type: "json",
		Log:  true,
	}, sub)

	// Get initial config
	config, err := cc.GetConfig()
	if err != nil {
		logx.Errorf("Failed to get config: %v", err)
		return
	}

	fmt.Printf("Initial config: %+v\n", config)

	// Add listener for config changes (hot reload)
	cc.AddListener(func() {
		newConfig, err := cc.GetConfig()
		if err != nil {
			logx.Errorf("Failed to get updated config: %v", err)
			return
		}
		fmt.Printf("Config updated: %+v\n", newConfig)

		// Apply new configuration
		// For example, update connection pool size, feature flags, etc.
		logx.Infof("Applying new config - MaxConns: %d, Timeout: %d",
			newConfig.MaxConns, newConfig.Timeout)
	})

	// Keep running to receive config updates
	fmt.Println("Listening for config changes... Press Ctrl+C to exit")
	select {}
}

// Example: Using specific key instead of entire namespace
func exampleSpecificKey() {
	// Get specific key from Apollo
	sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
		AppID:         "go-zero-demo",
		Cluster:       "default",
		NamespaceName: "application",
		MetaAddr:      "http://localhost:8080",
		Key:           "database.url", // Specific key to watch
	})

	cc := configurator.MustNewConfigCenter[string](configurator.Config{
		Type: "json",
	}, sub)

	dbUrl, _ := cc.GetConfig()
	fmt.Printf("Database URL: %s\n", dbUrl)
}

// Example: Using properties format
func examplePropertiesFormat() {
	sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
		AppID:         "go-zero-demo",
		Cluster:       "default",
		NamespaceName: "application.properties",
		MetaAddr:      "http://localhost:8080",
		Format:        "properties", // Properties format
	})

	cc := configurator.MustNewConfigCenter[string](configurator.Config{
		Type: "json",
	}, sub)

	properties, _ := cc.GetConfig()
	fmt.Printf("Properties: %s\n", properties)
}

// Example: Integration with existing go-zero service
type ServiceConfig struct {
	configurator.Configurator[AppConfig]
}

func (sc *ServiceConfig) UpdateConfig(newConfig AppConfig) {
	// Update service behavior based on new config
	logx.Infof("Service config updated: %+v", newConfig)
	// Update connection pools, timeouts, feature flags, etc.
}

func exampleServiceIntegration() {
	sub := subscriber.MustNewApolloSubscriber(subscriber.ApolloConf{
		AppID:    "my-service",
		MetaAddr: "http://localhost:8080",
	})

	cc := configurator.MustNewConfigCenter[AppConfig](
		configurator.Config{Type: "json"},
		sub,
	)

	svc := &ServiceConfig{
		Configurator: cc,
	}

	// Listen for changes
	cc.AddListener(func() {
		config, _ := cc.GetConfig()
		svc.UpdateConfig(config)
	})

	// Service logic continues...
	time.Sleep(time.Hour)
}
