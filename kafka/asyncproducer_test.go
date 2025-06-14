package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAsyncProducer(t *testing.T) {
	pc := ProducerConfig{Client: ClientConfig{ResourceName: "kafka-test-1"}, SharedProducerConfig: SharedProducerConfig{EnableRecovery: false}}
	_, err := NewAsyncProducer(pc)
	assert.Error(t, err)
}
