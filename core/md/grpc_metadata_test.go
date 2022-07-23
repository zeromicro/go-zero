package md

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGRPCMetadataCarrier_Carrier(t *testing.T) {
	t.Run("has metadata", func(t *testing.T) {
		carrier := GRPCMetadataCarrier{"metadata": {`{"colors":["v1"]}`}}
		metadata, err := carrier.Carrier()
		assert.NoError(t, err)
		assert.Equal(t, Metadata{"colors": {"v1"}}, metadata)
	})
	t.Run("no metadata", func(t *testing.T) {
		carrier := GRPCMetadataCarrier{}
		metadata, err := carrier.Carrier()
		assert.NoError(t, err)
		assert.Equal(t, Metadata{}, metadata)

		carrier = GRPCMetadataCarrier{"metadata": {}}
		metadata, err = carrier.Carrier()
		assert.NoError(t, err)
		assert.Equal(t, Metadata{}, metadata)
	})
	t.Run("no metadata", func(t *testing.T) {
		carrier := GRPCMetadataCarrier{"metadata": {`{1}`}}
		metadata, err := carrier.Carrier()
		assert.Error(t, err)
		assert.Nil(t, metadata)
	})

}
