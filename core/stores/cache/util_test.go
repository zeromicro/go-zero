package cache

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/lang"
)

func TestFormatKeys(t *testing.T) {
	assert.Equal(t, "a,b", formatKeys([]string{"a", "b"}))
}

func TestTotalWeights(t *testing.T) {
	val := TotalWeights([]NodeConf{
		{
			Weight: -1,
		},
		{
			Weight: 0,
		},
		{
			Weight: 1,
		},
	})
	assert.Equal(t, 1, val)
}

func createMiniRedis() (r *miniredis.Miniredis, clean func(), err error) {
	r, err = miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	return r, func() {
		ch := make(chan lang.PlaceholderType)
		go func() {
			r.Close()
			close(ch)
		}()
		select {
		case <-ch:
		case <-time.After(time.Second):
		}
	}, nil
}
