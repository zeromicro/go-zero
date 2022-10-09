package rest

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"
)

type (
	// A PrivateKeyConf is a private key config.
	PrivateKeyConf struct {
		Fingerprint string `json:"Fingerprint,omitempty" yaml:"Fingerprint"`
		KeyFile     string `json:"KeyFile,omitempty" yaml:"KeyFile"`
	}

	// A SignatureConf is a signature config.
	SignatureConf struct {
		Strict      bool             `json:"Strict,default=false" yaml:"Strict"`
		Expiry      time.Duration    `json:"Expiry,default=1h" yaml:"Expiry"`
		PrivateKeys []PrivateKeyConf `yaml:",inline"`
	}

	// AuthConf is a JWT config
	AuthConf struct {
		AccessSecret string `json:"AccessSecret" yaml:"AccessSecret"`
		AccessExpire int64  `json:"AccessExpire" yaml:"AccessExpire"`
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
		Host                string        `json:"Host,default=0.0.0.0" yaml:"Host"`
		Port                int           `json:"Port" yaml:"Port"`
		CertFile            string        `json:"CertFile,optional" yaml:"CertFile"`
		KeyFile             string        `json:"KeyFile,optional" yaml:"KeyFile"`
		Verbose             bool          `json:"Verbose,optional" yaml:"Verbose"`
		MaxConns            int           `json:"MaxConns,default=10000" yaml:"MaxConns"`
		MaxBytes            int64         `json:"MaxBytes,default=1048576" yaml:"MaxBytes"`
		Timeout             int64         `json:"Timeout,default=3000" yaml:"Timeout"` // milliseconds
		CpuThreshold        int64         `json:"CpuThreshold,default=900,range=[0:1000]" yaml:"CpuThreshold"`
		Signature           SignatureConf `json:"Signature,optional" yaml:"Signature"`
	}
)
