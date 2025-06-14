package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/metric"
)

var _ prometheus.Collector = (*collector[ProducerKey])(nil)

type labelContainer interface {
	Labels() []string
}

type ProducerKey struct {
	ProducerType string
	Brokers      string
	Topic        string
	ReferSite    string
	TrafficLabel string
}

func (k ProducerKey) Labels() []string {
	return []string{k.ProducerType, k.Brokers, k.Topic, k.ReferSite, k.TrafficLabel}
}

type ConsumerKey struct {
	Brokers      string
	Topic        string
	Partition    string
	ReferSite    string
	TrafficLabel string
}

func (k ConsumerKey) Labels() []string {
	return []string{k.Brokers, k.Topic, k.Partition, k.ReferSite, k.TrafficLabel}
}

type ConsumerGroupKey struct {
	Brokers      string
	GroupId      string
	Topic        string
	ReferSite    string
	TrafficLabel string
}

func (k ConsumerGroupKey) Labels() []string {
	return []string{k.Brokers, k.GroupId, k.Topic, k.ReferSite, k.TrafficLabel}
}

type keyType interface {
	ProducerKey | ConsumerKey | ConsumerGroupKey
}

type collector[T keyType] struct {
	lock   sync.Mutex
	values map[T]float64
	desc   *prometheus.Desc
}

func newCollector[T keyType](cfg *metric.GaugeVecOpts) *collector[T] {
	c := &collector[T]{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(cfg.Namespace, cfg.Subsystem, cfg.Name),
			cfg.Help, cfg.Labels, cfg.ConstLabels,
		),
		values: make(map[T]float64),
	}

	prometheus.MustRegister(c)

	return c
}

func (c *collector[T]) Set(value float64, key T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	curValue, exist := c.values[key]
	if !exist || curValue < value {
		c.values[key] = value
	}
}

// Describe implements the prometheus.Collector interface.
func (c *collector[T]) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *collector[T]) release() map[T]float64 {
	newMap := make(map[T]float64)
	c.lock.Lock()
	defer c.lock.Unlock()
	oldValue := c.values
	c.values = newMap
	return oldValue
}

// Collect implements the prometheus.Collector interface.
func (c *collector[T]) Collect(ch chan<- prometheus.Metric) {
	valueMap := c.release()
	for k, v := range valueMap {
		// k.(labelContainer) is invalid code
		var a any = k
		kt := a.(labelContainer)
		labels := kt.Labels()
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, v, labels...)
	}
}
