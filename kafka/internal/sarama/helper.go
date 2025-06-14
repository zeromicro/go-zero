package sarama

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/tools/tls"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

const availabilityZoneId = "MY_AVAILABLE_ZONE_ID"
const appName = "MY_PROJECT_NAME"

// metadata is metadata type for byone itself.
type metadata struct {
	// internal is internal metadata for byone.
	internal any
	// custom is custom metadata for user.
	custom any
}

func toSaramaConfig(c types.ClientConfig) (*sarama.Config, error) {
	if c.MaxRequestSize > sarama.MaxRequestSize {
		sarama.MaxRequestSize = c.MaxRequestSize
	}
	if c.MaxResponseSize > sarama.MaxResponseSize {
		sarama.MaxResponseSize = c.MaxResponseSize
	}
	sc := sarama.NewConfig()

	if len(c.Brokers) == 0 {
		return nil, errors.New("empty brokers")
	}

	// consumer return errors
	sc.Consumer.Return.Errors = true

	sc.Metadata.AllowAutoTopicCreation = c.AllowAutoTopicCreation

	if c.AuthType == types.PasswordAuthType {
		if c.SaslUsername == "" || c.SaslPassword == "" {
			return nil, errors.New("username and password are required when using password auth type")
		}
		sc.Net.SASL.Enable = true
		sc.Net.SASL.User = c.SaslUsername
		sc.Net.SASL.Password = c.SaslPassword

		switch c.SaslMechanism {
		case sarama.SASLTypeSCRAMSHA256:
			sc.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
			}
			sc.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case sarama.SASLTypeSCRAMSHA512:
			sc.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
			}
			sc.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			sc.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
	}

	if c.TLSEnabled {
		sc.Net.TLS.Enable = true

		tlsConfig, err := tls.NewConfig(c.TLSClientCert, c.TLSClientKey)
		if err != nil {
			return nil, err
		}
		if c.TLSRootCACerts != "" {
			rootCAsBytes, err := os.ReadFile(c.TLSRootCACerts)
			if err != nil {
				return nil, fmt.Errorf("failed to read root CA certificates: %v", err)
			}
			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(rootCAsBytes) {
				return nil, fmt.Errorf("failed to load root CA certificates from file: %s", c.TLSRootCACerts)
			}
			// Use specific root CA set vs the host's set
			tlsConfig.RootCAs = certPool
		}
		sc.Net.TLS.Config = tlsConfig
	}

	if len(c.Version) > 0 {
		v, err := sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			return nil, errors.New("invalid kafka version")
		}
		if !v.IsAtLeast(sarama.V2_4_0_0) {
			sc.Version = sarama.V2_4_0_0
			logx.Infof("kafka version IsAtLeast: V2_4_0_0")
		} else {
			sc.Version = v
		}
	} else {
		sc.Version = sarama.V2_4_0_0
	}

	if c.AzEnabled {
		sc.Metadata.RefreshFrequency = 10 * time.Second
		sc.RackID = os.Getenv(availabilityZoneId)
		logx.Infof("consumer az enable with rackId: %s", sc.RackID)
	}

	sc.ChannelBufferSize = c.ChannelBufferSize

	if c.ClientId != "" {
		sc.ClientID = c.ClientId
	} else {
		sc.ClientID = os.Getenv(appName)
	}
	if sc.ClientID == "" {
		sc.ClientID = "sarama"
	}
	return sc, nil
}

func requiredAcksFromString(s string) (sarama.RequiredAcks, error) {
	switch s {
	case types.RequiredAcksNone:
		return sarama.NoResponse, nil
	case types.RequiredAcksOne:
		return sarama.WaitForLocal, nil
	case types.RequiredAcksAll:
		return sarama.WaitForAll, nil
	default:
		return 0, errors.New("invalid RequiredAcks")
	}
}

func parseBalanceStrategy(value string) (sarama.BalanceStrategy, error) {
	switch value {
	case types.RangeBalanceStrategyName:
		return sarama.BalanceStrategyRange, nil
	case types.RoundRobinBalanceStrategyName:
		return sarama.BalanceStrategyRoundRobin, nil
	case types.StickyBalanceStrategyName:
		return sarama.BalanceStrategySticky, nil
	default:
		return nil, fmt.Errorf("kafka error: invalid BalanceStrategy: %s", value)
	}
}

func compressionFromString(s string) (sarama.CompressionCodec, error) {
	switch s {
	case types.CompressionNone:
		return sarama.CompressionNone, nil
	case types.CompressionGZIP:
		return sarama.CompressionGZIP, nil
	case types.CompressionSnappy:
		return sarama.CompressionSnappy, nil
	case types.CompressionLZ4:
		return sarama.CompressionLZ4, nil
	case types.CompressionZSTD:
		return sarama.CompressionZSTD, nil
	default:
		return 0, fmt.Errorf("invalid compression: %s", s)
	}
}

func partitionerFromString(s string) sarama.PartitionerConstructor {
	switch s {
	case types.HashPartitioner:
		return sarama.NewHashPartitioner
	case types.RandomPartitioner:
		return sarama.NewRandomPartitioner
	case types.ManualPartitioner:
		return sarama.NewManualPartitioner
	case types.RoundRobinPartitioner:
		return sarama.NewRoundRobinPartitioner
	default:
		return sarama.NewHashPartitioner
	}
}

func getCustomMetadata(m *sarama.ProducerMessage) any {
	if m.Metadata == nil {
		return nil
	}
	return m.Metadata.(*metadata).custom
}

func getInternalMetadata(m *sarama.ProducerMessage) any {
	if m.Metadata == nil {
		return nil
	}
	return m.Metadata.(*metadata).internal
}

func setCustomMetadata(m *sarama.ProducerMessage, meta any) {
	if m.Metadata == nil {
		m.Metadata = &metadata{
			custom: meta,
		}
	} else {
		m.Metadata.(*metadata).custom = meta
	}
}

func setInternalMetadata(m *sarama.ProducerMessage, meta any) {
	if m.Metadata == nil {
		m.Metadata = &metadata{
			internal: meta,
		}
	} else {
		m.Metadata.(*metadata).internal = meta
	}
}
