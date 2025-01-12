package redis

import (
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyKey is an error that indicates no redis key is set.
	ErrEmptyKey = errors.New("empty redis key")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Host     string
		Type     string `json:",default=node,options=node|cluster"`
		Pass     string `json:",optional"`
		Tls      bool   `json:",optional"`
		NonBlock bool   `json:",default=true"`
		// PingTimeout is the timeout for ping redis.
		PingTimeout time.Duration `json:",default=1s"`
		// VerifyBigKey is the config for big key hook.
		VerifyBigKey BigKeyHookConfig
	}

	// A RedisKeyConf is a redis config with key.
	RedisKeyConf struct {
		RedisConf
		Key string
	}
)

// NewRedis returns a Redis.
// Deprecated: use MustNewRedis or NewRedis instead.
func (rc RedisConf) NewRedis() *Redis {
	var opts []Option
	if rc.Type == ClusterType {
		opts = append(opts, Cluster())
	}
	if len(rc.Pass) > 0 {
		opts = append(opts, WithPass(rc.Pass))
	}
	if rc.Tls {
		opts = append(opts, WithTLS())
	}
	if rc.VerifyBigKey.Enable {
		bigKeyHook, err := NewBigKeyHook(rc.VerifyBigKey)
		if err != nil {
			logx.Errorf("NewBigKeyHook error: %v", err)
		} else {
			opts = append(opts, WithHook(bigKeyHook))
		}
	}

	return newRedis(rc.Host, opts...)
}

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

// Validate validates the RedisKeyConf.
func (rkc RedisKeyConf) Validate() error {
	if err := rkc.RedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
