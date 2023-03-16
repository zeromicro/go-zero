package subscriber

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber/internal"
)

type (
	// NacosSubscriber is a subscriber that subscribes to nacos.
	NacosSubscriber struct {
		conf   NacosConfig
		client config_client.IConfigClient
	}

	// NacosConfig is the configuration for nacos.
	NacosConfig struct {
		internal.NacosConf
		Group  string
		DataID string
	}
)

// MustNewNacosSubscriber returns a nacos Subscriber, exits on errors.
func MustNewNacosSubscriber(conf NacosConfig) Subscriber {
	subscriber, err := NewNacosSubscriber(conf)
	if err != nil {
		log.Fatal(err)
	}

	return subscriber
}

// NewNacosSubscriber returns a nacos Subscriber.
func NewNacosSubscriber(conf NacosConfig) (Subscriber, error) {
	client, err := conf.BuildConfigClient()
	if err != nil {
		return nil, err
	}

	return &NacosSubscriber{
		conf:   conf,
		client: client,
	}, nil
}

// AddListener adds a listener to the subscriber.
func (s *NacosSubscriber) AddListener(listener func()) error {
	return s.client.ListenConfig(vo.ConfigParam{
		DataId: s.conf.DataID,
		Group:  s.conf.Group,
		OnChange: func(_, _, _, _ string) {
			listener()
		},
	})
}

// Value returns the value of the subscriber.
func (s *NacosSubscriber) Value() (string, error) {
	content, err := s.client.GetConfig(vo.ConfigParam{
		DataId: s.conf.DataID,
		Group:  s.conf.Group,
	})
	if err != nil {
		return "", err
	}

	return content, nil
}
