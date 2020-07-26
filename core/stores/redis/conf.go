package redis

import "errors"

var (
	ErrEmptyHost = errors.New("empty redis host")
	ErrEmptyType = errors.New("empty redis type")
	ErrEmptyKey  = errors.New("empty redis key")
)

type (
	RedisConf struct {
		Host string
		Type string `json:",default=node,options=node|cluster"`
		Pass string `json:",optional"`
	}

	RedisKeyConf struct {
		RedisConf
		Key string `json:",optional"`
	}
)

func (rc RedisConf) NewRedis() *Redis {
	return NewRedis(rc.Host, rc.Type, rc.Pass)
}

func (rc RedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

func (rkc RedisKeyConf) Validate() error {
	if err := rkc.RedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
