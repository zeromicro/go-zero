package sqlx

import (
	"database/sql"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/metric"
)

const namespace = "sql_client"

var (
	metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "mysql client requests duration(ms).",
		Labels:    []string{"command"},
		Buckets:   []float64{0.25, 0.5, 1, 1.5, 2, 3, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	})
	metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "mysql client requests error count.",
		Labels:    []string{"command", "error"},
	})
	metricSlowCount = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "slow_total",
		Help:      "mysql client requests slow count.",
		Labels:    []string{"command"},
	})

	connLabels                         = []string{"db_name", "hash"}
	connCollector                      = newCollector()
	_             prometheus.Collector = (*collector)(nil)
)

type (
	statGetter struct {
		host      string
		dbName    string
		hash      string
		poolStats func() sql.DBStats
	}

	// collector collects statistics from a redis client.
	// It implements the prometheus.Collector interface.
	collector struct {
		maxOpenConnections *prometheus.Desc

		openConnections  *prometheus.Desc
		inUseConnections *prometheus.Desc
		idleConnections  *prometheus.Desc

		waitCount         *prometheus.Desc
		waitDuration      *prometheus.Desc
		maxIdleClosed     *prometheus.Desc
		maxIdleTimeClosed *prometheus.Desc
		maxLifetimeClosed *prometheus.Desc

		clients []*statGetter
		lock    sync.Mutex
	}
)

func newCollector() *collector {
	c := &collector{
		maxOpenConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_open_connections"),
			"Maximum number of open connections to the database.",
			connLabels, nil,
		),
		openConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "open_connections"),
			"The number of established connections both in use and idle.",
			connLabels, nil,
		),
		inUseConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "in_use_connections"),
			"The number of connections currently in use.",
			connLabels, nil,
		),
		idleConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "idle_connections"),
			"The number of idle connections.",
			connLabels, nil,
		),
		waitCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "wait_count_total"),
			"The total number of connections waited for.",
			connLabels, nil,
		),
		waitDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "wait_duration_seconds_total"),
			"The total time blocked waiting for a new connection.",
			connLabels, nil,
		),
		maxIdleClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_idle_closed_total"),
			"The total number of connections closed due to SetMaxIdleConns.",
			connLabels, nil,
		),
		maxIdleTimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_idle_time_closed_total"),
			"The total number of connections closed due to SetConnMaxIdleTime.",
			connLabels, nil,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_lifetime_closed_total"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			connLabels, nil,
		),
	}

	prometheus.MustRegister(c)

	return c
}

// Describe implements the prometheus.Collector interface.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnections
	ch <- c.openConnections
	ch <- c.inUseConnections
	ch <- c.idleConnections
	ch <- c.waitCount
	ch <- c.waitDuration
	ch <- c.maxIdleClosed
	ch <- c.maxLifetimeClosed
	ch <- c.maxIdleTimeClosed
}

// Collect implements the prometheus.Collector interface.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, client := range c.clients {
		dbName, hash := client.dbName, client.hash
		stats := client.poolStats()
		ch <- prometheus.MustNewConstMetric(c.maxOpenConnections, prometheus.GaugeValue,
			float64(stats.MaxOpenConnections), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.openConnections, prometheus.GaugeValue,
			float64(stats.OpenConnections), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.inUseConnections, prometheus.GaugeValue,
			float64(stats.InUse), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.idleConnections, prometheus.GaugeValue,
			float64(stats.Idle), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.waitCount, prometheus.CounterValue,
			float64(stats.WaitCount), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.waitDuration, prometheus.CounterValue,
			stats.WaitDuration.Seconds(), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.CounterValue,
			float64(stats.MaxIdleClosed), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.CounterValue,
			float64(stats.MaxLifetimeClosed), dbName, hash)
		ch <- prometheus.MustNewConstMetric(c.maxIdleTimeClosed, prometheus.CounterValue,
			float64(stats.MaxIdleTimeClosed), dbName, hash)
	}
}

func (c *collector) registerClient(client *statGetter) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.clients = append(c.clients, client)
}
