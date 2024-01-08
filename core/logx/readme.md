<img align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# logx

English | [简体中文](readme-cn.md)

## logx configurations

```go
type LogConf struct {
	ServiceName         string              `json:",optional"`
	Mode                string              `json:",default=console,options=[console,file,volume]"`
	Encoding            string              `json:",default=json,options=[json,plain]"`
	TimeFormat          string              `json:",optional"`
	Path                string              `json:",default=logs"`
	Level               string              `json:",default=info,options=[info,error,severe]"`
	Compress            bool                `json:",optional"`
	KeepDays            int                 `json:",optional"`
	StackCooldownMillis int                 `json:",default=100"`
	MaxBackups          int                 `json:",default=0"`
	MaxSize             int                 `json:",default=0"`
	Rotation            string              `json:",default=daily,options=[daily,size]"`
}
```

- `ServiceName`: set the service name, optional. on `volume` mode, the name is used to generate the log files. Within `rest/zrpc` services, the name will be set to the name of `rest` or `zrpc` automatically.
- `Mode`: the mode to output the logs, default is `console`.
  -  `console` mode writes the logs to `stdout/stderr`.
  - `file` mode writes the logs to the files specified by `Path`.
  - `volume` mode is used in docker, to write logs into mounted volumes.
- `Encoding`: indicates how to encode the logs, default is `json`.
  - `json` mode writes the logs in json format.
  - `plain` mode writes the logs with plain text, with terminal color enabled.
- `TimeFormat`: customize the time format, optional. Default is `2006-01-02T15:04:05.000Z07:00`.
- `Path`: set the log path, default to `logs`.
- `Level`: the logging level to filter logs. Default is `info`.
  - `info`, all logs are written.
  - `error`, `info` logs are suppressed.
  - `severe`, `info` and `error` logs are suppressed, only `severe` logs are written.
- `Compress`: whether or not to compress log files, only works with `file` mode.
- `KeepDays`: how many days that the log files are kept, after the given days, the outdated files will be deleted automatically. It has no effect on `console` mode.
- `StackCooldownMillis`: how many milliseconds to rewrite stacktrace again. It’s used to avoid stacktrace flooding.
- `MaxBackups`: represents how many backup log files will be kept. 0 means all files will be kept forever. Only take effect when `Rotation` is `size`. NOTE: the level of option `KeepDays` will be higher. Even though `MaxBackups` sets 0, log files will still be removed if the `KeepDays` limitation is reached.
- `MaxSize`: represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`. Only take effect when `Rotation` is `size`.
- `Rotation`: represents the type of log rotation rule. Default is `daily`.
  - `daily` rotate the logs by day.
  - `size` rotate the logs by size of logs.

## Logging methods

```go
type Logger interface {
	// Error logs a message at error level.
	Error(...any)
	// Errorf logs a message at error level.
	Errorf(string, ...any)
	// Errorv logs a message at error level.
	Errorv(any)
	// Errorw logs a message at error level.
	Errorw(string, ...LogField)
	// Info logs a message at info level.
	Info(...any)
	// Infof logs a message at info level.
	Infof(string, ...any)
	// Infov logs a message at info level.
	Infov(any)
	// Infow logs a message at info level.
	Infow(string, ...LogField)
	// Slow logs a message at slow level.
	Slow(...any)
	// Slowf logs a message at slow level.
	Slowf(string, ...any)
	// Slowv logs a message at slow level.
	Slowv(any)
	// Sloww logs a message at slow level.
	Sloww(string, ...LogField)
	// WithContext returns a new logger with the given context.
	WithContext(context.Context) Logger
	// WithDuration returns a new logger with the given duration.
	WithDuration(time.Duration) Logger
}
```

- `Error`, `Info`, `Slow`: write any kind of messages into logs, with like `fmt.Sprint(…)`.
- `Errorf`, `Infof`, `Slowf`: write messages with given format into logs.
- `Errorv`, `Infov`, `Slowv`: write any kind of messages into logs, with json marshalling to encode them.
- `Errorw`, `Infow`, `Sloww`: write the string message with given `key:value` fields.
- `WithContext`: inject the given ctx into the log messages, typically used to log `trace-id` and `span-id`.
- `WithDuration`: write elapsed duration into the log messages, with key `duration`.

## Integrating with third-party logging libs

- zap
  - implementation: [https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go)
  - usage example: [https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go)
- logrus
  - implementation: [https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go)
  - usage example: [https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go)

For more libs, please implement and PR to [https://github.com/zeromicro/zero-contrib](https://github.com/zeromicro/zero-contrib)

## Write the logs to specific stores

`logx` defined two interfaces to let you customize `logx` to write logs into any stores.

- `logx.NewWriter(w io.Writer)`
- `logx.SetWriter(writer logx.Writer)`

For example, if we want to write the logs into kafka instead of console or files, we can do it like below:

```go
type KafkaWriter struct {
	Pusher *kq.Pusher
}

func NewKafkaWriter(pusher *kq.Pusher) *KafkaWriter {
	return &KafkaWriter{
		Pusher: pusher,
	}
}

func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	// writing log with newlines, trim them.
	if err := w.Pusher.Push(strings.TrimSpace(string(p))); err != nil {
		return 0, err
	}

	return len(p), nil
}

func main() {
	pusher := kq.NewPusher([]string{"localhost:9092"}, "go-zero")
	defer pusher.Close()

	writer := logx.NewWriter(NewKafkaWriter(pusher))
	logx.SetWriter(writer)
  
	// more code
}
```

Complete code: [https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go)

## Filtering sensitive fields

If we need to prevent the `password` fields from logging, we can do it like below:

```go
type (
	Message struct {
		Name     string
		Password string
		Message  string
	}

	SensitiveLogger struct {
		logx.Writer
	}
)

func NewSensitiveLogger(writer logx.Writer) *SensitiveLogger {
	return &SensitiveLogger{
		Writer: writer,
	}
}

func (l *SensitiveLogger) Info(msg any, fields ...logx.LogField) {
	if m, ok := msg.(Message); ok {
		l.Writer.Info(Message{
			Name:     m.Name,
			Password: "******",
			Message:  m.Message,
		}, fields...)
	} else {
		l.Writer.Info(msg, fields...)
	}
}

func main() {
	// setup logx to make sure originalWriter not nil,
	// the injected writer is only for filtering, like a middleware.

	originalWriter := logx.Reset()
	writer := NewSensitiveLogger(originalWriter)
	logx.SetWriter(writer)

	logx.Infov(Message{
		Name:     "foo",
		Password: "shouldNotAppear",
		Message:  "bar",
	})
  
	// more code
}
```

Complete code: [https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go)

## More examples

[https://github.com/zeromicro/zero-examples/tree/main/logx](https://github.com/zeromicro/zero-examples/tree/main/logx)

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!
