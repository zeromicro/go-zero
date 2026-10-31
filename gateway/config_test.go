package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteMapping_GetMethod(t *testing.T) {
	tests := []struct {
		name     string
		mapping  RouteMapping
		expected string
	}{
		{
			name: "top-level method only",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
			},
			expected: "GET",
		},
		{
			name: "match method takes precedence",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
				Match: &Match{
					Method: "POST",
					Path:   "/api/users",
				},
			},
			expected: "POST",
		},
		{
			name: "match with empty method falls back to top-level",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
				Match: &Match{
					Path: "/api/users",
				},
			},
			expected: "GET",
		},
		{
			name: "match only, no top-level",
			mapping: RouteMapping{
				Match: &Match{
					Method: "PUT",
					Path:   "/api/users",
				},
			},
			expected: "PUT",
		},
		{
			name:     "nil match with empty top-level",
			mapping:  RouteMapping{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.mapping.GetMethod())
		})
	}
}

func TestRouteMapping_GetPath(t *testing.T) {
	tests := []struct {
		name     string
		mapping  RouteMapping
		expected string
	}{
		{
			name: "top-level path only",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
			},
			expected: "/users",
		},
		{
			name: "match path takes precedence",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
				Match: &Match{
					Method: "GET",
					Path:   "/api/v1/users",
				},
			},
			expected: "/api/v1/users",
		},
		{
			name: "match with empty path falls back to top-level",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
				Match: &Match{
					Method: "POST",
				},
			},
			expected: "/users",
		},
		{
			name: "match only, no top-level",
			mapping: RouteMapping{
				Match: &Match{
					Method: "GET",
					Path:   "/api/v1/users/:id",
				},
			},
			expected: "/api/v1/users/:id",
		},
		{
			name:     "nil match with empty top-level",
			mapping:  RouteMapping{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.mapping.GetPath())
		})
	}
}

func TestRouteMapping_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mapping RouteMapping
		wantErr bool
	}{
		{
			name: "valid top-level fields",
			mapping: RouteMapping{
				Method: "GET",
				Path:   "/users",
			},
		},
		{
			name: "valid match fields",
			mapping: RouteMapping{
				Match: &Match{
					Method: "POST",
					Path:   "/api/users",
				},
			},
		},
		{
			name: "valid mixed: match path + top-level method",
			mapping: RouteMapping{
				Method: "GET",
				Match: &Match{
					Path: "/api/users",
				},
			},
		},
		{
			name:    "empty mapping",
			mapping: RouteMapping{},
			wantErr: true,
		},
		{
			name: "missing method",
			mapping: RouteMapping{
				Path: "/users",
			},
			wantErr: true,
		},
		{
			name: "missing path",
			mapping: RouteMapping{
				Method: "GET",
			},
			wantErr: true,
		},
		{
			name: "match with method only, no path anywhere",
			mapping: RouteMapping{
				Match: &Match{
					Method: "GET",
				},
			},
			wantErr: true,
		},
		{
			name: "match with path only, no method anywhere",
			mapping: RouteMapping{
				Match: &Match{
					Path: "/users",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mapping.Validate()
			if tt.wantErr {
				assert.ErrorIs(t, err, errMissingMethodPath)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
