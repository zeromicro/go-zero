package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertVersion(t *testing.T) {
	number, tag := convertVersion("1.1.10")
	assert.Equal(t, 1.110, number)
	assert.Equal(t, "", tag)

	number, tag = convertVersion("0.1.11")
	assert.Equal(t, 0.111, number)
	assert.Equal(t, "", tag)

	number, tag = convertVersion("1.11-pre")
	assert.Equal(t, 1.11, number)
	assert.Equal(t, "pre", tag)

	number, tag = convertVersion("1.11-beta-v1")
	assert.Equal(t, 1.11, number)
	assert.Equal(t, "betav1", tag)
}

func Test_IsVersionGatherThan(t *testing.T) {
	assert.False(t, IsVersionGreaterThan("0.11", "1.1"))
	assert.True(t, IsVersionGreaterThan("0.112", "0.1"))
	assert.True(t, IsVersionGreaterThan("1.1.10", "1.0.111"))
	assert.True(t, IsVersionGreaterThan("1.1.10", "1.1.10-pre"))
	assert.True(t, IsVersionGreaterThan("1.1.11-pre", "1.1.10"))
}
