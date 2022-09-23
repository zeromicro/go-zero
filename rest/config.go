package rest

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"
)

type (
	// A PrivateKeyConf is a private key config.
	PrivateKeyConf struct {
		Fingerprint string
		KeyFile     string
	}

	// A SignatureConf is a signature config.
	SignatureConf struct {
		Strict      bool          `json:",default=false"`
		Expiry      time.Duration `json:",default=1h"`
		PrivateKeys []PrivateKeyConf
	}

	// A RestConf is a http service config.
	// Why not name it as Conf, because we need to consider usage like:
	//  type Config struct {
	//     zrpc.RpcConf
	//     rest.RestConf
	//  }
	// if with the name Conf, there will be two Conf inside Config.
	RestConf struct {
		service.ServiceConf `yaml:",inline"`
		Host                string        `json:"host,default=0.0.0.0" yaml:"Host"`
		Port                int           `json:"port" yaml:"Port"`
		CertFile            string        `json:"certFile,optional" yaml:"CertFile"`
		KeyFile             string        `json:"KeyFile,optional" yaml:"KeyFile"`
		Verbose             bool          `json:"verbose,optional" yaml:"Verbose"`
		MaxConns            int           `json:"maxConns,default=10000" yaml:"MaxConns"`
		MaxBytes            int64         `json:"maxBytes,default=1048576" yaml:"MaxBytes"`
		Timeout             int64         `json:"timeout,default=3000" yaml:"Timeout"` // milliseconds
		CpuThreshold        int64         `json:"cpuThreshold,default=900,range=[0:1000]" yaml:"CpuThreshold"`
		Signature           SignatureConf `json:"signature,optional" yaml:"Signature"`
	}
)
