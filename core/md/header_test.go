package md

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderCarrier_Carrier(t *testing.T) {
	t.Run("has metadata", func(t *testing.T) {
		carrier := HeaderCarrier(http.Header{"Colors": {"v1"}})
		metadata, err := carrier.Carrier()
		assert.NoError(t, err)
		assert.Equal(t, Metadata{"colors": {"v1"}}, metadata)
	})

	t.Run("no metadata", func(t *testing.T) {
		carrier := HeaderCarrier(http.Header{})
		metadata, err := carrier.Carrier()
		assert.NoError(t, err)
		assert.Equal(t, Metadata{}, metadata)
	})
}
