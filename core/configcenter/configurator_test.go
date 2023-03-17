package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigCenter(t *testing.T) {
	_, err := NewConfigCenter[any](Config{}, &mockSubscriber{})
	assert.Error(t, err)

	_, err = NewConfigCenter[any](Config{Type: "json"}, &mockSubscriber{})
	assert.NoError(t, err)
}

func TestConfigCenter_GetConfig(t *testing.T) {
	mock := &mockSubscriber{}
	type Data struct {
		Name string `json:"name"`
	}

	mock.v = `{"name": "go-zero"}`
	c1, err := NewConfigCenter[Data](Config{Type: "json"}, mock)
	assert.NoError(t, err)

	data, err := c1.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "go-zero", data.Name)

	mock.v = `{"name": 111"}`
	c2, err := NewConfigCenter[Data](Config{Type: "json"}, mock)
	assert.NoError(t, err)

	data, err = c2.GetConfig()
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
}

func TestConfigCenter_AddListener(t *testing.T) {
	mock := &mockSubscriber{}
	mock.v = "1234"

	c, err := NewConfigCenter[any](Config{Type: "json"}, mock)
	assert.NoError(t, err)

	cc := c.(*configCenter[any])

	assert.Equal(t, cc.Value(), "1234")
}

func TestConfigCenter_Value(t *testing.T) {
	c, err := NewConfigCenter[any](Config{Type: "json"}, &mockSubscriber{})
	assert.NoError(t, err)

	cc := c.(*configCenter[any])
	cc.AddListener(func() {})
	cc.AddListener(func() {})

	assert.Equal(t, 2, len(cc.listeners))
}

type mockSubscriber struct {
	v        string
	listener func()
}

func (m *mockSubscriber) AddListener(listener func()) error {
	m.listener = listener
	return nil
}

func (m *mockSubscriber) Value() (string, error) {
	return m.v, nil
}

func (m *mockSubscriber) change() {
	if m.listener != nil {
		m.listener()
	}
}
