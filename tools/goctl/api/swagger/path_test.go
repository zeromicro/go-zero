package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestSpec2PathsWithRootRoute(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		routePath    string
		expectedPath string
	}{
		{
			name:         "prefix with root route",
			prefix:       "/api/v1/shoppings",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "prefix with sub route",
			prefix:       "/api/v1/shoppings",
			routePath:    "/list",
			expectedPath: "/api/v1/shoppings/list",
		},
		{
			name:         "empty prefix with root route",
			prefix:       "",
			routePath:    "/",
			expectedPath: "/",
		},
		{
			name:         "empty prefix with sub route",
			prefix:       "",
			routePath:    "/list",
			expectedPath: "/list",
		},
		{
			name:         "prefix with trailing slash and root route",
			prefix:       "/api/v1/shoppings/",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "prefix without leading slash and root route",
			prefix:       "api/v1/shoppings",
			routePath:    "/",
			expectedPath: "/api/v1/shoppings",
		},
		{
			name:         "single level prefix with root route",
			prefix:       "/api",
			routePath:    "/",
			expectedPath: "/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := spec.Service{
				Groups: []spec.Group{
					{
						Annotation: spec.Annotation{
							Properties: map[string]string{
								propertyKeyPrefix: tt.prefix,
							},
						},
						Routes: []spec.Route{
							{
								Method:  "get",
								Path:    tt.routePath,
								Handler: "TestHandler",
							},
						},
					},
				},
			}

			ctx := testingContext(t)
			paths := spec2Paths(ctx, srv)

			assert.Contains(t, paths.Paths, tt.expectedPath,
				"Expected path %s not found in generated paths. Got: %v",
				tt.expectedPath, paths.Paths)
		})
	}
}
