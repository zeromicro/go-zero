package types

import (
	"errors"
	"sort"
	"strings"
	"time"
)

const (
	PasswordAuthType = "password"
	NoAuthType       = "none"

	OffsetOldest = "oldest"
	OffsetNewest = "newest"

	SaslMechanismPlain  = "PLAIN"
	SaslMechanismSha256 = "SCRAM-SHA-256"
	SaslMechanismSha512 = "SCRAM-SHA-512"

	RequiredAcksNone = "none"
	RequiredAcksOne  = "one"
	RequiredAcksAll  = "all"

	// RangeBalanceStrategyName identifies strategies that use the range partition assignment strategy
	RangeBalanceStrategyName = "range"
	// RoundRobinBalanceStrategyName identifies strategies that use the round-robin partition assignment strategy
	RoundRobinBalanceStrategyName = "roundrobin"
	// StickyBalanceStrategyName identifies strategies that use the sticky-partition assignment strategy
	StickyBalanceStrategyName = "sticky"

	// CompressionNone no compression
	CompressionNone = ""
	// CompressionGZIP compression using GZIP
	CompressionGZIP = "gzip"
	// CompressionSnappy compression using snappy
	CompressionSnappy = "snappy"
	// CompressionLZ4 compression using LZ4
	CompressionLZ4 = "lz4"
	// CompressionZSTD compression using ZSTD
	CompressionZSTD = "zstd"

	HashPartitioner       = "hash"
	RandomPartitioner     = "random"
	ManualPartitioner     = "manual"
	RoundRobinPartitioner = "roundrobin"
)

