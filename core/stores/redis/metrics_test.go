package redis

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	red "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/internal/devserver"
)

func TestRedisMetric(t *testing.T) {
	cfg := devserver.Config{}
	_ = conf.FillDefault(&cfg)
	server := devserver.NewServer(cfg)
	server.StartAsync(cfg)
	time.Sleep(time.Second)

	metricReqDur.Observe(8, "test-cmd")
	metricReqErr.Inc("test-cmd", "internal-error")
	metricSlowCount.Inc("test-cmd")

	url := "http://127.0.0.1:6060/metrics"
	resp, err := http.Get(url)
	assert.Nil(t, err)
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	content := string(s)
	assert.Contains(t, content, "redis_client_requests_duration_ms_sum{command=\"test-cmd\"} 8\n")
	assert.Contains(t, content, "redis_client_requests_duration_ms_count{command=\"test-cmd\"} 1\n")
	assert.Contains(t, content, "redis_client_requests_error_total{command=\"test-cmd\",error=\"internal-error\"} 1\n")
	assert.Contains(t, content, "redis_client_requests_slow_total{command=\"test-cmd\"} 1\n")
}

func Test_newCollector(t *testing.T) {
	prometheus.Unregister(connCollector)
	c := newCollector()
	c.registerClient(&statGetter{
		clientType: "node",
		key:        "test1",
		poolSize:   10,
		poolStats: func() *red.PoolStats {
			return &red.PoolStats{
				Hits:       10000,
				Misses:     10,
				Timeouts:   5,
				TotalConns: 100,
				IdleConns:  20,
				StaleConns: 1,
			}
		},
	})
	c.registerClient(&statGetter{
		clientType: "node",
		key:        "test2",
		poolSize:   11,
		poolStats: func() *red.PoolStats {
			return &red.PoolStats{
				Hits:       10001,
				Misses:     11,
				Timeouts:   6,
				TotalConns: 101,
				IdleConns:  21,
				StaleConns: 2,
			}
		},
	})
	c.registerClient(&statGetter{
		clientType: "cluster",
		key:        "test3",
		poolSize:   5,
		poolStats: func() *red.PoolStats {
			return &red.PoolStats{
				Hits:       20000,
				Misses:     20,
				Timeouts:   10,
				TotalConns: 200,
				IdleConns:  40,
				StaleConns: 2,
			}
		},
	})
	val := `
		# HELP redis_client_pool_conn_idle_current Current number of idle connections in the pool
		# TYPE redis_client_pool_conn_idle_current gauge
		redis_client_pool_conn_idle_current{client_type="cluster",key="test3"} 40
		redis_client_pool_conn_idle_current{client_type="node",key="test1"} 20
		redis_client_pool_conn_idle_current{client_type="node",key="test2"} 21
	 	# HELP redis_client_pool_conn_max Max number of connections in the pool
	 	# TYPE redis_client_pool_conn_max counter
		redis_client_pool_conn_max{client_type="cluster",key="test3"} 5
	 	redis_client_pool_conn_max{client_type="node",key="test1"} 10
		redis_client_pool_conn_max{client_type="node",key="test2"} 11
	 	# HELP redis_client_pool_conn_stale_total Number of times a connection was removed from the pool because it was stale
	 	# TYPE redis_client_pool_conn_stale_total counter
		redis_client_pool_conn_stale_total{client_type="cluster",key="test3"} 2
	 	redis_client_pool_conn_stale_total{client_type="node",key="test1"} 1
		redis_client_pool_conn_stale_total{client_type="node",key="test2"} 2
	 	# HELP redis_client_pool_conn_total_current Current number of connections in the pool
	 	# TYPE redis_client_pool_conn_total_current gauge
		redis_client_pool_conn_total_current{client_type="cluster",key="test3"} 200
	 	redis_client_pool_conn_total_current{client_type="node",key="test1"} 100
		redis_client_pool_conn_total_current{client_type="node",key="test2"} 101
	 	# HELP redis_client_pool_hit_total Number of times a connection was found in the pool
	 	# TYPE redis_client_pool_hit_total counter
		redis_client_pool_hit_total{client_type="cluster",key="test3"} 20000
	 	redis_client_pool_hit_total{client_type="node",key="test1"} 10000
		redis_client_pool_hit_total{client_type="node",key="test2"} 10001
	 	# HELP redis_client_pool_miss_total Number of times a connection was not found in the pool
	 	# TYPE redis_client_pool_miss_total counter
		redis_client_pool_miss_total{client_type="cluster",key="test3"} 20
	 	redis_client_pool_miss_total{client_type="node",key="test1"} 10
		redis_client_pool_miss_total{client_type="node",key="test2"} 11
	 	# HELP redis_client_pool_timeout_total Number of times a timeout occurred when looking for a connection in the pool
	 	# TYPE redis_client_pool_timeout_total counter
		redis_client_pool_timeout_total{client_type="cluster",key="test3"} 10
	 	redis_client_pool_timeout_total{client_type="node",key="test1"} 5
		redis_client_pool_timeout_total{client_type="node",key="test2"} 6
`

	err := testutil.CollectAndCompare(c, strings.NewReader(val))
	assert.NoError(t, err)
}
