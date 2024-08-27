package logx

// A LogConf is a logging config.
type LogConf struct {
	// ServiceName represents the service name.
	ServiceName string `json:",optional"`
	// Mode represents the logging mode, default is `console`.
	// console: log to console.
	// file: log to file.
	// volume: used in k8s, prepend the hostname to the log file name.
	Mode string `json:",default=console,options=[console,file,volume]"`
	// Encoding represents the encoding type, default is `json`.
	// json: json encoding.
	// plain: plain text encoding, typically used in development.
	Encoding string `json:",default=json,options=[json,plain]"`
	// TimeFormat represents the time format, default is `2006-01-02T15:04:05.000Z07:00`.
	TimeFormat string `json:",optional"`
	// Path represents the log file path, default is `logs`.
	Path string `json:",default=logs"`
	// Level represents the log level, default is `info`.
	Level string `json:",default=info,options=[debug,info,error,severe]"`
	// MaxContentLength represents the max content bytes, default is no limit.
	MaxContentLength uint32 `json:",optional"`
	// Compress represents whether to compress the log file, default is `false`.
	Compress bool `json:",optional"`
	// Stat represents whether to log statistics, default is `true`.
	Stat bool `json:",default=true"`
	// KeepDays represents how many days the log files will be kept. Default to keep all files.
	// Only take effect when Mode is `file` or `volume`, both work when Rotation is `daily` or `size`.
	KeepDays int `json:",optional"`
	// StackCooldownMillis represents the cooldown time for stack logging, default is 100ms.
	StackCooldownMillis int `json:",default=100"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even though `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `json:",default=0"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int `json:",default=0"`
	// Rotation represents the type of log rotation rule. Default is `daily`.
	// daily: daily rotation.
	// size: size limited rotation.
	Rotation string `json:",default=daily,options=[daily,size]"`
	// FileTimeFormat represents the time format for file name, default is `2006-01-02T15:04:05.000Z07:00`.
	FileTimeFormat string `json:",optional"`
}
