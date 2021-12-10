package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	ztrace "github.com/tal-tech/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestOtelHandler(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			spanCtx := trace.SpanContextFromContext(ctx)
			assert.Equal(t, true, spanCtx.IsValid())
		}),
	)
	defer ts.Close()

	client := ts.Client()
	err := func(ctx context.Context) error {
		ctx, span := otel.Tracer("httptrace/client").Start(ctx, "test")
		defer span.End()

		req, _ := http.NewRequest("GET", ts.URL, nil)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		res, err := client.Do(req)
		assert.Equal(t, err, nil)
		_ = res.Body.Close()
		return nil
	}(context.Background())

	assert.Equal(t, err, nil)
}
