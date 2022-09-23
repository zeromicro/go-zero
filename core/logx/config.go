package logx

// A LogConf is a logging config.
type LogConf struct {
	ServiceName         string `json:"serviceName,optional" yaml:"ServiceName"`
	Mode                string `json:"mode,default=console,options=[console,file,volume]" yaml:"Mode"`
	Encoding            string `json:"encoding,default=json,options=[json,plain]" yaml:"Encoding"`
	TimeFormat          string `json:"timeFormat,optional" yaml:"TimeFormat"`
	Path                string `json:"path,default=logs" yaml:"Path"`
	Level               string `json:"level,default=info,options=[debug,info,error,severe]" yaml:"Level"`
	Compress            bool   `json:"compress,optional" yaml:"Compress"`
	KeepDays            int    `json:"keepDays,optional" yaml:"KeepDays"`
	StackCooldownMillis int    `json:"stackCooldownMillis,default=100" yaml:"StackCooldownMillis"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even thougth `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `json:"maxBackups,default=0" yaml:"MaxBackups"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int `json:"maxSize,default=0" yaml:"MaxSize"`
	// RotationRuleType represents the type of log rotation rule. Default is `daily`.
	// daily: daily rotation.
	// size: size limited rotation.
	Rotation string `json:"rotation,default=daily,options=[daily,size]" yaml:"Rotation"`
}
