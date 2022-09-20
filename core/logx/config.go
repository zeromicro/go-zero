package logx

// A LogConf is a logging config.
type LogConf struct {
	ServiceName         string `json:",optional"`
	Mode                string `default:"console" validate:"oneof=console file volume"`
	Encoding            string `default:"json" validate:"oneof=plain json"`
	TimeFormat          string `json:",optional"` 
	Path                string `default:"logs"`
	Level               string `default:"info" validate:"oneof=debug info error severe"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCooldownMillis int    `default:"100"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even thougth `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `default:"0"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int `default:"0"`
	// RotationRuleType represents the type of log rotation rule. Default is `daily`.
	// daily: daily rotation.
	// size: size limited rotation.
	Rotation string `default:"daily" validate:"oneof=daily size"`
}
