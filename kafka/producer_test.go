package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func TestNewProducer(t *testing.T) {
	pc := ProducerConfig{Client: ClientConfig{ResourceName: "kafka-test-1"}}
	_, err := NewProducer(pc)
	assert.Error(t, err)
}

func TestBuildDelayMessage(t *testing.T) {
	var m = Message{
		Topic: "test",
	}
	m.BuildDelayMessage(5)
	assert.Equal(t, m.Topic, types.NeptuneTopic)
	assert.NotEmpty(t, m.GetHeader(types.DelayKafkaMsgId))
	assert.Equal(t, "test", m.GetHeader(types.NeptuneRealTopic))
	assert.NotEmpty(t, m.GetHeader(types.DelayKeyDelayTimestamp))
}
