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
		Host     string
		Type     string `json:",default=node,options=node|cluster"`
		User     string `json:",optional"`
		Pass     string `json:",optional"`
		Tls      bool   `json:",optional"`
		NonBlock bool   `json:",default=true"`
		// DisableIdentity is used to disable CLIENT SETINFO command on connect.
		//
		// Some redis versions/proxies do not support CLIENT SETINFO and return an
		// error on connect; since that command runs through the breaker hook it can
		// trip the breaker. Set this to true to skip it on such servers. Together
		// with the default MaintNotifications=disabled (and the always-ignored
		// HELLO command), this keeps the connect-time commands from tripping the
		// breaker on incompatible servers, without forcing RESP2.
		//
		// default: false
		DisableIdentity bool `json:",default=false"`
		// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
		//
		// default: 3.
		Protocol int `json:",default=3"`
		// MaintNotifications controls the CLIENT MAINT_NOTIFICATIONS handshake mode
		// (go-redis MaintNotificationsConfig.Mode):
		//   - disabled: never send the command (avoids tripping the breaker on servers
		//     that don't support it; keeps RESP3 intact)
		//   - auto: try, silently fall back on error (go-redis default)
		//   - enabled: force, fail the connection on error
		//
		// default: disabled
		MaintNotifications string `json:",default=disabled,options=disabled|enabled|auto"`
		// PingTimeout is the timeout for ping redis.
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
	if len(rc.User) > 0 {
		opts = append(opts, WithUser(rc.User))
	}
	if len(rc.Pass) > 0 {
		opts = append(opts, WithPass(rc.Pass))
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
