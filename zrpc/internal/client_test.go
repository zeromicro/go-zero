package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithDialOption(t *testing.T) {
	var options ClientOptions
	agent := grpc.WithUserAgent("chrome")
	opt := WithDialOption(agent)
	opt(&options)
	assert.Contains(t, options.DialOptions, agent)
}

func TestWithTimeout(t *testing.T) {
	var options ClientOptions
	opt := WithTimeout(time.Second)
	opt(&options)
	assert.Equal(t, time.Second, options.Timeout)
}

func TestBuildDialOptions(t *testing.T) {
	agent := grpc.WithUserAgent("chrome")
	opts := buildDialOptions(WithDialOption(agent))
	assert.Contains(t, opts, agent)
}
