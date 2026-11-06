<img align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# logx

[English](readme.md) | [简体中文](readme-cn.md) | 한국어

## logx 설정

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

- `ServiceName`: 서비스 이름을 설정합니다. 선택 사항입니다. `volume` 모드에서는 이 이름이 로그 파일 생성에 사용됩니다. `rest/zrpc` 서비스에서는 이름이 `rest` 또는 `zrpc`의 이름으로 자동 설정됩니다.
- `Mode`: 로그 출력 모드입니다. 기본값은 `console`입니다.
  - `console` 모드는 로그를 `stdout/stderr`에 씁니다.
  - `file` 모드는 `Path`로 지정한 파일에 로그를 씁니다.
  - `volume` 모드는 docker에서 사용하며, 마운트된 볼륨에 로그를 씁니다.
- `Encoding`: 로그 인코딩 방식을 나타냅니다. 기본값은 `json`입니다.
  - `json` 모드는 로그를 json 형식으로 씁니다.
  - `plain` 모드는 터미널 색상이 활성화된 일반 텍스트로 로그를 씁니다.
- `TimeFormat`: 시간 형식을 사용자 지정합니다. 선택 사항입니다. 기본값은 `2006-01-02T15:04:05.000Z07:00`입니다.
- `Path`: 로그 경로를 설정합니다. 기본값은 `logs`입니다.
- `Level`: 로그를 필터링할 로깅 레벨입니다. 기본값은 `info`입니다.
  - `info`: 모든 로그가 기록됩니다.
  - `error`: `info` 로그가 억제됩니다.
  - `severe`: `info`와 `error` 로그가 억제되고 `severe` 로그만 기록됩니다.
- `Compress`: 로그 파일 압축 여부입니다. `file` 모드에서만 동작합니다.
- `KeepDays`: 로그 파일을 보관할 일수입니다. 지정한 일수가 지나면 오래된 파일이 자동으로 삭제됩니다. `console` 모드에는 영향을 주지 않습니다.
- `StackCooldownMillis`: 스택 트레이스를 다시 기록하기까지의 밀리초입니다. 스택 트레이스 로그 폭주를 방지하는 데 사용됩니다.
- `MaxBackups`: 보관할 백업 로그 파일 개수입니다. 0은 모든 파일을 영구 보관한다는 의미입니다. `Rotation`이 `size`일 때만 적용됩니다. 참고: `KeepDays` 옵션의 우선순위가 더 높습니다. `MaxBackups`가 0이더라도 `KeepDays` 제한에 도달하면 로그 파일은 삭제됩니다.
- `MaxSize`: 현재 기록 중인 로그 파일이 차지할 수 있는 최대 공간입니다. 0은 제한 없음을 의미합니다. 단위는 `MB`입니다. `Rotation`이 `size`일 때만 적용됩니다.
- `Rotation`: 로그 로테이션 규칙의 유형입니다. 기본값은 `daily`입니다.
  - `daily`: 날짜 단위로 로그를 회전합니다.
  - `size`: 로그 크기 단위로 로그를 회전합니다.

## 로깅 메서드

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

- `Error`, `Info`, `Slow`: `fmt.Sprint(…)`처럼 모든 종류의 메시지를 로그에 씁니다.
- `Errorf`, `Infof`, `Slowf`: 지정한 형식으로 메시지를 로그에 씁니다.
- `Errorv`, `Infov`, `Slowv`: 모든 종류의 메시지를 json 마샬링으로 인코딩해 로그에 씁니다.
- `Errorw`, `Infow`, `Sloww`: 지정한 `key:value` 필드와 함께 문자열 메시지를 씁니다.
- `WithContext`: 지정한 ctx를 로그 메시지에 주입합니다. 일반적으로 `trace-id`와 `span-id`를 기록하는 데 사용됩니다.
- `WithDuration`: 경과 시간을 `duration` 키로 로그 메시지에 씁니다.

## 타사 로깅 라이브러리와 통합

- zap
  - 구현: [https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/zapx/zap.go)
  - 사용 예시: [https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/zaplog/main.go)
- logrus
  - 구현: [https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go](https://github.com/zeromicro/zero-contrib/blob/main/logx/logrusx/logrus.go)
  - 사용 예시: [https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/logrus/main.go)

더 많은 라이브러리는 직접 구현한 뒤 [https://github.com/zeromicro/zero-contrib](https://github.com/zeromicro/zero-contrib)에 PR을 보내주세요.

## 특정 저장소에 로그 쓰기

`logx`는 로그를 원하는 저장소에 쓸 수 있도록 사용자 지정할 수 있는 두 인터페이스를 정의합니다.

- `logx.NewWriter(w io.Writer)`
- `logx.SetWriter(writer logx.Writer)`

예를 들어 로그를 콘솔이나 파일 대신 kafka에 쓰고 싶다면 아래처럼 할 수 있습니다.

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

전체 코드: [https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/tokafka/main.go)

## 민감한 필드 필터링

`password` 필드가 로그에 기록되지 않도록 하려면 아래처럼 할 수 있습니다.

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

전체 코드: [https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go](https://github.com/zeromicro/zero-examples/blob/main/logx/filterfields/main.go)

## 더 많은 예제

[https://github.com/zeromicro/zero-examples/tree/main/logx](https://github.com/zeromicro/zero-examples/tree/main/logx)

## 별을 눌러주세요! ⭐

이 프로젝트가 마음에 들거나 학습 또는 자체 솔루션을 시작하는 데 사용 중이라면 star를 눌러주세요. 감사합니다!
