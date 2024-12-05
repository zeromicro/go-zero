package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/baggage"
)

func TestAddBaggagesFromMap(t *testing.T) {
	ctx := AddBaggagesFromMap(context.Background(), map[string]string{"test": "test", "aaa": "aaa"})
	b := baggage.FromContext(ctx)
	assert.Equal(t, "test", b.Member("test").Value())
	assert.Equal(t, "aaa", b.Member("aaa").Value())
}

func TestGetBaggageValue(t *testing.T) {
	ctx := AddBaggagesFromMap(context.Background(), map[string]string{"test": "aaa"})

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			"not exists",
			args{
				ctx: context.Background(),
				key: "test",
			},
			"",
			false,
		},
		{
			"exists",
			args{
				ctx: ctx,
				key: "test",
			},
			"aaa",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetBaggageValue(tt.args.ctx, tt.args.key)
			if got != tt.want {
				t.Errorf("GetBaggageValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetBaggageValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestWithBaggage(t *testing.T) {
	ctx := context.Background()
	ctx = WithBaggage(ctx, "aaa", "aaa")
	val, ok := GetBaggageValue(ctx, "aaa")
	assert.True(t, ok)
	assert.Equal(t, "aaa", val)

	ctx = WithBaggage(ctx, "aaa", "bbb")
	val, ok = GetBaggageValue(ctx, "aaa")
	assert.True(t, ok)
	assert.Equal(t, "bbb", val)
}
