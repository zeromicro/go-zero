package redistest

import (
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

func CreateRedis() (r *redis.Redis, clean func(), err error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	return redis.NewRedis(mr.Addr(), redis.NodeType), func() {
		ch := make(chan lang.PlaceholderType)
		go func() {
			mr.Close()
			close(ch)
		}()
		select {
		case <-ch:
		case <-time.After(time.Second):
		}
	}, nil
}
