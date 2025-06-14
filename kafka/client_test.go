package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	config := UniversalClientConfig{Client: ClientConfig{
		ResourceName: "kafka-erms-test-1",
	}}
	_, err := NewClient(config)
	assert.Error(t, err)
}

func adfa() {

}
