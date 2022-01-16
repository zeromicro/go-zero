package logx

// A LogConf is a logging config.
type LogConf struct {
	ServiceName         string `json:",optional"`
	Mode                string `json:",default=console,options=[console,file,volume]"`
	Encoding            string `json:",default=json,options=[json,plain]"`
	TimeFormat          string `json:",optional"`
	Path                string `json:",default=logs"`
	Level               string `json:",default=info,options=[info,error,severe]"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCooldownMillis int    `json:",default=100"`
}
