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
		Strict      bool          `default:"false"`
		Expiry      time.Duration `default:"1h"`
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
		service.ServiceConf `mapstructure:",squash"`
		Host     string `default:"0.0.0.0"`
		Port     int
		CertFile string `json:",optional"`
		KeyFile  string `json:",optional"`
		Verbose  bool   `json:",optional"`
		MaxConns int    `default:"10000"`
		MaxBytes int64  `default:"1048576"`
		// milliseconds
		Timeout      int64         `default:"3000"`
		CpuThreshold int64         `default:"900" validate:"min=0,max=1000"`
		Signature    SignatureConf `json:",optional"`
	}
)
