package redistest

import (
	"time"

	"github.com/3Rivers/go-zero/core/lang"
	"github.com/3Rivers/go-zero/core/stores/redis"
	"github.com/alicebob/miniredis/v2"
)

func CreateRedis() (r *redis.Redis, clean func(), err error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	return redis.NewRedis(mr.Addr(), redis.NodeType, true), func() {
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
