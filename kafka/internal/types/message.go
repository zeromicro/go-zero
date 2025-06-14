package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/utils"
)

const (
	ContentTypeKey     = "content-type"
	OriginAppNameKey   = "x-origin-app-name"
	ClientIDKey        = "x-clientid"
	CreateTimestampKey = "x-create-timestamp"

	ContentTypeJSON = "application/json"
	KafkaBroker     = "messaging.broker"
	KafkaTopic      = "messaging.destination"

	DelayKafkaMsgId        = "KAFKA_MSG_ID"
	NeptuneRealTopic       = "NEPTUNE_REAL_TOPIC"
	NeptuneTopic           = "NEPTUNE_DELAY_XXX"
	DelayKeyDelayTimestamp = "NEPTUNE_DELAY_TIMESTAMP"
	DelayTimeout           = 90 * 24 * 3600
)

type (
	// Message is a data structure representing kafka messages.
	Message struct {
		// Topic indicates which topic this message was consumed from via Reader.
		//
		// When being used with Writer, this can be used to configure the topic if
		// not already specified on the writer itself.
		Topic string

		// Partition is read-only and MUST NOT be set when writing messages
		Partition     int
		Offset        int64
		HighWaterMark int64
		Key           []byte
		Value         []byte
		Headers       []Header

		// If not set at the creation, Time will be automatically set when
		// writing the message.
		Time time.Time

		// This field is used to hold arbitrary data you wish to include so it
		// will be available when receiving on the Successes and Errors channels.
		// Sarama completely ignores this field and is only to be used for
		// pass-through data.
		// Only for producer message.
		Metadata any
	}

	// Header represents a single entry in a list of record headers.
	Header struct {
		Key   string
		Value []byte
	}
)

func (m *Message) GetHeader(key string) string {
	for _, h := range m.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (m *Message) SetHeader(key, val string) {
	// Ensure uniqueness of keys
	for i := 0; i < len(m.Headers); i++ {
		if m.Headers[i].Key == key {
			m.Headers = append(m.Headers[:i], m.Headers[i+1:]...)
			i--
		}
	}
	m.Headers = append(m.Headers, Header{
		Key:   key,
		Value: []byte(val),
	})
}

func (m *Message) BuildDelayMessage(delaySeconds int64) {
	m.SetHeader(DelayKafkaMsgId, createUniqID())
	m.SetHeader(NeptuneRealTopic, m.Topic)
	m.SetHeader(DelayKeyDelayTimestamp, fmt.Sprint(time.Now().Unix()+delaySeconds))
	m.Topic = NeptuneTopic
}

func createUniqID() string {
	str := utils.NewUuid()
	uuid := strings.Replace(str, "-", "", -1)
	return uuid
}
