package trace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
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

func TestCreateExporter_InvalidFilePath(t *testing.T) {
	logx.Disable()

	c := Config{
		Name:     "test-invalid-file",
		Endpoint: "/non-existent-directory/trace.log",
		Batcher:  kindFile,
	}

	_, err := createExporter(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file exporter endpoint error")
}

func TestCreateExporter_UnknownBatcher(t *testing.T) {
	logx.Disable()

	c := Config{
		Name:     "test-unknown",
		Endpoint: "localhost:1234",
		Batcher:  "unknown-batcher-type",
	}

	_, err := createExporter(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown exporter")
}

func TestCreateExporter_ValidExporters(t *testing.T) {
	logx.Disable()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file exporter",
			config: Config{
				Name:     "file-test",
				Endpoint: "/tmp/trace-test.log",
				Batcher:  kindFile,
			},
			wantErr: false,
		},
		{
			name: "invalid file path",
			config: Config{
				Name:     "file-test-invalid",
				Endpoint: "/invalid-path/that/does/not/exist/trace.log",
				Batcher:  kindFile,
			},
			wantErr: true,
			errMsg:  "file exporter endpoint error",
		},
		{
			name: "unknown batcher",
			config: Config{
				Name:     "unknown-test",
				Endpoint: "localhost:1234",
				Batcher:  "invalid-batcher",
			},
			wantErr: true,
			errMsg:  "unknown exporter",
		},
		{
			name: "jaeger http",
			config: Config{
				Name:     "jaeger-http",
				Endpoint: "http://localhost:14268/api/traces",
				Batcher:  kindJaeger,
			},
			wantErr: false,
		},
		{
			name: "jaeger udp",
			config: Config{
				Name:     "jaeger-udp",
				Endpoint: "udp://localhost:6831",
				Batcher:  kindJaeger,
			},
			wantErr: false,
		},
		{
			name: "zipkin",
			config: Config{
				Name:     "zipkin",
				Endpoint: "http://localhost:9411/api/v2/spans",
				Batcher:  kindZipkin,
			},
			wantErr: false,
		},
		{
			name: "otlpgrpc",
			config: Config{
				Name:     "otlpgrpc",
				Endpoint: "localhost:4317",
				Batcher:  kindOtlpGrpc,
			},
			wantErr: false,
		},
		{
			name: "otlpgrpc with headers",
			config: Config{
				Name:     "otlpgrpc-headers",
				Endpoint: "localhost:4317",
				Batcher:  kindOtlpGrpc,
				OtlpHeaders: map[string]string{
					"authorization": "Bearer token123",
					"x-custom-key":  "custom-value",
				},
			},
			wantErr: false,
		},
		{
			name: "otlphttp",
			config: Config{
				Name:     "otlphttp",
				Endpoint: "localhost:4318",
				Batcher:  kindOtlpHttp,
			},
			wantErr: false,
		},
		{
			name: "otlphttp with headers",
			config: Config{
				Name:     "otlphttp-headers",
				Endpoint: "localhost:4318",
				Batcher:  kindOtlpHttp,
				OtlpHeaders: map[string]string{
					"authorization": "Bearer token456",
					"x-api-key":     "api-key-value",
				},
			},
			wantErr: false,
		},
		{
			name: "otlphttp with headers and path",
			config: Config{
				Name:         "otlphttp-headers-path",
				Endpoint:     "localhost:4318",
				Batcher:      kindOtlpHttp,
				OtlpHttpPath: "/v1/traces",
				OtlpHeaders: map[string]string{
					"authorization":  "Bearer token789",
					"x-custom-trace": "trace-id",
				},
			},
			wantErr: false,
		},
		{
			name: "otlphttp with secure connection",
			config: Config{
				Name:           "otlphttp-secure",
				Endpoint:       "localhost:4318",
				Batcher:        kindOtlpHttp,
				OtlpHttpSecure: true,
				OtlpHeaders: map[string]string{
					"authorization": "Bearer secure-token",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exporter, err := createExporter(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, exporter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, exporter)
				// Clean up the exporter
				if exporter != nil {
					_ = exporter.Shutdown(context.Background())
				}
			}
		})
	}
}

func TestStopAgent(t *testing.T) {
	logx.Disable()

	// StopAgent should be idempotent and safe to call multiple times
	assert.NotPanics(t, func() {
		StopAgent()
		StopAgent()
		StopAgent()
	})
}

func TestStartAgent_WithEndpoint(t *testing.T) {
	logx.Disable()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "empty endpoint - no exporter created",
			config: Config{
				Name:    "test-no-endpoint",
				Sampler: 1.0,
			},
			wantErr: false,
		},
		{
			name: "valid endpoint with file exporter",
			config: Config{
				Name:     "test-with-endpoint",
				Endpoint: "/tmp/test-trace.log",
				Batcher:  kindFile,
				Sampler:  1.0,
			},
			wantErr: false,
		},
		{
			name: "endpoint with invalid exporter type",
			config: Config{
				Name:     "test-invalid-batcher",
				Endpoint: "localhost:1234",
				Batcher:  "invalid-type",
				Sampler:  1.0,
			},
			wantErr: true,
		},
		{
			name: "endpoint with invalid file path",
			config: Config{
				Name:     "test-invalid-path",
				Endpoint: "/non/existent/path/trace.log",
				Batcher:  kindFile,
				Sampler:  1.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset tp for each test
			originalTp := tp
			tp = nil
			defer func() {
				if tp != nil {
					_ = tp.Shutdown(context.Background())
				}
				tp = originalTp
			}()

			err := startAgent(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tp, "TracerProvider should be created")
			}
		})
	}
}

func TestStartAgent_ErrorHandler(t *testing.T) {
	// Setup a tracer provider to test error handler
	originalTp := tp
	tp = nil
	defer func() {
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
		tp = originalTp
	}()

	// Call startAgent to set up the error handler
	config := Config{
		Name:    "test-error-handler",
		Sampler: 1.0,
	}
	err := startAgent(config)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	// Verify the error handler was set and can be called without panicking
	// We test this by calling otel.Handle which will invoke the registered error handler
	testErr := errors.New("test otel error")
	assert.NotPanics(t, func() {
		otel.Handle(testErr)
	}, "Error handler should handle errors without panicking")
}
