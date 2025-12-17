package configurator

import (
	"errors"
	"fmt"
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
	errEmptyConfig            = errors.New("empty config value")
	errMissingUnmarshalerType = errors.New("missing unmarshaler type")
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
		// Log is the flag to control logging.
		Log bool `json:",default=true"`
	}

	configCenter[T any] struct {
		conf        Config
		unmarshaler LoaderFn
		subscriber  subscriber.Subscriber
		listeners   []func()
		lock        sync.Mutex
		snapshot    atomic.Value
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
	logx.Must(err)
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
	}

	if err := cc.loadConfig(); err != nil {
		return nil, err
	}

	if err := cc.subscriber.AddListener(cc.onChange); err != nil {
		return nil, err
	}

	if _, err := cc.GetConfig(); err != nil {
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
	v := c.value()
	if v == nil || len(v.data) == 0 {
		var empty T
		return empty, errEmptyConfig
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
	if err := c.loadConfig(); err != nil {
		return
	}

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
	if len(data) == 0 {
		return v
	}

	t := reflect.TypeOf(v.marshalData)
	// if the type is nil, it means that the user has not set the type of the configuration.
	if t == nil {
		v.err = errMissingUnmarshalerType
		return v
	}

	t = mapping.Deref(t)
	switch t.Kind() {
	case reflect.Struct, reflect.Array, reflect.Slice:
		if err := c.unmarshaler([]byte(data), &v.marshalData); err != nil {
			v.err = err
			if c.conf.Log {
				logx.Errorf("ConfigCenter unmarshal configuration failed, err: %+v, content [%s]",
					err.Error(), data)
			}
		}
	case reflect.String:
		if str, ok := any(data).(T); ok {
			v.marshalData = str
		} else {
			v.err = errMissingUnmarshalerType
		}
	default:
		if c.conf.Log {
			logx.Errorf("ConfigCenter unmarshal configuration missing unmarshaler for type: %s, content [%s]",
				t.Kind(), data)
		}
		v.err = errMissingUnmarshalerType
	}

	return v
}
