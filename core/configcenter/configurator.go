package configurator

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

// Configurator is the interface for configuration center.
type Configurator[T any] interface {
	// GetConfig returns the subscription value.
	GetConfig() (T, error)
}

type (
	// Config is the configuration for Configurator.
	Config struct {
		// Type is the value type, yaml, json or toml.
		Type string `json:",default=yaml,options=[yaml,json,toml]"`
		// Log indicates whether to log the configuration.
		Log bool `json:",default=ture"`
	}

	configCenter[T any] struct {
		conf        Config
		unmarshaler LoaderFn

		subscriber subscriber.Subscriber

		listeners []func()
		lock      sync.Mutex
		snapshot  atomic.Value
	}
)

// MustNewConfigCenter returns a Configurator, exits on errors.
func MustNewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) Configurator[T] {
	cc, err := NewConfigCenter[T](c, subscriber)
	if err != nil {
		log.Fatalf("NewConfigCenter failed: %v", err)
	}

	_, err = cc.GetConfig()
	if err != nil {
		log.Fatalf("NewConfigCenter.GetConfig failed: %v", err)
	}

	return cc
}

// NewConfigCenter returns a Configurator.
func NewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) (Configurator[T], error) {
	unmarshaler, ok := Unmarshaler(strings.ToLower(c.Type))
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", c.Type)
	}

	cc := &configCenter[T]{
		conf:        c,
		unmarshaler: unmarshaler,
		subscriber:  subscriber,
		listeners:   nil,
		lock:        sync.Mutex{},
		snapshot:    atomic.Value{},
	}

	if err := cc.loadConfig(); err != nil {
		return nil, err
	}

	if err := cc.subscriber.AddListener(cc.onChange); err != nil {
		return nil, err
	}

	return cc, nil
}

// AddListener adds listener to s.
func (c *configCenter[T]) AddListener(listener func()) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.listeners = append(c.listeners, listener)
}

// GetConfig return structured config.
func (c *configCenter[T]) GetConfig() (T, error) {
	var r T
	err := c.unmarshaler([]byte(c.Value()), &r)
	return r, err
}

// Value returns the subscription value.
func (c *configCenter[T]) Value() string {
	content := c.snapshot.Load()
	if content == nil {
		return ""
	}
	return content.(string)
}

func (c *configCenter[T]) loadConfig() error {
	v, err := c.subscriber.Value()
	if err != nil {
		if c.conf.Log {
			logx.Errorf("ConfigCenter loads changed configuration, error: %v", err)
		}
		return err
	}

	if c.conf.Log {
		logx.Infof("ConfigCenter loads changed configuration, content [%s]", v)
	}

	c.snapshot.Store(v)
	return nil
}

func (c *configCenter[T]) onChange() {
	c.loadConfig()

	c.lock.Lock()
	listeners := make([]func(), len(c.listeners))
	copy(listeners, c.listeners)
	c.lock.Unlock()

	for _, l := range listeners {
		threading.GoSafe(l)
	}
}
