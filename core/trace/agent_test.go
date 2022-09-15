package trace

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestAwsXray(t *testing.T) {
	logx.Disable()

	c := Config{
		Name:        "grpc",
		Endpoint:    "localhost:1235",
		Batcher:     "grpc",
		IdGenerator: "xray",
		Propagator:  "xray",
	}
	assert.Len(t, agents, 0)
	StartAgent(c)
	assert.Len(t, agents, 1)
	a := agents[c.Endpoint]
	assert.NotNil(t, a)
	assert.IsType(t, otel.GetTextMapPropagator(), xray.Propagator{})
	_, span := otel.GetTracerProvider().Tracer("test").Start(context.TODO(), "N")
	// copy from xray test
	previousTime := time.Now().Unix()
	traceID := span.SpanContext().TraceID()
	expectedTraceIDLength := 32
	assert.Equal(t, len(traceID.String()), expectedTraceIDLength, "TraceID has incorrect length.")
	currentTime, err := strconv.ParseInt(traceID.String()[0:8], 16, 64)
	require.NoError(t, err)
	nextTime := time.Now().Unix()
	assert.LessOrEqual(t, previousTime, currentTime, "TraceID is generated incorrectly with the wrong timestamp.")
	assert.LessOrEqual(t, currentTime, nextTime, "TraceID is generated incorrectly with the wrong timestamp.")

	spanID := span.SpanContext().SpanID()
	expectedSpanIDLength := 16
	assert.Equal(t, len(spanID.String()), expectedSpanIDLength, "SpanID has incorrect length")

}

func TestStartAgent(t *testing.T) {
	logx.Disable()

	const (
		endpoint1 = "localhost:1234"
		endpoint2 = "remotehost:1234"
		endpoint3 = "localhost:1235"
		endpoint4 = "/stdout"
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
		Name:     "grpc",
		Endpoint: endpoint3,
		Batcher:  "grpc",
	}
	c6 := Config{
		Name:     "stdout",
		Endpoint: endpoint4,
		Batcher:  "stdout",
	}

	StartAgent(c1)
	StartAgent(c1)
	StartAgent(c2)
	StartAgent(c3)
	StartAgent(c4)
	StartAgent(c5)
	StartAgent(c6)

	lock.Lock()
	defer lock.Unlock()

	// because remotehost cannot be resolved
	assert.Equal(t, 4, len(agents))
	_, ok := agents[""]
	assert.True(t, ok)
	_, ok = agents[endpoint1]
	assert.True(t, ok)
	_, ok = agents[endpoint2]
	assert.False(t, ok)
	_, ok = agents[endpoint3]
	assert.True(t, ok)
	_, ok = agents[endpoint4]
	assert.True(t, ok)
}
