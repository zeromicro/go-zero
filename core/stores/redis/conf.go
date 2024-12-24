package redis

import (
	"errors"
	"time"
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
		Host        string
		Type        string        `json:",default=node,options=node|cluster"`
		Username    string        `json:",optional"`
		Password    string        `json:",optional"`
		Tls         bool          `json:",optional"`
		NonBlock    bool          `json:",default=true"`
		PingTimeout time.Duration `json:",default=1s"`
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
	if len(rc.Username) > 0 {
		opts = append(opts, WithUsername(rc.Username))
	}
	if len(rc.Password) > 0 {
		opts = append(opts, WithPassword(rc.Password))
	}
	if rc.Tls {
		opts = append(opts, WithTLS())
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
