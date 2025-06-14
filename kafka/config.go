package kafka

import "github.com/zeromicro/go-zero/kafka/internal/types"

type (
	// ClientConfig is shared kafka client config.
	ClientConfig = types.ClientConfig
	// RetryConfig is retry config for consumer handler.
	RetryConfig = types.RetryConfig
	// ProducerConfig is config for kafka producer.
	ProducerConfig = types.ProducerConfig
	// ConsumerGroupConfig is config for kafka consumer group.
	ConsumerGroupConfig = types.ConsumerGroupConfig
	// ConsumerConfig is config for kafka consumer.
	ConsumerConfig = types.ConsumerConfig

	// UniversalClientConfig 是新 API 创建共享 kafka client 的初始化配置.
	// 包括 kafka client 配置, producer, consumer group 共享配置.
	UniversalClientConfig = types.UniversalClientConfig
	// GroupConfig 是共享 client 创建 kafka consumer group 的初始化配置.
	GroupConfig = types.GroupConfig
	// SharedProducerConfig 是共享 kafka producer 相关参数配置.
	SharedProducerConfig = types.SharedProducerConfig
	// SharedConsumerConfig 是共享 kafka consumer group 相关参数配置.
	SharedConsumerConfig = types.SharedConsumerConfig
)
