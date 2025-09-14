// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1-alpha

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Trans struct {
		Secret     string
		PrevSecret string
	}
}
