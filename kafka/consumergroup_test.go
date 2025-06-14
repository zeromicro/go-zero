package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func TestNewConsumerGroup(t *testing.T) {
	config := ConsumerGroupConfig{Client: ClientConfig{
		ResourceName: "kafka-erms-test-1",
	},
		GroupConfig: GroupConfig{
			EnableRecovery: false,
		}}
	_, err := NewConsumerGroup(config, nil)
	assert.Error(t, err)
}

type consumerGroup0 struct {
}

func (c consumerGroup0) Start() {
}

func (c consumerGroup0) Close() error {
	return nil
}

func (c consumerGroup0) MarkMessage(message *types.Message) {
}

func (c consumerGroup0) Commit() {
}

func TestMarkMessage(t *testing.T) {
	logx.Disable()
	var cp = &consumerGroup{
		c: &consumerGroup0{},
	}
	cp.MarkMessage(&Message{Topic: "", Value: []byte("")})
}
