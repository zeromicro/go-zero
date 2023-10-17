package service

import (
	"github.com/zeromicro/go-zero/internal/devserver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestServiceConf(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "console",
		},
		Mode: "dev",
		DevServer: devserver.Config{
			Port:       6470,
			HealthPath: "/healthz",
		},
	}
	c.MustSetUp()
}

func TestServiceConfWithMetricsUrl(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "volume",
		},
		Mode:       "dev",
		MetricsUrl: "http://localhost:8080",
	}
	assert.NoError(t, c.SetUp())
}
