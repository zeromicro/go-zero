package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestNoneSelector(t *testing.T) {
	selector := noneSelector{}
	assert.Equal(t, "", selector.Name())
	assert.Equal(t, []Conn(nil), selector.Select(nil, balancer.PickInfo{}))
}
