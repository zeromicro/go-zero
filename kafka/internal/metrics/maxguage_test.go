package metrics

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/metric"
)

func TestMaxGaugeVec_Set(t *testing.T) {
	c := newCollector[ConsumerKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "TestMaxGaugeVec_Set",
		Labels:    []string{"brokers", "topic", "partition", "refer_site", "traffic_label"},
	})

	key := ConsumerKey{
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		Partition:    "partition-1",
		ReferSite:    "",
		TrafficLabel: "",
	}

	// insert
	c.Set(10, key)
	assert.Equal(t, float64(10), c.values[key])

	// ignore
	c.Set(5, key)
	assert.Equal(t, float64(10), c.values[key])

	// update
	c.Set(20, key)
	assert.Equal(t, float64(20), c.values[key])

	// clear
	c.release()
	assert.Empty(t, c.values)

	// insert
	c.Set(10, ConsumerKey{
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		Partition:    "partition-1",
		ReferSite:    "",
		TrafficLabel: "",
	})
	// insert
	c.Set(20, ConsumerKey{
		Brokers:      "127.0.0.1",
		Topic:        "topic-2",
		Partition:    "partition-2",
		ReferSite:    "",
		TrafficLabel: "",
	})
	assert.Equal(t, 2, len(c.values))
}

func TestMetricCollector_Producer(t *testing.T) {
	c := newCollector[ProducerKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "TestMetricCollector_Producer",
		Help:      "Help",
		Labels:    []string{"producer_type", "brokers", "topic", "refer_site", "traffic_label"},
	})

	c.Set(10, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		ReferSite:    "",
		TrafficLabel: "",
	})
	c.Set(20, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-2",
		ReferSite:    "",
		TrafficLabel: "",
	})
	val := `
# HELP kafka_TestMetricCollector_Producer Help
# TYPE kafka_TestMetricCollector_Producer gauge
kafka_TestMetricCollector_Producer{brokers="127.0.0.1",producer_type="sync",refer_site="",topic="topic-1",traffic_label=""} 10
kafka_TestMetricCollector_Producer{brokers="127.0.0.1",producer_type="sync",refer_site="",topic="topic-2",traffic_label=""} 20
`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(val)))

	c.Set(25, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		ReferSite:    "",
	})
	c.Set(15, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		ReferSite:    "",
	})
	c.Set(10, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-2",
		ReferSite:    "",
	})
	c.Set(30, ProducerKey{
		ProducerType: "sync",
		Brokers:      "127.0.0.1",
		Topic:        "topic-2",
		ReferSite:    "",
	})
	val = `
# HELP kafka_TestMetricCollector_Producer Help
# TYPE kafka_TestMetricCollector_Producer gauge
kafka_TestMetricCollector_Producer{brokers="127.0.0.1",producer_type="sync",refer_site="",topic="topic-1",traffic_label=""} 25
kafka_TestMetricCollector_Producer{brokers="127.0.0.1",producer_type="sync",refer_site="",topic="topic-2",traffic_label=""} 30
`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(val)))
}

func TestMetricCollector_Consumer(t *testing.T) {
	c := newCollector[ConsumerKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "TestMetricCollector_Consumer",
		Help:      "Help",
		Labels:    []string{"brokers", "topic", "partition", "refer_site", "traffic_label"},
	})

	c.Set(10, ConsumerKey{
		Brokers:      "127.0.0.1",
		Topic:        "topic-1",
		Partition:    "1",
		ReferSite:    "",
		TrafficLabel: "",
	})
	c.Set(20, ConsumerKey{
		Brokers:      "127.0.0.1",
		Topic:        "topic-2",
		Partition:    "2",
		ReferSite:    "",
		TrafficLabel: "",
	})
	val := `
# HELP kafka_TestMetricCollector_Consumer Help
# TYPE kafka_TestMetricCollector_Consumer gauge
kafka_TestMetricCollector_Consumer{brokers="127.0.0.1",partition="1",refer_site="",topic="topic-1",traffic_label=""} 10
kafka_TestMetricCollector_Consumer{brokers="127.0.0.1",partition="2",refer_site="",topic="topic-2",traffic_label=""} 20
`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(val)))
}

func TestMetricCollector_ConsumerGroup(t *testing.T) {
	c := newCollector[ConsumerGroupKey](&metric.GaugeVecOpts{
		Namespace: serverNamespace,
		Name:      "TestMetricCollector_ConsumerGroup",
		Help:      "Help",
		Labels:    []string{"brokers", "group_id", "topic", "refer_site", "traffic_label"},
	})

	c.Set(10, ConsumerGroupKey{
		Brokers:      "127.0.0.1",
		GroupId:      "group-1",
		Topic:        "topic-1",
		ReferSite:    "",
		TrafficLabel: "",
	})
	c.Set(20, ConsumerGroupKey{
		Brokers:      "127.0.0.1",
		GroupId:      "group-2",
		Topic:        "topic-2",
		ReferSite:    "",
		TrafficLabel: "",
	})
	val := `
# HELP kafka_TestMetricCollector_ConsumerGroup Help
# TYPE kafka_TestMetricCollector_ConsumerGroup gauge
kafka_TestMetricCollector_ConsumerGroup{brokers="127.0.0.1",group_id="group-1",refer_site="",topic="topic-1",traffic_label=""} 10
kafka_TestMetricCollector_ConsumerGroup{brokers="127.0.0.1",group_id="group-2",refer_site="",topic="topic-2",traffic_label=""} 20
`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(val)))
}
