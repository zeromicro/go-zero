package config

import "github.com/l306287405/go-zero/core/logx"

// Config defines a service configure for goctl update
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
