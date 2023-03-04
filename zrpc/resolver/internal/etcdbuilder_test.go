package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdBuilder_Scheme(t *testing.T) {
	assert.Equal(t, EtcdScheme, new(etcdBuilder).Scheme())
}
