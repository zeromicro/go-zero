package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubeBuilder_Scheme(t *testing.T) {
	var b kubeBuilder
	assert.Equal(t, KubernetesScheme, b.Scheme())
}
