package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type mockConn struct {
	attributes *attributes.Attributes
}

func (c mockConn) Address() resolver.Address {
	return resolver.Address{
		BalancerAttributes: c.attributes,
	}
}

func Test_defaultSelector(t *testing.T) {
	selector := defaultSelector{}
	assert.Equal(t, "defaultSelector", selector.Name())
}
