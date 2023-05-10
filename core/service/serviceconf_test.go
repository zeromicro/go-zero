package service

import (
	"testing"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestServiceConf(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "console",
		},
		SqlLog: sqlx.LogConf{
			DisableStmtLog: false,
			DisableSqlLog:  false,
		},
		Mode: "dev",
	}
	c.MustSetUp()
}

func TestServiceConfWithMetricsUrl(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "volume",
		},
		SqlLog: sqlx.LogConf{
			DisableStmtLog: true,
			DisableSqlLog:  true,
		},
		Mode:       "dev",
		MetricsUrl: "http://localhost:8080",
	}
	assert.NoError(t, c.SetUp())
}
