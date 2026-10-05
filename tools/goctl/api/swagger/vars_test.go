package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractPathPlaceholders(t *testing.T) {
	tests := []struct {
		name                string
		path                string
		expectedPlaceholders []string
	}{
		{
			name:                "empty path",
			path:                "",
			expectedPlaceholders: []string{},
		},
		{
			name:                "no placeholders",
			path:                "/api/v1/users",
			expectedPlaceholders: []string{},
		},
		{
			name:                "single :id style placeholder",
			path:                "/api/v1/users/:id",
			expectedPlaceholders: []string{"id"},
		},
		{
			name:                "single {id} style placeholder",
			path:                "/api/v1/users/{id}",
			expectedPlaceholders: []string{"id"},
		},
		{
			name:                "multiple :id style placeholders",
			path:                "/api/v1/:namespace/:id",
			expectedPlaceholders: []string{"namespace", "id"},
		},
		{
			name:                "multiple {id} style placeholders",
			path:                "/api/v1/{namespace}/{id}",
			expectedPlaceholders: []string{"namespace", "id"},
		},
		{
			name:                "mixed style placeholders",
			path:                "/api/v1/:namespace/users/{id}",
			expectedPlaceholders: []string{"namespace", "id"},
		},
		{
			name:                "placeholder at root",
			path:                "/:id",
			expectedPlaceholders: []string{"id"},
		},
		{
			name:                "placeholder after static segments",
			path:                "/foo/bar/:id",
			expectedPlaceholders: []string{"id"},
		},
		{
			name:                "empty placeholder name with colon",
			path:                "/api/v1/:/users",
			expectedPlaceholders: []string{},
		},
		{
			name:                "empty placeholder name with braces",
			path:                "/api/v1/{}/users",
			expectedPlaceholders: []string{},
		},
		{
			name:                "trailing slash with placeholder",
			path:                "/api/v1/:id/",
			expectedPlaceholders: []string{"id"},
		},
		{
			name:                "complex path with multiple placeholders",
			path:                "/api/v1/:org/:namespace/:id",
			expectedPlaceholders: []string{"org", "namespace", "id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeholders := extractPathPlaceholders(tt.path)

			// Check that we have the expected number of placeholders
			assert.Equal(t, len(tt.expectedPlaceholders), len(placeholders))

			// Check that all expected placeholders are present
			for _, expected := range tt.expectedPlaceholders {
				assert.True(t, placeholders[expected], "Expected placeholder '%s' to be present", expected)
			}
		})
	}
}

func TestExtractPathPlaceholders_Duplicates(t *testing.T) {
	// When the same placeholder appears multiple times, it should be deduplicated
	path := "/api/v1/:id/other/:id"
	placeholders := extractPathPlaceholders(path)

	assert.Equal(t, 1, len(placeholders))
	assert.True(t, placeholders["id"])
}

func TestExtractPathPlaceholders_CaseSensitive(t *testing.T) {
	// Placeholders should be case-sensitive
	path := "/api/v1/:ID/:id/:Id"
	placeholders := extractPathPlaceholders(path)

	assert.Equal(t, 3, len(placeholders))
	assert.True(t, placeholders["ID"])
	assert.True(t, placeholders["id"])
	assert.True(t, placeholders["Id"])
}