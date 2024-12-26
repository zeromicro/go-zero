package configurator

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigCenter(t *testing.T) {
	_, err := NewConfigCenter[any](Config{
		Log: true,
	}, &mockSubscriber{})
	assert.Error(t, err)

	_, err = NewConfigCenter[any](Config{
		Type: "json",
		Log:  true,
	}, &mockSubscriber{})
	assert.Error(t, err)
}

func TestConfigCenter_GetConfig(t *testing.T) {
	mock := &mockSubscriber{}
	type Data struct {
		Name string `json:"name"`
	}

	mock.v = `{"name": "go-zero"}`
	c1, err := NewConfigCenter[Data](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.NoError(t, err)

	data, err := c1.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "go-zero", data.Name)

	mock.v = `{"name": "111"}`
	c2, err := NewConfigCenter[Data](Config{Type: "json"}, mock)
	assert.NoError(t, err)

	mock.v = `{}`
	c3, err := NewConfigCenter[string](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.NoError(t, err)
	_, err = c3.GetConfig()
	assert.NoError(t, err)

	data, err = c2.GetConfig()
	assert.NoError(t, err)

	mock.lisErr = errors.New("mock error")
	_, err = NewConfigCenter[Data](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.Error(t, err)
}

func TestConfigCenter_onChange(t *testing.T) {
	mock := &mockSubscriber{}
	type Data struct {
		Name string `json:"name"`
	}

	mock.v = `{"name": "go-zero"}`
	c1, err := NewConfigCenter[Data](Config{Type: "json", Log: true}, mock)
	assert.NoError(t, err)

	data, err := c1.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "go-zero", data.Name)

	mock.v = `{"name": "go-zero2"}`
	mock.change()

	data, err = c1.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "go-zero2", data.Name)

	mock.valErr = errors.New("mock error")
	_, err = NewConfigCenter[Data](Config{Type: "json", Log: false}, mock)
	assert.Error(t, err)
}

func TestConfigCenter_Value(t *testing.T) {
	mock := &mockSubscriber{}
	mock.v = "1234"

	c, err := NewConfigCenter[string](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.NoError(t, err)

	cc := c.(*configCenter[string])

	assert.Equal(t, cc.Value(), "1234")

	mock.valErr = errors.New("mock error")

	_, err = NewConfigCenter[any](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.Error(t, err)
}

func TestConfigCenter_AddListener(t *testing.T) {
	mock := &mockSubscriber{}
	mock.v = "1234"
	c, err := NewConfigCenter[string](Config{
		Type: "json",
		Log:  true,
	}, mock)
	assert.NoError(t, err)

	cc := c.(*configCenter[string])
	var a, b int
	var mutex sync.Mutex
	cc.AddListener(func() {
		mutex.Lock()
		a = 1
		mutex.Unlock()
	})
	cc.AddListener(func() {
		mutex.Lock()
		b = 2
		mutex.Unlock()
	})

	assert.Equal(t, 2, len(cc.listeners))

	mock.change()

	time.Sleep(time.Millisecond * 100)

	mutex.Lock()
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	mutex.Unlock()
}

func TestConfigCenter_genValue(t *testing.T) {
	t.Run("data is empty", func(t *testing.T) {
		c := &configCenter[string]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue("")
		assert.Equal(t, "", v.data)
	})

	t.Run("invalid template type", func(t *testing.T) {
		c := &configCenter[any]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue("xxxx")
		assert.Equal(t, errMissingUnmarshalerType, v.err)
	})

	t.Run("unsupported template type", func(t *testing.T) {
		c := &configCenter[int]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue("1")
		assert.Equal(t, errMissingUnmarshalerType, v.err)
	})

	t.Run("supported template string type", func(t *testing.T) {
		c := &configCenter[string]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue("12345")
		assert.NoError(t, v.err)
		assert.Equal(t, "12345", v.data)
	})

	t.Run("unmarshal fail", func(t *testing.T) {
		c := &configCenter[struct {
			Name string `json:"name"`
		}]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue(`{"name":"new name}`)
		assert.Equal(t, `{"name":"new name}`, v.data)
		assert.Error(t, v.err)
	})

	t.Run("success", func(t *testing.T) {
		c := &configCenter[struct {
			Name string `json:"name"`
		}]{
			unmarshaler: registry.unmarshalers["json"],
			conf:        Config{Log: true},
		}
		v := c.genValue(`{"name":"new name"}`)
		assert.Equal(t, `{"name":"new name"}`, v.data)
		assert.Equal(t, "new name", v.marshalData.Name)
		assert.NoError(t, v.err)
	})
}

type mockSubscriber struct {
	v              string
	lisErr, valErr error
	listener       func()
}

func (m *mockSubscriber) AddListener(listener func()) error {
	m.listener = listener
	return m.lisErr
}

func (m *mockSubscriber) Value() (string, error) {
	return m.v, m.valErr
}

func (m *mockSubscriber) change() {
	if m.listener != nil {
		m.listener()
	}
}
