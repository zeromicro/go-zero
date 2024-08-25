package configurator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/threading"
)

var (
	ErrorEmptyConfig         = errors.New("empty config value")
	ErrorMissUnmarshalerType = errors.New("miss unmarshaler type")
)

// Configurator is the interface for configuration center.
type Configurator[T any] interface {
	// GetConfig returns the subscription value.
	GetConfig() (T, error)
	// AddListener adds a listener to the subscriber.
	AddListener(listener func())
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

	value[T any] struct {
		data        string
		marshalData T
		err         error
	}
)

// Configurator is the interface for configuration center.
var _ Configurator[any] = (*configCenter[any])(nil)

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
	v := c.value()
	if v == nil || len(v.data) < 1 {
		return r, ErrorEmptyConfig
	}

	return v.marshalData, v.err
}

// Value returns the subscription value.
func (c *configCenter[T]) Value() string {
	v := c.value()
	if v == nil {
		return ""
	}
	return v.data
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

	c.snapshot.Store(c.genValue(v))
	return nil
}

func (c *configCenter[T]) onChange() {
	_ = c.loadConfig()

	c.lock.Lock()
	listeners := make([]func(), len(c.listeners))
	copy(listeners, c.listeners)
	c.lock.Unlock()

	for _, l := range listeners {
		threading.GoSafe(l)
	}
}

func (c *configCenter[T]) value() *value[T] {
	content := c.snapshot.Load()
	if content == nil {
		return nil
	}
	return content.(*value[T])
}

func (c *configCenter[T]) genValue(data string) *value[T] {
	v := &value[T]{
		data: data,
	}
	if len(data) <= 0 {
		return v
	}

	t := reflect.TypeOf(v.marshalData)
	// if the type is nil, it means that the user has not set the type of the configuration.
	if t == nil {
		v.err = ErrorMissUnmarshalerType
		return v
	}

	t = mapping.Deref(t)

	switch t.Kind() {
	case reflect.Struct, reflect.Array, reflect.Slice:
		err := c.unmarshaler([]byte(data), &v.marshalData)
		if err != nil {
			v.err = err
			if c.conf.Log {
				logx.Errorf("ConfigCenter unmarshal configuration failed, err: %+v, content [%s]", err.Error(), data)
			}
		}
	case reflect.String:
		if str, ok := any(data).(T); ok {
			v.marshalData = str
		} else {
			v.err = ErrorMissUnmarshalerType
		}
	default:
		if c.conf.Log {
			logx.Errorf("ConfigCenter unmarshal configuration missing unmarshaler for type: %s, content [%s]",
				t.Kind(), data)
		}
		v.err = ErrorMissUnmarshalerType
	}

	return v
}
