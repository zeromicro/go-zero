package sarama

import (
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func TestClient_GetOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().GetOffset("test", int32(0), sarama.OffsetNewest).Return(int64(10), nil)
	c := &Client{saramaClient: sc}
	offset, err := c.GetOffset("test", 0, sarama.OffsetNewest)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), offset)
}

func TestClient_Partitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().Partitions("test").Return([]int32{0, 1}, nil)
	c := &Client{saramaClient: sc}
	partitions, err := c.Partitions("test")
	assert.NoError(t, err)
	assert.Equal(t, []int32{0, 1}, partitions)
}

func TestClient_waitTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	c := &Client{saramaClient: sc}

	wg := sync.WaitGroup{}
	wg.Add(1)
	now := time.Now()
	c.waitTimeout(time.Second, func() {
		time.Sleep(2 * time.Minute)
		wg.Done()
	}, &wg)
	assert.Less(t, time.Since(now).Seconds(), float64(2))
}

func TestClient_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
	c := &Client{saramaClient: sc}
	err := c.Close()
	assert.NoError(t, err)
}

func TestNewClient(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		_, err := NewClient(types.UniversalClientConfig{})
		assert.Error(t, err)

		_, err = NewClient(types.UniversalClientConfig{
			Client: types.ClientConfig{Brokers: []string{"b1"}},
			Producer: types.SharedProducerConfig{
				RequiredAcks: types.RequiredAcksNone,
				Compression:  types.CompressionGZIP,
			},
		})
		assert.Error(t, err)

		_, err = NewClient(types.UniversalClientConfig{
			Client: types.ClientConfig{Brokers: []string{"b1"}},
			Producer: types.SharedProducerConfig{
				RequiredAcks: types.RequiredAcksNone,
				Compression:  types.CompressionGZIP,
			},
			Consumer: types.SharedConsumerConfig{
				BalanceStrategy: types.RangeBalanceStrategyName,
			},
		})
		assert.Error(t, err)
	})
}
func TestClientNewConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
	sc.EXPECT().Closed().MinTimes(1).MaxTimes(1).Return(false)
	sc.EXPECT().Config().MinTimes(1).Return(&sarama.Config{})
	c := &Client{saramaClient: sc, config: types.UniversalClientConfig{}, saramaConfig: &sarama.Config{}}
	_, err := c.NewConsumer()
	assert.NoError(t, err)
}

func TestClientNewConsumerGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
	sc.EXPECT().Closed().MinTimes(1).MaxTimes(1).Return(false)
	c := &Client{saramaClient: sc}
	_, err := c.NewConsumerGroup(types.GroupConfig{
		EnableRecovery: false,
	}, nil)
	assert.Error(t, err)
}

func runWithClient(t *testing.T, f func(t *testing.T, client *Client)) {
	logx.Disable()

	seedBroker := sarama.NewMockBroker(t, 1)
	seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetController(seedBroker.BrokerID()).
			SetBroker(seedBroker.Addr(), seedBroker.BrokerID()),
	})
	defer seedBroker.Close()

	var config types.UniversalClientConfig
	conf.FillDefault(&config)
	config.Client.Brokers = []string{seedBroker.Addr()}

	client, err := NewClient(config)
	assert.Nil(t, err)
	defer client.Close()

	f(t, client)
}
