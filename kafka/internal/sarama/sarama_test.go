package sarama

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func Test_convertToSaramaHeaders(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []sarama.RecordHeader{}, convertToSaramaHeaders(nil))
	})

	t.Run("should works", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		inputs := make([]types.Header, 0, 100)
		for i := 0; i < 100; i++ {
			val := make([]byte, 100)
			_, _ = rand.Read(val)
			inputs = append(inputs, types.Header{
				Key:   strconv.Itoa(i),
				Value: val,
			})
		}

		res := convertToSaramaHeaders(inputs)
		for i := 0; i < 100; i++ {
			actual := res[i]
			origin := inputs[i]
			assert.Equal(t, origin.Key, string(actual.Key))
			assert.True(t, bytes.Equal(origin.Value, actual.Value))
		}
	})
}

func Test_convertToHeaders(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []types.Header{}, convertToHeaders(nil))
	})

	t.Run("should works", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		inputs := make([]*sarama.RecordHeader, 0, 100)
		for i := 0; i < 100; i++ {
			val := make([]byte, 100)
			_, _ = rand.Read(val)
			inputs = append(inputs, &sarama.RecordHeader{
				Key:   []byte(strconv.Itoa(i)),
				Value: val,
			})
		}

		res := convertToHeaders(inputs)
		for i := 0; i < 100; i++ {
			actual := res[i]
			origin := inputs[i]
			assert.Equal(t, string(origin.Key), actual.Key)
			assert.True(t, bytes.Equal(origin.Value, actual.Value))
		}
	})
}

func Test_convertToHeaders2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []types.Header{}, convertToHeaders2(nil))
	})

	t.Run("should works", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		inputs := make([]sarama.RecordHeader, 0, 100)
		for i := 0; i < 100; i++ {
			val := make([]byte, 100)
			_, _ = rand.Read(val)
			inputs = append(inputs, sarama.RecordHeader{
				Key:   []byte(strconv.Itoa(i)),
				Value: val,
			})
		}

		res := convertToHeaders2(inputs)
		for i := 0; i < 100; i++ {
			actual := res[i]
			origin := inputs[i]
			assert.Equal(t, string(origin.Key), actual.Key)
			assert.True(t, bytes.Equal(origin.Value, actual.Value))
		}
	})
}

func Test_message2ProducerMessage(t *testing.T) {
	assert.Nil(t, message2ProducerMessage(nil))

	t.Run("should omit partition and offset", func(t *testing.T) {
		testByte := []byte("test")
		msg := &types.Message{
			Topic:         "aaa",
			Partition:     10,
			Offset:        10,
			HighWaterMark: 10,
			Key:           testByte,
			Value:         testByte,
			Headers:       []types.Header{{Key: "test", Value: testByte}},
			Time:          time.Now(),
		}

		resp := message2ProducerMessage(msg)
		assert.Equal(t, msg.Topic, resp.Topic)
		assert.Equal(t, int32(10), resp.Partition)
		assert.Equal(t, int64(0), resp.Offset)
		assert.True(t, bytes.Equal(msg.Key, byteEncoderToByte(t, resp.Key)))
		assert.True(t, bytes.Equal(msg.Value, byteEncoderToByte(t, resp.Value)))
		assert.True(t, resp.Timestamp.Equal(msg.Time))
		assert.Equal(t, 1, len(resp.Headers))
		assert.Nil(t, getCustomMetadata(resp))
	})

	t.Run("should passthrough user's metadata", func(t *testing.T) {
		testByte := []byte("test")
		type abc struct {
			foo string
		}
		ref := &abc{
			foo: "xxx",
		}

		msg := &types.Message{
			Topic:    "aaa",
			Value:    testByte,
			Metadata: ref,
		}
		resp := message2ProducerMessage(msg)
		assert.Equal(t, ref, getCustomMetadata(resp))
	})
}

func byteEncoderToByte(t *testing.T, encoder sarama.Encoder) []byte {
	resp, err := encoder.Encode()
	assert.Nil(t, err)
	return resp
}

func Test_consumerMessage2Message(t *testing.T) {
	assert.Nil(t, consumerMessage2Message(nil))

	t.Run("should works", func(t *testing.T) {
		testByte := []byte("test")
		msg := &sarama.ConsumerMessage{
			Headers:        []*sarama.RecordHeader{{Key: testByte, Value: testByte}},
			Timestamp:      time.Now(),
			BlockTimestamp: time.Time{},
			Key:            testByte,
			Value:          testByte,
			Topic:          "test",
			Partition:      100,
			Offset:         1000,
		}

		resp := consumerMessage2Message(msg)
		assert.Equal(t, msg.Topic, resp.Topic)
		assert.Equal(t, 100, resp.Partition)
		assert.Equal(t, int64(1000), resp.Offset)
		assert.True(t, bytes.Equal(msg.Key, resp.Key))
		assert.True(t, bytes.Equal(msg.Value, resp.Value))
		assert.True(t, resp.Time.Equal(msg.Timestamp))
		assert.Equal(t, 1, len(resp.Headers))
		header := resp.Headers[0]
		assert.Equal(t, string(testByte), header.Key)
		assert.True(t, bytes.Equal(testByte, header.Value))
	})
}

func Test_producerMessage2Message(t *testing.T) {
	assert.Nil(t, producerMessage2Message(nil))

	t.Run("should works", func(t *testing.T) {
		testByte := []byte("test")
		msg := &sarama.ProducerMessage{
			Headers:   []sarama.RecordHeader{{Key: testByte, Value: testByte}},
			Timestamp: time.Now(),
			Key:       sarama.ByteEncoder(testByte),
			Value:     sarama.ByteEncoder(testByte),
			Topic:     "test",
			Partition: 100,
			Offset:    1000,
		}

		resp := producerMessage2Message(msg)
		assert.Equal(t, msg.Topic, resp.Topic)
		assert.Equal(t, 100, resp.Partition)
		assert.Equal(t, int64(1000), resp.Offset)
		assert.True(t, bytes.Equal(msg.Key.(sarama.ByteEncoder), resp.Key))
		assert.True(t, bytes.Equal(msg.Value.(sarama.ByteEncoder), resp.Value))
		assert.True(t, resp.Time.Equal(msg.Timestamp))
		assert.Equal(t, 1, len(resp.Headers))
		header := resp.Headers[0]
		assert.Equal(t, string(testByte), header.Key)
		assert.True(t, bytes.Equal(testByte, header.Value))
		assert.Nil(t, resp.Metadata)
	})

	t.Run("should pass through user's metadata", func(t *testing.T) {
		testByte := []byte("test")
		type abc struct {
			foo string
		}
		ref := &abc{
			foo: "xxx",
		}

		msg := &sarama.ProducerMessage{
			Key:   sarama.ByteEncoder(testByte),
			Value: sarama.ByteEncoder(testByte),
			Topic: "test",
		}
		setCustomMetadata(msg, ref)

		resp := producerMessage2Message(msg)
		assert.Equal(t, ref, resp.Metadata)
	})
}
