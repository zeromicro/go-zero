package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	config := ConsumerConfig{Client: ClientConfig{
		ResourceName: "kafka-test-1",
	}}
	_, err := NewConsumer(config)
	assert.Error(t, err)
}
