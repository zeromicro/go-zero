package fullpath

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPath(t *testing.T) {
	expect := []string{"/", "/api/v1", "/api/user/:name", "/api/user/:name/:id", "/:name/:id/:age"}
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	for _, p := range expect {
		r = WithFullPath(r, p)
		assert.EqualValues(t, p, FullPath(r))
	}

}
