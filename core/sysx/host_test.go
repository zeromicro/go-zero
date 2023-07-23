package sysx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostname(t *testing.T) {
	assert.True(t, len(Hostname()) > 0)
}
