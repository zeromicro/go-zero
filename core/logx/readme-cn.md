<IMG align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# logx

[English](readme.md) | 简体中文

## logx 配置

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

- `ServiceName`：设置服务名称，可选。在 `volume` 模式下，该名称用于生成日志文件。在 `rest/zrpc` 服务中，名称将被自动设置为 `rest`或`zrpc` 的名称。
- `Mode`：输出日志的模式，默认是 `console`
  - `console` 模式将日志写到 `stdout/stderr`
  - `file` 模式将日志写到 `Path` 指定目录的文件中
  - `volume` 模式在 docker 中使用，将日志写入挂载的卷中
- `Encoding`: 指示如何对日志进行编码，默认是 `json`
  - `json`模式以 json 格式写日志
  - `plain`模式用纯文本写日志，并带有终端颜色显示
- `TimeFormat`：自定义时间格式，可选。默认是 `2006-01-02T15:04:05.000Z07:00`
- `Path`：设置日志路径，默认为 `logs`
- `Level`: 用于过滤日志的日志级别。默认为 `info`
  - `info`，所有日志都被写入
  - `error`, `info` 的日志被丢弃
  - `severe`, `info` 和 `error` 日志被丢弃，只有 `severe` 日志被写入
- `Compress`: 是否压缩日志文件，只在 `file` 模式下工作
- `KeepDays`：日志文件被保留多少天，在给定的天数之后，过期的文件将被自动删除。对 `console` 模式没有影响
- `StackCooldownMillis`：多少毫秒后再次写入堆栈跟踪。用来避免堆栈跟踪日志过多
- `MaxBackups`: 多少个日志文件备份将被保存。0代表所有备份都被保存。当`Rotation`被设置为`size`时才会起作用。注意：`KeepDays`选项的优先级会比`MaxBackups`高，即使`MaxBackups`被设置为0，当达到`KeepDays`上限时备份文件同样会被删除。
- `MaxSize`: 当前被写入的日志文件最大可占用多少空间。0代表没有上限。单位为`MB`。当`Rotation`被设置为`size`时才会起作用。
- `Rotation`: 日志轮转策略类型。默认为`daily`（按天轮转）。
  - `daily` 按天轮转。
  - `size` 按日志大小轮转。


## 打印日志方法

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

- `Error`, `Info`, `Slow`: 将任何类型的信息写进日志，使用 `fmt.Sprint(...)` 来转换为 `string`
- `Errorf`, `Infof`, `Slowf`: 将指定格式的信息写入日志
- `Errorv`, `Infov`, `Slowv`: 将任何类型的信息写入日志，用 `json marshal` 编码
- `Errorw`, `Infow`, `Sloww`: 写日志，并带上给定的 `key:value` 字段
- `WithContext`：将给定的 ctx 注入日志信息，例如用于记录 `trace-id`和`span-id`
- `WithDuration`: 将指定的时间写入日志信息中，字段名为 `duration`

## 与第三方日志库集成

- zap
  - 实现：[https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go)
  - 使用示例：[https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go)
- logrus
  - 实现：[https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go)
  - 使用示例：[https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go)

对于其它的日志库，请参考上面示例实现，并欢迎提交 `PR` 到 [https://github.com/zeromicro/zero-contrib](https://github.com/zeromicro/zero-contrib)

## 将日志写到指定的存储

`logx`定义了两个接口，方便自定义 `logx`，将日志写入任何存储。

- `logx.NewWriter(w io.Writer)`
- `logx.SetWriter(write logx.Writer)`

例如，如果我们想把日志写进kafka，而不是控制台或文件，我们可以像下面这样做。

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

完整代码：[https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go)

## 过滤敏感字段

如果我们需要防止  `password` 字段被记录下来，我们可以像下面这样实现。

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

完整代码：[https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go)

## 更多示例

[https://github.com/zeromicro/zero-examples/tree/main/logx](https://github.com/zeromicro/zero-examples/tree/main/logx)

## Give a Star! ⭐

如果你正在使用或者觉得这个项目对你有帮助，请 **star** 支持，感谢！
