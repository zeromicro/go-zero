package internal

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	errEmptyServerConfigs = errors.New("empty nacos ServerConfigs")
)

const (
	namespacePublic            = "public"
	namespacePublicReplacement = ""
)

type (
	// ClientOption is alias of nacos ClientOption.
	ClientOption = constant.ClientOption

	// ServerConfig is nacos server config,
	// Address is host:port of http protocol,
	// GrpcPort is port of grpc protocol which is optional.
	ServerConfig struct {
		Address  string
		GrpcPort uint64 `json:",optional"`
	}

	// NacosConf is the config item with the given key on nacos.
	NacosConf struct {
		// nacos server config.
		ServerConfigs []ServerConfig // nacos server address config

		// nacos client config.
		TimeoutMs            uint64 `json:",default=10000"`            // timeout for requesting Nacos server, default value is 10000ms
		AppName              string `json:",optional"`                 // the appName
		NamespaceId          string `json:",optional"`                 // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
		CacheDir             string `json:",default=data/nacos/cache"` // the directory for persist nacos service info,default value is current path
		NotLoadCacheAtStart  bool   `json:",default=true"`             // not to load persistent nacos service info in CacheDir at start time
		UpdateCacheWhenEmpty bool   `json:",default=true"`             // update cache when get empty service instance from server
		Username             string `json:",optional"`                 // the username for nacos auth
		Password             string `json:",optional,security"`        // the password for nacos auth
		LogDir               string `json:",default=data/nacos/log"`   // the directory for log, default is current path
		LogLevel             string `json:",default=info"`             // the level of log, it's must be debug,info,warn,error, default value is info
	}
)

// BuildNamingClient create a nacos naming instance from current config.
func (c NacosConf) BuildNamingClient(opts ...ClientOption) (naming_client.INamingClient, error) {
	config, err := c.clientParam(opts...)
	if err != nil {
		return nil, err
	}

	return clients.NewNamingClient(*config)
}

// BuildConfigClient create a nacos config instance from current config.
func (c NacosConf) BuildConfigClient(opts ...ClientOption) (config_client.IConfigClient, error) {
	config, err := c.clientParam(opts...)
	if err != nil {
		return nil, err
	}

	return clients.NewConfigClient(*config)
}

// Validate validate the config.
func (c NacosConf) Validate() error {
	if len(c.ServerConfigs) == 0 {
		return errEmptyServerConfigs
	}

	return nil
}

// clientParam validate config and convert it into nacos vo.NacosClientParam.
func (c NacosConf) clientParam(opts ...ClientOption) (*vo.NacosClientParam, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	serverConfigs := make([]constant.ServerConfig, 0)

	for _, sc := range c.ServerConfigs {
		h, p, err := net.SplitHostPort(sc.Address)
		if err != nil {
			return nil, fmt.Errorf("invalid address: %w", err)
		}

		port, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}

		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr:   h,
			Port:     uint64(port),
			GrpcPort: sc.GrpcPort,
		})
	}

	cc := c.clientConfig()

	for _, opt := range opts {
		opt(cc)
	}

	return &vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: serverConfigs,
	}, nil
}

func (c NacosConf) clientConfig() *constant.ClientConfig {
	// the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	namespaceId := c.NamespaceId
	if namespaceId == namespacePublic {
		namespaceId = namespacePublicReplacement
	}

	return &constant.ClientConfig{
		TimeoutMs:            c.TimeoutMs,
		NamespaceId:          namespaceId,
		CacheDir:             c.CacheDir,
		NotLoadCacheAtStart:  c.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty: c.UpdateCacheWhenEmpty,
		Username:             c.Username,
		Password:             c.Password,
		LogDir:               c.LogDir,
		LogLevel:             c.LogLevel,
	}
}
