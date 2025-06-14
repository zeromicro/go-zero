package sarama

import (
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func convertToSaramaHeaders(headers []types.Header) []sarama.RecordHeader {
	res := make([]sarama.RecordHeader, 0, len(headers))
	for _, h := range headers {
		res = append(res, sarama.RecordHeader{
			Key:   []byte(h.Key),
			Value: h.Value,
		})
	}
	return res
}

func convertToHeaders(headers []*sarama.RecordHeader) []types.Header {
	res := make([]types.Header, 0, len(headers))
	for _, h := range headers {
		res = append(res, types.Header{
			Key:   string(h.Key),
			Value: h.Value,
		})
	}
	return res
}

func convertToHeaders2(headers []sarama.RecordHeader) []types.Header {
	res := make([]types.Header, 0, len(headers))
	for _, h := range headers {
		res = append(res, types.Header{
			Key:   string(h.Key),
			Value: h.Value,
		})
	}
	return res
}

func message2ProducerMessage(msg *types.Message) *sarama.ProducerMessage {
	if msg == nil {
		return nil
	}

	smsg := &sarama.ProducerMessage{
		Partition: int32(msg.Partition),
		Topic:     msg.Topic,
		Key:       sarama.ByteEncoder(msg.Key),
		Value:     sarama.ByteEncoder(msg.Value),
		Headers:   convertToSaramaHeaders(msg.Headers),
		Timestamp: msg.Time,
	}
	setCustomMetadata(smsg, msg.Metadata)

	return smsg
}

func producerMessage2Message(smsg *sarama.ProducerMessage) *types.Message {
	if smsg == nil {
		return nil
	}

	return &types.Message{
		Topic:     smsg.Topic,
		Partition: int(smsg.Partition),
		Offset:    smsg.Offset,
		Key:       smsg.Key.(sarama.ByteEncoder), // raw sarama message not exposed to user, so this assertion is safe
		Value:     smsg.Value.(sarama.ByteEncoder),
		Headers:   convertToHeaders2(smsg.Headers),
		Time:      smsg.Timestamp,
		Metadata:  getCustomMetadata(smsg),
	}
}

func consumerMessage2Message(cmsg *sarama.ConsumerMessage) *types.Message {
	if cmsg == nil {
		return nil
	}
	return &types.Message{
		Topic:     cmsg.Topic,
		Partition: int(cmsg.Partition),
		Offset:    cmsg.Offset,
		Key:       cmsg.Key,
		Value:     cmsg.Value,
		Headers:   convertToHeaders(cmsg.Headers),
		Time:      cmsg.Timestamp, // todo: Timestamp or BlockTimestamp
	}
}
