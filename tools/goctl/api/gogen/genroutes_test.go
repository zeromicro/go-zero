package gogen

import (
	"testing"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func Test_formatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0 * time.Nanosecond"},
		{time.Nanosecond, "1 * time.Nanosecond"},
		{100 * time.Nanosecond, "100 * time.Nanosecond"},
		{500 * time.Microsecond, "500 * time.Microsecond"},
		{2 * time.Millisecond, "2 * time.Millisecond"},
		{time.Second, "1000 * time.Millisecond"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %v; want %v", test.duration, result, test.expected)
		}
	}
}

func TestSSESupport(t *testing.T) {
	// Test API spec with SSE enabled
	apiSpec := &spec.ApiSpec{
		Service: spec.Service{
			Groups: []spec.Group{
				{
					Annotation: spec.Annotation{
						Properties: map[string]string{
							"sse":    "true",
							"prefix": "/api/v1",
						},
					},
					Routes: []spec.Route{
						{
							Method:  "get",
							Path:    "/events",
							Handler: "StreamEvents",
						},
					},
				},
			},
		},
	}

	groups, err := getRoutes(apiSpec)
	if err != nil {
		t.Fatalf("getRoutes failed: %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	group := groups[0]
	if !group.sseEnabled {
		t.Error("Expected SSE to be enabled")
	}

	if group.prefix != "/api/v1" {
		t.Errorf("Expected prefix '/api/v1', got '%s'", group.prefix)
	}

	if len(group.routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(group.routes))
	}

	route := group.routes[0]
	if route.method != "http.MethodGet" {
		t.Errorf("Expected method 'http.MethodGet', got '%s'", route.method)
	}

	if route.path != "/events" {
		t.Errorf("Expected path '/events', got '%s'", route.path)
	}
}

func TestSSEWithOtherFeatures(t *testing.T) {
	// Test API spec with SSE and other features
	apiSpec := &spec.ApiSpec{
		Service: spec.Service{
			Groups: []spec.Group{
				{
					Annotation: spec.Annotation{
						Properties: map[string]string{
							"sse":        "true",
							"jwt":        "Auth",
							"signature":  "true",
							"prefix":     "/api/v1",
							"timeout":    "30s",
							"middleware": "AuthMiddleware,LogMiddleware",
						},
					},
					Routes: []spec.Route{
						{
							Method:  "get",
							Path:    "/events",
							Handler: "StreamEvents",
						},
					},
				},
			},
		},
	}

	groups, err := getRoutes(apiSpec)
	if err != nil {
		t.Fatalf("getRoutes failed: %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	group := groups[0]

	// Verify all features are enabled
	if !group.sseEnabled {
		t.Error("Expected SSE to be enabled")
	}

	if !group.jwtEnabled {
		t.Error("Expected JWT to be enabled")
	}

	if !group.signatureEnabled {
		t.Error("Expected signature to be enabled")
	}

	if group.authName != "Auth" {
		t.Errorf("Expected authName 'Auth', got '%s'", group.authName)
	}

	if group.prefix != "/api/v1" {
		t.Errorf("Expected prefix '/api/v1', got '%s'", group.prefix)
	}

	if group.timeout != "30s" {
		t.Errorf("Expected timeout '30s', got '%s'", group.timeout)
	}

	expectedMiddlewares := []string{"AuthMiddleware", "LogMiddleware"}
	if len(group.middlewares) != len(expectedMiddlewares) {
		t.Errorf("Expected %d middlewares, got %d", len(expectedMiddlewares), len(group.middlewares))
	}

	for i, expected := range expectedMiddlewares {
		if group.middlewares[i] != expected {
			t.Errorf("Expected middleware[%d] '%s', got '%s'", i, expected, group.middlewares[i])
		}
	}
}

func TestSSEDisabled(t *testing.T) {
	// Test API spec without SSE
	apiSpec := &spec.ApiSpec{
		Service: spec.Service{
			Groups: []spec.Group{
				{
					Annotation: spec.Annotation{
						Properties: map[string]string{
							"prefix": "/api/v1",
						},
					},
					Routes: []spec.Route{
						{
							Method:  "get",
							Path:    "/status",
							Handler: "GetStatus",
						},
					},
				},
			},
		},
	}

	groups, err := getRoutes(apiSpec)
	if err != nil {
		t.Fatalf("getRoutes failed: %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	group := groups[0]
	if group.sseEnabled {
		t.Error("Expected SSE to be disabled")
	}
}
