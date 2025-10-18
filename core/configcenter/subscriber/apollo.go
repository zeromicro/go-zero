package subscriber

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"
)

type (
	// ApolloConf is the configuration for Apollo.
	ApolloConf struct {
		AppID          string
		Cluster        string `json:",default=default"`
		NamespaceName  string `json:",default=application"`
		IP             string `json:",optional"`
		MetaAddr       string // Apollo meta server address (required)
		Secret         string `json:",optional"`
		IsBackupConfig bool   `json:",optional"`
		BackupPath     string `json:",optional"`
		MustStart      bool   `json:",optional"`
		Format         string `json:",default=json,options=json|yaml|properties"`
		Key            string `json:",optional"` // Specific key in namespace, empty means use all content
	}

	// apolloSubscriber is a subscriber that subscribes to Apollo config center.
	apolloSubscriber struct {
		client    agollo.Client
		conf      ApolloConf
		listeners []func()
		lock      sync.RWMutex
		cache     string
		cacheLock sync.RWMutex
	}
)

var (
	// ErrEmptyMetaAddr indicates that Apollo meta server address is empty.
	ErrEmptyMetaAddr = errors.New("empty Apollo meta server address")
	// ErrEmptyAppID indicates that Apollo app ID is empty.
	ErrEmptyAppID = errors.New("empty Apollo app ID")
)

// MustNewApolloSubscriber returns an Apollo Subscriber, exits on errors.
func MustNewApolloSubscriber(conf ApolloConf) Subscriber {
	s, err := NewApolloSubscriber(conf)
	logx.Must(err)
	return s
}

// NewApolloSubscriber returns an Apollo Subscriber.
func NewApolloSubscriber(conf ApolloConf) (Subscriber, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	apolloConf := buildApolloConfig(conf)
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	if err != nil {
		return nil, err
	}

	sub := &apolloSubscriber{
		client: client,
		conf:   conf,
	}

	// Load initial value
	if err := sub.loadValue(); err != nil {
		return nil, err
	}

	// Register change listener
	client.AddChangeListener(sub)

	return sub, nil
}

// Validate validates the ApolloConf.
func (c ApolloConf) Validate() error {
	if len(c.MetaAddr) == 0 {
		return ErrEmptyMetaAddr
	}
	if len(c.AppID) == 0 {
		return ErrEmptyAppID
	}
	return nil
}

// AddListener adds a listener to the subscriber.
func (s *apolloSubscriber) AddListener(listener func()) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners = append(s.listeners, listener)
	return nil
}

// Value returns the value of the subscriber.
func (s *apolloSubscriber) Value() (string, error) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	return s.cache, nil
}

// OnChange is called when Apollo config changes.
// Implements agollo.ChangeListener interface.
func (s *apolloSubscriber) OnChange(event *storage.ChangeEvent) {
	if err := s.loadValue(); err != nil {
		logx.Errorf("Apollo config reload failed: %v", err)
		return
	}

	s.lock.RLock()
	listeners := make([]func(), len(s.listeners))
	copy(listeners, s.listeners)
	s.lock.RUnlock()

	for _, listener := range listeners {
		listener()
	}
}

// OnNewestChange is called when Apollo config changes to newest.
// Implements agollo.ChangeListener interface.
func (s *apolloSubscriber) OnNewestChange(event *storage.FullChangeEvent) {
	// Trigger reload on any change
	if err := s.loadValue(); err != nil {
		logx.Errorf("Apollo config reload failed: %v", err)
		return
	}

	s.lock.RLock()
	listeners := make([]func(), len(s.listeners))
	copy(listeners, s.listeners)
	s.lock.RUnlock()

	for _, listener := range listeners {
		listener()
	}
}

func (s *apolloSubscriber) loadValue() error {
	var value string
	var err error

	// If specific key is set, get that key's value
	if len(s.conf.Key) > 0 {
		val := s.client.GetValue(s.conf.Key)
		if len(val) == 0 {
			return errors.New("key not found in Apollo namespace")
		}
		value = val
	} else {
		// Get all content from namespace
		value, err = s.getAllContent()
		if err != nil {
			return err
		}
	}

	s.cacheLock.Lock()
	s.cache = value
	s.cacheLock.Unlock()

	return nil
}

func (s *apolloSubscriber) getAllContent() (string, error) {
	allConfig := make(map[string]interface{})

	// Get all keys from the namespace
	cache := s.client.GetConfigCache(s.conf.NamespaceName)
	cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			allConfig[k] = value
		}
		return true
	})

	var result []byte
	var err error

	switch s.conf.Format {
	case "json":
		result, err = json.Marshal(allConfig)
	case "yaml":
		result, err = yaml.Marshal(allConfig)
	case "properties":
		// For properties format, convert to key=value format
		props := ""
		for k, v := range allConfig {
			props += k + "=" + toString(v) + "\n"
		}
		result = []byte(props)
	default:
		result, err = json.Marshal(allConfig)
	}

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func buildApolloConfig(conf ApolloConf) *config.AppConfig {
	apolloConf := &config.AppConfig{
		AppID:          conf.AppID,
		Cluster:        conf.Cluster,
		NamespaceName:  conf.NamespaceName,
		IP:             conf.MetaAddr,
		IsBackupConfig: conf.IsBackupConfig,
		MustStart:      conf.MustStart,
	}

	if len(conf.IP) > 0 {
		apolloConf.IP = conf.IP
	}

	if len(conf.Secret) > 0 {
		apolloConf.Secret = conf.Secret
	}

	if len(conf.BackupPath) > 0 {
		apolloConf.BackupConfigPath = conf.BackupPath
	}

	return apolloConf
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	// Try to convert to string via json
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
