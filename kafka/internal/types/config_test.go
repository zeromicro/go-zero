package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientConfig_Validate(t *testing.T) {
	assert.Error(t, (&ClientConfig{}).Validate())
	assert.Error(t, (&ClientConfig{Brokers: []string{"b1"}, AuthType: PasswordAuthType}).Validate())
	assert.NoError(t, (&ClientConfig{Brokers: []string{"b1"}}).Validate())
}

func TestClientConfig_GetClientName(t *testing.T) {
	assert.Equal(t, "b1,b2", (&ClientConfig{Brokers: []string{"b1", "b2"}}).GetClientName())
}
