package config

import "github.com/zeromicro/go-zero/core/logx"

// Config defines a service configure for cztctl update
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
