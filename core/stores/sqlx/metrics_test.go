package sqlx

import (
	"database/sql"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/internal/devserver"
)

func TestSqlxMetric(t *testing.T) {
	cfg := devserver.Config{}
	_ = conf.FillDefault(&cfg)
	cfg.Port = 6480
	server := devserver.NewServer(cfg)
	server.StartAsync(cfg)
	time.Sleep(time.Second)

	metricReqDur.Observe(8, "test-cmd")
	metricReqErr.Inc("test-cmd", "internal-error")
	metricSlowCount.Inc("test-cmd")

	url := "http://127.0.0.1:6480/metrics"
	resp, err := http.Get(url)
	assert.Nil(t, err)
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	content := string(s)
	assert.Contains(t, content, "sql_client_requests_duration_ms_sum{command=\"test-cmd\"} 8\n")
	assert.Contains(t, content, "sql_client_requests_duration_ms_count{command=\"test-cmd\"} 1\n")
	assert.Contains(t, content, "sql_client_requests_error_total{command=\"test-cmd\",error=\"internal-error\"} 1\n")
	assert.Contains(t, content, "sql_client_requests_slow_total{command=\"test-cmd\"} 1\n")
}

func TestMetricCollector(t *testing.T) {
	prometheus.Unregister(connCollector)
	c := newCollector()
	c.registerClient(&statGetter{
		dbName: "db-1",
		hash:   "hash-1",
		poolStats: func() sql.DBStats {
			return sql.DBStats{
				MaxOpenConnections: 1,
				OpenConnections:    2,
				InUse:              3,
				Idle:               4,
				WaitCount:          5,
				WaitDuration:       6 * time.Second,
				MaxIdleClosed:      7,
				MaxIdleTimeClosed:  8,
				MaxLifetimeClosed:  9,
			}
		},
	})
	c.registerClient(&statGetter{
		dbName: "db-1",
		hash:   "hash-2",
		poolStats: func() sql.DBStats {
			return sql.DBStats{
				MaxOpenConnections: 10,
				OpenConnections:    20,
				InUse:              30,
				Idle:               40,
				WaitCount:          50,
				WaitDuration:       60 * time.Second,
				MaxIdleClosed:      70,
				MaxIdleTimeClosed:  80,
				MaxLifetimeClosed:  90,
			}
		},
	})
	c.registerClient(&statGetter{
		dbName: "db-2",
		hash:   "hash-2",
		poolStats: func() sql.DBStats {
			return sql.DBStats{
				MaxOpenConnections: 100,
				OpenConnections:    200,
				InUse:              300,
				Idle:               400,
				WaitCount:          500,
				WaitDuration:       600 * time.Second,
				MaxIdleClosed:      700,
				MaxIdleTimeClosed:  800,
				MaxLifetimeClosed:  900,
			}
		},
	})
	val := `
		# HELP sql_client_idle_connections The number of idle connections.
		# TYPE sql_client_idle_connections gauge
		sql_client_idle_connections{db_name="db-1",hash="hash-1"} 4
		sql_client_idle_connections{db_name="db-1",hash="hash-2"} 40
		sql_client_idle_connections{db_name="db-2",hash="hash-2"} 400
		# HELP sql_client_in_use_connections The number of connections currently in use.
		# TYPE sql_client_in_use_connections gauge
		sql_client_in_use_connections{db_name="db-1",hash="hash-1"} 3
		sql_client_in_use_connections{db_name="db-1",hash="hash-2"} 30
		sql_client_in_use_connections{db_name="db-2",hash="hash-2"} 300
		# HELP sql_client_max_idle_closed_total The total number of connections closed due to SetMaxIdleConns.
		# TYPE sql_client_max_idle_closed_total counter
		sql_client_max_idle_closed_total{db_name="db-1",hash="hash-1"} 7
		sql_client_max_idle_closed_total{db_name="db-1",hash="hash-2"} 70
		sql_client_max_idle_closed_total{db_name="db-2",hash="hash-2"} 700
		# HELP sql_client_max_idle_time_closed_total The total number of connections closed due to SetConnMaxIdleTime.
		# TYPE sql_client_max_idle_time_closed_total counter
		sql_client_max_idle_time_closed_total{db_name="db-1",hash="hash-1"} 8
		sql_client_max_idle_time_closed_total{db_name="db-1",hash="hash-2"} 80
		sql_client_max_idle_time_closed_total{db_name="db-2",hash="hash-2"} 800
		# HELP sql_client_max_lifetime_closed_total The total number of connections closed due to SetConnMaxLifetime.
		# TYPE sql_client_max_lifetime_closed_total counter
		sql_client_max_lifetime_closed_total{db_name="db-1",hash="hash-1"} 9
		sql_client_max_lifetime_closed_total{db_name="db-1",hash="hash-2"} 90
		sql_client_max_lifetime_closed_total{db_name="db-2",hash="hash-2"} 900
		# HELP sql_client_max_open_connections Maximum number of open connections to the database.
		# TYPE sql_client_max_open_connections gauge
		sql_client_max_open_connections{db_name="db-1",hash="hash-1"} 1
		sql_client_max_open_connections{db_name="db-1",hash="hash-2"} 10
		sql_client_max_open_connections{db_name="db-2",hash="hash-2"} 100
		# HELP sql_client_open_connections The number of established connections both in use and idle.
		# TYPE sql_client_open_connections gauge
		sql_client_open_connections{db_name="db-1",hash="hash-1"} 2
		sql_client_open_connections{db_name="db-1",hash="hash-2"} 20
		sql_client_open_connections{db_name="db-2",hash="hash-2"} 200
		# HELP sql_client_wait_count_total The total number of connections waited for.
		# TYPE sql_client_wait_count_total counter
		sql_client_wait_count_total{db_name="db-1",hash="hash-1"} 5
		sql_client_wait_count_total{db_name="db-1",hash="hash-2"} 50
		sql_client_wait_count_total{db_name="db-2",hash="hash-2"} 500
		# HELP sql_client_wait_duration_seconds_total The total time blocked waiting for a new connection.
		# TYPE sql_client_wait_duration_seconds_total counter
		sql_client_wait_duration_seconds_total{db_name="db-1",hash="hash-1"} 6
		sql_client_wait_duration_seconds_total{db_name="db-1",hash="hash-2"} 60
		sql_client_wait_duration_seconds_total{db_name="db-2",hash="hash-2"} 600
`

	err := testutil.CollectAndCompare(c, strings.NewReader(val))
	assert.NoError(t, err)
}
