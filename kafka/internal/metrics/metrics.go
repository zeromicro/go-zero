package metrics

import "github.com/zeromicro/go-zero/core/metric"

const (
	CodeOK    = "OK"
	CodeError = "Error"

	ProducerTypeSync  = "sync"
	ProducerTypeASync = "async"

	serverNamespace = "kafka"
)

var (
	KafkaPublishCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Name:      "publish_total",
		Help:      "kafka producer publish messages total.",
		Labels:    []string{"producer_type", "brokers", "topic", "code"},
	})

	KafkaPublishHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Name:      "publish_seconds",
		Help:      "kafka producer publish messages cost(s)",
		Labels:    []string{"producer_type", "brokers", "topic"},
		Buckets:   []float64{.025, .05, .1, .25, .5, 1, 2.5, 5},
	})

	KafkaPublishGauge = newCollector[ProducerKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "publish_max_seconds",
		Help:      "kafka producer publish messages max cost(s)",
		Labels:    []string{"producer_type", "brokers", "topic"},
	})

	KafkaConsumerGroupCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_group_handle_total",
		Help:      "kafka consumergroup consume messages total.",
		Labels:    []string{"brokers", "group_id", "topic", "code"},
	})

	KafkaConsumerGroupHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_group_handle_seconds",
		Help:      "kafka consumergroup handler handle cost(s)",
		Labels:    []string{"brokers", "group_id", "topic"},
		Buckets:   []float64{.025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
	})

	KafkaConsumerGroupGauge = newCollector[ConsumerGroupKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_group_handle_max_seconds",
		Help:      "kafka consumergroup handler max handle cost(s)",
		Labels:    []string{"brokers", "group_id", "topic"},
	})

	KafkaConsumerCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_handle_total",
		Help:      "kafka consumer consume messages total.",
		Labels:    []string{"brokers", "topic", "partition"},
	})

	KafkaConsumerHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_handle_seconds",
		Help:      "kafka consumer handler handle cost(s)",
		Labels:    []string{"brokers", "topic", "partition"},
		Buckets:   []float64{.025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
	})

	KafkaConsumerGauge = newCollector[ConsumerKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "consumer_handle_max_seconds",
		Help:      "kafka consumer handler max handle cost(s)",
		Labels:    []string{"brokers", "topic", "partition"},
	})

	KafkaPayloadSizeHistogram = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Name:      "payload_size",
		Help:      "kafka producer publish payload size(byte)",
		Labels:    []string{"topic", "partition"},
		Buckets:   []float64{1 * 1024 * 1024, 10 * 1024 * 1024, 50 * 1024 * 1024, 100 * 1024 * 1024, 200 * 1024 * 1024},
	})
)
