package md

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderCarrier_Extract(t *testing.T) {
	carrier := HeaderCarrier(http.Header{"A": {"a1", "a2"}})
	ctx, err := carrier.Extract(context.Background())
	assert.NoError(t, err)
	assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}}, FromContext(ctx))
}

func TestHeaderCarrier_Injection(t *testing.T) {
	header := http.Header{}
	carrier := HeaderCarrier(header)
	err := carrier.Inject(NewContext(context.Background(), map[string][]string{"a": {"a1", "a2"}}))
	assert.NoError(t, err)
	assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}}, header)
}
