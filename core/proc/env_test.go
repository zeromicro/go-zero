package proc

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	assert.True(t, len(Env("any")) == 0)
	envLock.RLock()
	val, ok := envs["any"]
	envLock.RUnlock()
	assert.True(t, len(val) == 0)
	assert.True(t, ok)
	assert.True(t, len(Env("any")) == 0)
}

func TestEnvInt(t *testing.T) {
	val, ok := EnvInt("any")
	assert.Equal(t, 0, val)
	assert.False(t, ok)
	err := os.Setenv("anyInt", "10")
	assert.Nil(t, err)
	val, ok = EnvInt("anyInt")
	assert.Equal(t, 10, val)
	assert.True(t, ok)
	err = os.Setenv("anyString", "a")
	assert.Nil(t, err)
	val, ok = EnvInt("anyString")
	assert.Equal(t, 0, val)
	assert.False(t, ok)
}