type (
	// ClientConfig 老版本的 kafka client 相关配置,
	// 内嵌在 ProducerConfig, ConsumerConfig, ConsumerGroupConfig 中.
	ClientConfig struct {
		ResourceName           string   `json:",optional"`                                                // include Brokers/AuthType/SaslUsername/SaslPassword/SaslMechanism
		Brokers                []string `json:",optional=!ResourceName"`                                  // kafka broker addresses
		AuthType               string   `json:",default=none,options=none|password"`                      // auth type, password / none. default: none
		SaslUsername           string   `json:",optional"`                                                // required when authType is password
		SaslPassword           string   `json:",optional"`                                                // required when authType is password
		SaslMechanism          string   `json:",default=PLAIN,options=PLAIN|SCRAM-SHA-256|SCRAM-SHA-512"` // kafka SASL.Mechanism, PLAIN|SCRAM-SHA-256|SCRAM-SHA-512, default PLAIN
		Version                string   `json:",default=2.4.0"`                                           // kafka version, default 2.4.0
		AllowAutoTopicCreation bool     `json:",optional"`                                                // If true, the broker may auto-create topics which do not exist, default is false

		TLSEnabled     bool   `json:",default=false"` // enable TLS, default false
		TLSClientCert  string `json:",optional"`
		TLSClientKey   string `json:",optional"`
		TLSRootCACerts string `json:",optional"`

		AzEnabled         bool  `json:",default=true"`      // enable AZ consume, default true, set rack.id from env
		ChannelBufferSize int   `json:",default=10240"`     // The number of events to buffer in internal and external channels
		MaxRequestSize    int32 `json:",default=524288000"` // is the maximum size (in bytes) of any request that Sarama will attempt to send.
		MaxResponseSize   int32 `json:",default=524288000"` // is the maximum size (in bytes) of any response that Sarama will attempt to parse

		ClientId string `json:",optional"`
	}

	RetryConfig struct {
		MaxRetries int64 `json:",optional"` // max retry times default 0
	}

	// SharedProducerConfig 是共享 kafka producer 相关参数配置.
	SharedProducerConfig struct {
		MaxMessageBytes int    `json:",default=524288000"`
		AppName         string `json:",optional"`                                           // todo: load from env?
		RequiredAcks    string `json:",default=all,options=none|one|all"`                   // none=0, one=1, all=-1, default is all
		Compression     string `json:",default=lz4"`                                        // compression type, can be one of gzip|snappy|lz4|zstd, default is none
		Idempotent      bool   `json:",default=true"`                                       // default true, the producer will ensure that exactly one copy of each message is write
		Partitioner     string `json:",default=hash,options=hash|random|manual|roundrobin"` // partitioner type, can be one of hash|random|manual|roundrobin, default is hash
		EnableRecovery  bool   `json:",default=true"`                                       //callback handler panic is recovery when EnableRecovery enabled, default true

		Flush struct {
			Bytes       int           `json:",default=65536"`  // The best-effort number of bytes needed to trigger a flush. Use the global sarama.MaxRequestSize to set a hard upper limit.
			Messages    int           `json:",default=0"`      // The best-effort number of messages needed to trigger a flush. Use`MaxMessages` to set a hard upper limit.
			Frequency   time.Duration `json:",default=1ms"`    // The best-effort frequency of flushes. Equivalent to `queue.buffering.max.ms` setting of JVM producer.
			MaxMessages int           `json:",default=100000"` // The maximum number of messages the producer will send in a single broker request. Defaults to 0 for unlimited. Similar to `queue.buffering.max.messages` in the JVM producer.
		}
	}

	// ProducerConfig 是老版本 API kafka producer 初始化配置.
	ProducerConfig struct {
		Client ClientConfig
		// Deprecated: please set Topic in ProducerMessage directly
		Topic string `json:",optional"` // default topic
		SharedProducerConfig
	}

	// SharedConsumerConfig 是共享 kafka consumer group 相关参数配置.
	SharedConsumerConfig struct {
		// kafka fetch.min.bytes, default is 1,
		// fetch requests are answered as soon as a single byte of data is available
		FetchMinBytes int `json:",default=1"`
		// kafka fetch.max.bytes, default is 10e6(10MB)
		FetchMaxBytes int `json:",default=10000000"`
		// kafka fetch.wait.max.ms, default is 500ms
		FetchMaxWaitTime time.Duration `json:",default=50ms"`
		// timeout for consuming a message, default is 0, no timeout
		ConsumeTimeout time.Duration `json:",optional"`
		// kafka consumer group rebalance strategy, can be one of range|roundrobin|sticky,
		// default is range.
		BalanceStrategy string `json:",default=range,options=range|roundrobin|sticky"`
		// kafka MaxProcessingTime, default is 100ms
		MaxProcessingTime time.Duration `json:",default=100ms"`
	}

	// GroupConfig 是共享 client 创建 kafka consumer 相关参数.
	GroupConfig struct {
		Topic           string
		GroupID         string
		InitialOffset   string      `json:",default=newest,options=oldest|newest"`
		RetryConfig     RetryConfig `json:",optional"`
		LogOffsets      bool        `json:",optional"`
		EnableRecovery  bool        `json:",default=true"` //callback handler panic is recovery when EnableRecovery enabled, default true
		AutoCommit      bool        `json:",default=true"`
		DisableAutoMark bool        `json:",default=false"`
	}

	// ConsumerGroupConfig 是老版本 API kafka consumer group 初始化配置.
	ConsumerGroupConfig struct {
		Client ClientConfig
		SharedConsumerConfig
		GroupConfig
	}

	// ConsumerConfig 是老版本 API 创建 kafka consumer 初始化配置.
	ConsumerConfig struct {
		Client ClientConfig
		SharedConsumerConfig
	}

	UniversalClientConfig struct {
		Client   ClientConfig
		Producer SharedProducerConfig
		Consumer SharedConsumerConfig
	}
)

// GetClientName combine brokers to stable string which is useful for metric label.
func (c *ClientConfig) GetClientName() string {
	sort.Strings(c.Brokers)
	return strings.Join(c.Brokers, ",")
}

func (c *ClientConfig) Validate() error {
	if len(c.Brokers) == 0 {
		return errors.New("empty brokers")
	}

	if c.AuthType == PasswordAuthType {
		if c.SaslUsername == "" || c.SaslPassword == "" {
			return errors.New("username and password are required when using password auth type")
		}
	}

	return nil
}
