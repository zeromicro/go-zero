package md

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGRPCMetadataCarrier_Extract(t *testing.T) {
	t.Run("no err", func(t *testing.T) {
		carrier := GRPCMetadataCarrier(metadata.MD{"metadata": {`{"a":["a1","a2"]}`, "v2"}})
		ctx, err := carrier.Extract(context.Background())
		assert.NoError(t, err)
		md := FromContext(ctx)
		assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}}, md)
	})

	t.Run("err", func(t *testing.T) {
		carrier := GRPCMetadataCarrier(metadata.MD{"metadata": {`{"a":["a1","a2"1]}`, "v2"}})
		ctx, err := carrier.Extract(context.Background())
		assert.Error(t, err)
		md := FromContext(ctx)
		assert.EqualValues(t, map[string][]string{}, md)
	})

	t.Run("no metadata", func(t *testing.T) {
		carrier := GRPCMetadataCarrier(metadata.MD{"metadata": {}})
		ctx, err := carrier.Extract(context.Background())
		assert.NoError(t, err)
		md := FromContext(ctx)
		assert.EqualValues(t, map[string][]string{}, md)

		carrier = GRPCMetadataCarrier(metadata.MD{})
		ctx, err = carrier.Extract(context.Background())
		assert.NoError(t, err)
		md = FromContext(ctx)
		assert.EqualValues(t, map[string][]string{}, md)
	})
}

func TestGRPCMetadataCarrier_Injection(t *testing.T) {
	t.Run("has metadata", func(t *testing.T) {
		md := metadata.MD{}
		err := GRPCMetadataCarrier(md).Inject(NewContext(context.Background(), Metadata{"a": {"a1", "a2"}}))
		assert.NoError(t, err)
		assert.EqualValues(t, map[string][]string{"metadata": {`{"a":["a1","a2"]}`}}, md)
	})

	t.Run("no metadata", func(t *testing.T) {
		md := metadata.MD{}
		err := GRPCMetadataCarrier(md).Inject(NewContext(context.Background(), Metadata{}))
		assert.NoError(t, err)
		assert.EqualValues(t, map[string][]string{}, md)
	})
}
