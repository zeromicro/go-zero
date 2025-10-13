package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestStartAgent(t *testing.T) {
	logx.Disable()

	const (
		endpoint1  = "localhost:1234"
		endpoint2  = "remotehost:1234"
		endpoint3  = "localhost:1235"
		endpoint4  = "localhost:1236"
		endpoint5  = "udp://localhost:6831"
		endpoint6  = "localhost:1237"
		endpoint71 = "/tmp/trace.log"
		endpoint72 = "/not-exist-fs/trace.log"
	)
	c1 := Config{
		Name: "foo",
	}
	c2 := Config{
		Name:     "bar",
		Endpoint: endpoint1,
		Batcher:  kindJaeger,
	}
	c3 := Config{
		Name:     "any",
		Endpoint: endpoint2,
		Batcher:  kindZipkin,
	}
	c4 := Config{
		Name:     "bla",
		Endpoint: endpoint3,
		Batcher:  "otlp",
	}
	c5 := Config{
		Name:     "otlpgrpc",
		Endpoint: endpoint3,
		Batcher:  kindOtlpGrpc,
		OtlpHeaders: map[string]string{
			"uptrace-dsn": "http://project2_secret_token@localhost:14317/2",
		},
	}
	c6 := Config{
		Name:     "otlphttp",
		Endpoint: endpoint4,
		Batcher:  kindOtlpHttp,
		OtlpHeaders: map[string]string{
			"uptrace-dsn": "http://project2_secret_token@localhost:14318/2",
		},
		OtlpHttpPath: "/v1/traces",
	}
	c7 := Config{
		Name:     "UDP",
		Endpoint: endpoint5,
		Batcher:  kindJaeger,
	}
	c8 := Config{
		Disabled: true,
		Endpoint: endpoint6,
		Batcher:  kindJaeger,
	}
	c9 := Config{
		Name:     "file",
		Endpoint: endpoint71,
		Batcher:  kindFile,
	}
	c10 := Config{
		Name:     "file",
		Endpoint: endpoint72,
		Batcher:  kindFile,
	}

	StartAgent(c1)
	StartAgent(c1)
	StartAgent(c2)
	StartAgent(c3)
	StartAgent(c4)
	StartAgent(c5)
	StartAgent(c6)
	StartAgent(c7)
	StartAgent(c8)
	StartAgent(c9)
	StartAgent(c10)
	defer StopAgent()

	// With sync.Once, only the first non-disabled config (c1) takes effect.
	// Subsequent calls are ignored, which is the desired behavior to prevent
	// multiple servers (REST + RPC) from reinitializing the global tracer.
	assert.NotNil(t, tp)
}
