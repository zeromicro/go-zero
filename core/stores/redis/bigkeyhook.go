package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/threading"
)

type (
	bigKeyHook struct {
		config BigKeyHookConfig
		buffer chan bigKeyData
	}

	BigKeyHookConfig struct {
		Enable       bool          `json:",default=true"`
		LimitSize    int           `json:",default=10240"`
		LimitCount   int           `json:",default=5"`
		BufferLen    int           `json:",default=100"`
		StatInterval time.Duration `json:",default=5m"`
	}

	bigKeyData struct {
		key   string
		size  int
		count int
	}
)

func NewBigKeyHook(config BigKeyHookConfig) (red.Hook, error) {
	if config.LimitSize <= 0 {
		return nil, errors.New("limit size must be greater than 0")
	}
	if config.LimitCount <= 0 {
		config.LimitCount = 5
	}

	if config.BufferLen > 0 && config.StatInterval <= 0 {
		return nil, errors.New("stat interval must be greater than 0")
	}

	hook := &bigKeyHook{
		config: config,
		buffer: make(chan bigKeyData, config.BufferLen),
	}

	threading.GoSafe(hook.stat)

	return hook, nil
}

func (h *bigKeyHook) DialHook(next red.DialHook) red.DialHook {
	return next
}

func (h *bigKeyHook) ProcessHook(next red.ProcessHook) red.ProcessHook {
	return func(ctx context.Context, cmd red.Cmder) error {
		if !h.config.Enable {
			return next(ctx, cmd)
		}

		defer h.cmdCheck(ctx, cmd)

		return next(ctx, cmd)
	}
}

func (h *bigKeyHook) ProcessPipelineHook(next red.ProcessPipelineHook) red.ProcessPipelineHook {
	return func(ctx context.Context, cmds []red.Cmder) error {
		if !h.config.Enable {
			return next(ctx, cmds)
		}

		defer func() {
			for _, cmd := range cmds {
				h.cmdCheck(ctx, cmd)
			}
		}()

		return next(ctx, cmds)
	}
}

func (h *bigKeyHook) cmdCheck(ctx context.Context, cmd red.Cmder) {
	if h.config.LimitSize <= 0 || len(cmd.Args()) < 2 || cmd.Err() != nil {
		return
	}

	var (
		size int
		key  = mapping.Repr(cmd.Args()[1])
	)

	switch strings.ToLower(cmd.Name()) {
	case "get":
		c, ok := cmd.(*red.StringCmd)
		if !ok {
			return
		}
		size = len(c.Val())
	case "set", "setnx":
		if len(cmd.Args()) >= 3 {
			size = len(mapping.Repr(cmd.Args()[2]))
		}
	case "getset":
		c, ok := cmd.(*red.StringCmd)
		if !ok {
			return
		}

		if c.Err() == nil && len(c.Val()) > 0 {
			size = len(c.Val())
		} else if len(c.Args()) >= 3 {
			size = len(mapping.Repr(c.Args()[2]))
		}
	case "hgetall":
		c, ok := cmd.(*red.MapStringStringCmd)
		if !ok {
			return
		}
		println(c.Val())
		for _, v := range c.Val() {
			size += len(v)
		}
	case "hget":
		if cmd.Err() != nil {
			return
		}
		c, ok := cmd.(*red.StringCmd)
		if !ok {
			return
		}
		if len(cmd.Args()) >= 3 {
			key += ":" + mapping.Repr(cmd.Args()[2])
		}
		size = len(c.Val())
	case "hmget":
		c, ok := cmd.(*red.SliceCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(mapping.Repr(v))
		}
	case "hset", "hsetnx":
		if len(cmd.Args()) >= 4 {
			key += ":" + mapping.Repr(cmd.Args()[2])
			size = len(mapping.Repr(cmd.Args()[3]))
		}
	case "hmset":
		for i := 3; i < len(cmd.Args()); i += 2 {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "sadd":
		for i := 2; i < len(cmd.Args()); i++ {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "smembers":
		c, ok := cmd.(*red.StringSliceCmd)
		if !ok {
			return
		}
		for _, v := range c.Val() {
			size += len(v)
		}
	case "zrange":
		switch cmd.(type) {
		case *red.StringSliceCmd:
			for _, v := range cmd.(*red.StringSliceCmd).Val() {
				size += len(v)
			}
		case *red.ZSliceCmd:
			for _, v := range cmd.(*red.ZSliceCmd).Val() {
				size += len(mapping.Repr(v.Member))
			}
		}
	case "zadd":
		for i := 3; i < len(cmd.Args()); i += 2 {
			size += len(mapping.Repr(cmd.Args()[i]))
		}
	case "zrangebyscore":
		c, ok := cmd.(*red.ZSliceCmd)
		if !ok {
			return
		}

		for _, v := range c.Val() {
			size += len(mapping.Repr(v.Member))
		}
	default:
		return
	}

	if size > h.config.LimitSize {
		if h.config.BufferLen <= 0 {
			logc.Infof(ctx, "[REDIS] BigKey limit, key: %s, size: %d", key, size)
		} else {
			select {
			case h.buffer <- bigKeyData{key: key, size: size}:
			default:
				logc.Infof(ctx, "[REDIS] BigKey limit, key: %s, size: %d", key, size)
			}
		}
	}

	return
}

func (h *bigKeyHook) stat() {
	if !h.config.Enable || h.config.BufferLen <= 0 {
		return
	}

	// log the big key.
	for {
		for key, data := range h.getIntervalData() {
			if data.count >= h.config.LimitCount {
				logx.Infof("[REDIS] BigKey limit, key: %s, size: %d, count: %d", key, data.size, data.count)
			}
		}
	}
}

func (h *bigKeyHook) getIntervalData() map[string]bigKeyData {
	var m = make(map[string]bigKeyData)

	timeout := time.NewTimer(h.config.StatInterval)

	for {
		select {
		case data := <-h.buffer:
			if _, ok := m[data.key]; !ok {
				m[data.key] = data
			}

			m[data.key] = bigKeyData{
				key:   data.key,
				size:  data.size,
				count: m[data.key].count + 1,
			}
		case <-timeout.C:
			return m
		}
	}
}
