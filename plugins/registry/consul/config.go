package consul

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/hashicorp/consul/api"
)

const (
	allEths    = "0.0.0.0"
	envPodIP   = "POD_IP"
	consulTags = "consul_tags"
)

// Conf is the config item with the given key on etcd.
type Conf struct {
	Host     string
	ListenOn string
	Key      string
	Token    string            `json:",optional"`
	Tag      []string          `json:",optional"`
	Meta     map[string]string `json:",optional"`
	TTL      int               `json:"ttl,optional"`
}

// Validate validates c.
func (c Conf) Validate() error {
	if len(c.Host) == 0 {
		return errors.New("empty consul hosts")
	}
	if len(c.Key) == 0 {
		return errors.New("empty consul key")
	}
	if c.TTL == 0 {
		c.TTL = 20
	}

	return nil
}

// NewClient create new client
func (c Conf) NewClient() (*api.Client, error) {
	client, err := api.NewClient(&api.Config{Scheme: "http", Address: c.Host, Token: c.Token})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// LoadYAMLConf load config from consul kv
func LoadYAMLConf(client *api.Client, key string, v any) {
	kv := client.KV()

	data, _, err := kv.Get(key, nil)
	err = yaml.Unmarshal(data.Value, v)
	logx.Must(err)
}

// LoadJSONConf load config from consul kv
func LoadJSONConf(client *api.Client, key string, v any) {
	kv := client.KV()

	data, _, err := kv.Get(key, nil)
	err = json.Unmarshal(data.Value, v)
	logx.Must(err)
}
