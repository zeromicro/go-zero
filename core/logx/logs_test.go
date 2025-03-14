package logx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	s           = []byte("Sending #11 notification (id: 1451875113812010473) in #1 connection")
	pool        = make(chan []byte, 1)
	_    Writer = (*mockWriter)(nil)
)

func init() {
	ExitOnFatal.Set(false)
}

type mockWriter struct {
	lock    sync.Mutex
	builder strings.Builder
}

func (mw *mockWriter) Alert(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelAlert, v)
}

func (mw *mockWriter) Debug(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelDebug, v, fields...)
}

func (mw *mockWriter) Error(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v, fields...)
}

func (mw *mockWriter) Info(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelInfo, v, fields...)
}

func (mw *mockWriter) Severe(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSevere, v)
}

func (mw *mockWriter) Slow(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSlow, v, fields...)
}

func (mw *mockWriter) Stack(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v)
}

func (mw *mockWriter) Stat(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelStat, v, fields...)
}

func (mw *mockWriter) Close() error {
	return nil
}

func (mw *mockWriter) Contains(text string) bool {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return strings.Contains(mw.builder.String(), text)
}

func (mw *mockWriter) Reset() {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.builder.Reset()
}

func (mw *mockWriter) String() string {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return mw.builder.String()
}

func TestField(t *testing.T) {
	tests := []struct {
		name string
		f    LogField
		want map[string]any
	}{
		{
			name: "error",
			f:    Field("foo", errors.New("bar")),
			want: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "errors",
			f:    Field("foo", []error{errors.New("bar"), errors.New("baz")}),
			want: map[string]any{
				"foo": []any{"bar", "baz"},
			},
		},
		{
			name: "strings",
			f:    Field("foo", []string{"bar", "baz"}),
			want: map[string]any{
				"foo": []any{"bar", "baz"},
			},
		},
		{
			name: "duration",
			f:    Field("foo", time.Second),
			want: map[string]any{
				"foo": "1s",
			},
		},
		{
			name: "durations",
			f:    Field("foo", []time.Duration{time.Second, 2 * time.Second}),
			want: map[string]any{
				"foo": []any{"1s", "2s"},
			},
		},
		{
			name: "times",
			f: Field("foo", []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			}),
			want: map[string]any{
				"foo": []any{"2020-01-01 00:00:00 +0000 UTC", "2020-01-02 00:00:00 +0000 UTC"},
			},
		},
		{
			name: "stringer",
			f:    Field("foo", ValStringer{val: "bar"}),
			want: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "stringers",
			f:    Field("foo", []fmt.Stringer{ValStringer{val: "bar"}, ValStringer{val: "baz"}}),
			want: map[string]any{
				"foo": []any{"bar", "baz"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			w := new(mockWriter)
			old := writer.Swap(w)
			defer writer.Store(old)

			Infow("foo", test.f)
			validateFields(t, w.String(), test.want)
		})
	}
}

func TestFileLineFileMode(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	file, line := getFileLine()
	Error("anything")
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))

	file, line = getFileLine()
	Errorf("anything %s", "format")
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))
}

func TestFileLineConsoleMode(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	file, line := getFileLine()
	Error("anything")
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))

	w.Reset()
	file, line = getFileLine()
	Errorf("anything %s", "format")
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))
}

func TestMust(t *testing.T) {
	assert.Panics(t, func() {
		Must(errors.New("foo"))
	})
}

func TestStructedLogAlert(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelAlert, w, func(v ...any) {
		Alert(fmt.Sprint(v...))
	})
}

func TestStructedLogDebug(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelDebug, w, func(v ...any) {
		Debug(v...)
	})
}

func TestStructedLogDebugf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelDebug, w, func(v ...any) {
		Debugf(fmt.Sprint(v...))
	})
}

func TestStructedLogDebugfn(t *testing.T) {
	t.Run("debugfn with output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLog(t, levelDebug, w, func(v ...any) {
			Debugfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})

	t.Run("debugfn without output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLogEmpty(t, w, InfoLevel, func(v ...any) {
			Debugfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})
}

func TestStructedLogDebugv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelDebug, w, func(v ...any) {
		Debugv(fmt.Sprint(v...))
	})
}

func TestStructedLogDebugw(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelDebug, w, func(v ...any) {
		Debugw(fmt.Sprint(v...), Field("foo", time.Second))
	})
}

func TestStructedLogError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...any) {
		Error(v...)
	})
}

func TestStructedLogErrorf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...any) {
		Errorf("%s", fmt.Sprint(v...))
	})
}

func TestStructedLogErrorfn(t *testing.T) {
	t.Run("errorfn with output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLog(t, levelError, w, func(v ...any) {
			Errorfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})

	t.Run("errorfn without output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLogEmpty(t, w, SevereLevel, func(v ...any) {
			Errorfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})
}

func TestStructedLogErrorv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...any) {
		Errorv(fmt.Sprint(v...))
	})
}

func TestStructedLogErrorw(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...any) {
		Errorw(fmt.Sprint(v...), Field("foo", "bar"))
	})
}

func TestStructedLogInfo(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...any) {
		Info(v...)
	})
}

func TestStructedLogInfof(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...any) {
		Infof("%s", fmt.Sprint(v...))
	})
}

func TestStructedInfofn(t *testing.T) {
	t.Run("infofn with output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLog(t, levelInfo, w, func(v ...any) {
			Infofn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})

	t.Run("infofn without output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLogEmpty(t, w, ErrorLevel, func(v ...any) {
			Infofn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})
}

func TestStructedLogInfov(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...any) {
		Infov(fmt.Sprint(v...))
	})
}

func TestStructedLogInfow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...any) {
		Infow(fmt.Sprint(v...), Field("foo", "bar"))
	})
}

func TestStructedLogFieldNil(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	assert.NotPanics(t, func() {
		var s *string
		Infow("test", Field("bb", s))
		var d *nilStringer
		Infow("test", Field("bb", d))
		var e *nilError
		Errorw("test", Field("bb", e))
	})
	assert.NotPanics(t, func() {
		var p panicStringer
		Infow("test", Field("bb", p))
		var ps innerPanicStringer
		Infow("test", Field("bb", ps))
	})
}

func TestStructedLogInfoConsoleAny(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...any) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Infov(v)
	})
}

func TestStructedLogInfoConsoleAnyString(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...any) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Infov(fmt.Sprint(v...))
	})
}

func TestStructedLogInfoConsoleAnyError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...any) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Infov(errors.New(fmt.Sprint(v...)))
	})
}

func TestStructedLogInfoConsoleAnyStringer(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...any) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Infov(ValStringer{
			val: fmt.Sprint(v...),
		})
	})
}

func TestStructedLogInfoConsoleText(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...any) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Info(fmt.Sprint(v...))
	})
}

func TestInfofnWithErrorLevel(t *testing.T) {
	called := false
	SetLevel(ErrorLevel)
	defer SetLevel(DebugLevel)
	Infofn(func() any {
		called = true
		return "info log"
	})
	assert.False(t, called)
}

func TestStructedLogSlow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...any) {
		Slow(v...)
	})
}

func TestStructedLogSlowf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...any) {
		Slowf(fmt.Sprint(v...))
	})
}

func TestStructedLogSlowfn(t *testing.T) {
	t.Run("slowfn with output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLog(t, levelSlow, w, func(v ...any) {
			Slowfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})

	t.Run("slowfn without output", func(t *testing.T) {
		w := new(mockWriter)
		old := writer.Swap(w)
		defer writer.Store(old)

		doTestStructedLogEmpty(t, w, SevereLevel, func(v ...any) {
			Slowfn(func() any {
				return fmt.Sprint(v...)
			})
		})
	})
}

func TestStructedLogSlowv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...any) {
		Slowv(fmt.Sprint(v...))
	})
}

func TestStructedLogSloww(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...any) {
		Sloww(fmt.Sprint(v...), Field("foo", time.Second))
	})
}

func TestStructedLogStat(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelStat, w, func(v ...any) {
		Stat(v...)
	})
}

func TestStructedLogStatf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelStat, w, func(v ...any) {
		Statf(fmt.Sprint(v...))
	})
}

func TestStructedLogSevere(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSevere, w, func(v ...any) {
		Severe(v...)
	})
}

func TestStructedLogSeveref(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSevere, w, func(v ...any) {
		Severef(fmt.Sprint(v...))
	})
}

func TestStructedLogWithDuration(t *testing.T) {
	const message = "hello there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Info(message)
	var entry map[string]any
	if err := json.Unmarshal([]byte(w.String()), &entry); err != nil {
		t.Error(err)
	}
	assert.Equal(t, levelInfo, entry[levelKey])
	assert.Equal(t, message, entry[contentKey])
	assert.Equal(t, "1000.0ms", entry[durationKey])
}

func TestSetLevel(t *testing.T) {
	SetLevel(ErrorLevel)
	const message = "hello there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	Info(message)
	assert.Equal(t, 0, w.builder.Len())
}

func TestSetLevelTwiceWithMode(t *testing.T) {
	testModes := []string{
		"console",
		"volumn",
		"mode",
	}
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	for _, mode := range testModes {
		testSetLevelTwiceWithMode(t, mode, w)
	}
}

func TestSetLevelWithDuration(t *testing.T) {
	SetLevel(ErrorLevel)
	const message = "hello there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Info(message)
	assert.Equal(t, 0, w.builder.Len())
}

func TestErrorfWithWrappedError(t *testing.T) {
	SetLevel(ErrorLevel)
	const message = "there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	Errorf("hello %s", errors.New(message))
	assert.True(t, strings.Contains(w.String(), "hello there"))
}

func TestMustNil(t *testing.T) {
	Must(nil)
}

func TestSetup(t *testing.T) {
	defer func() {
		SetLevel(InfoLevel)
		atomic.StoreUint32(&encoding, jsonEncodingType)
	}()

	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		Encoding:    "json",
		TimeFormat:  timeFormat,
	})
	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		TimeFormat:  timeFormat,
	})
	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "file",
		Path:        os.TempDir(),
	})
	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "volume",
		Path:        os.TempDir(),
	})
	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		TimeFormat:  timeFormat,
	})
	setupOnce = sync.Once{}
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		Encoding:    plainEncoding,
	})

	defer os.RemoveAll("CD01CB7D-2705-4F3F-889E-86219BF56F10")
	assert.NotNil(t, setupWithVolume(LogConf{}))
	assert.Nil(t, setupWithVolume(LogConf{
		ServiceName: "CD01CB7D-2705-4F3F-889E-86219BF56F10",
	}))
	assert.Nil(t, setupWithVolume(LogConf{
		ServiceName: "CD01CB7D-2705-4F3F-889E-86219BF56F10",
		Rotation:    sizeRotationRule,
	}))
	assert.NotNil(t, setupWithFiles(LogConf{}))
	assert.Nil(t, setupWithFiles(LogConf{
		ServiceName: "any",
		Path:        os.TempDir(),
		Compress:    true,
		KeepDays:    1,
		MaxBackups:  3,
		MaxSize:     1024 * 1024,
	}))
	setupLogLevel(LogConf{
		Level: levelInfo,
	})
	setupLogLevel(LogConf{
		Level: levelError,
	})
	setupLogLevel(LogConf{
		Level: levelSevere,
	})
	_, err := createOutput("")
	assert.NotNil(t, err)
	Disable()
	SetLevel(InfoLevel)
	atomic.StoreUint32(&encoding, jsonEncodingType)
}

func TestDisable(t *testing.T) {
	Disable()
	defer func() {
		SetLevel(InfoLevel)
		atomic.StoreUint32(&encoding, jsonEncodingType)
	}()

	var opt logOptions
	WithKeepDays(1)(&opt)
	WithGzip()(&opt)
	WithMaxBackups(1)(&opt)
	WithMaxSize(1024)(&opt)
	assert.Nil(t, Close())
	assert.Nil(t, Close())
	assert.Equal(t, uint32(disableLevel), atomic.LoadUint32(&logLevel))
}

func TestDisableStat(t *testing.T) {
	DisableStat()

	const message = "hello there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)
	Stat(message)
	assert.Equal(t, 0, w.builder.Len())
}

func TestAddWriter(t *testing.T) {
	const message = "hello there"
	w := new(mockWriter)
	AddWriter(w)
	w1 := new(mockWriter)
	AddWriter(w1)
	Error(message)
	assert.Contains(t, w.String(), message)
	assert.Contains(t, w1.String(), message)
}

func TestSetWriter(t *testing.T) {
	atomic.StoreUint32(&logLevel, 0)
	Reset()
	SetWriter(nopWriter{})
	assert.NotNil(t, writer.Load())
	assert.True(t, writer.Load() == nopWriter{})
	mocked := new(mockWriter)
	SetWriter(mocked)
	assert.Equal(t, mocked, writer.Load())
}

func TestWithGzip(t *testing.T) {
	fn := WithGzip()
	var opt logOptions
	fn(&opt)
	assert.True(t, opt.gzipEnabled)
}

func TestWithKeepDays(t *testing.T) {
	fn := WithKeepDays(1)
	var opt logOptions
	fn(&opt)
	assert.Equal(t, 1, opt.keepDays)
}

func BenchmarkCopyByteSliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf []byte
		buf = append(buf, getTimestamp()...)
		buf = append(buf, ' ')
		buf = append(buf, s...)
		_ = buf
	}
}

func BenchmarkCopyByteSliceAllocExactly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now := []byte(getTimestamp())
		buf := make([]byte, len(now)+1+len(s))
		n := copy(buf, now)
		buf[n] = ' '
		copy(buf[n+1:], s)
	}
}

func BenchmarkCopyByteSlice(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf = make([]byte, len(s))
		copy(buf, s)
	}
	fmt.Fprint(io.Discard, buf)
}

func BenchmarkCopyOnWriteByteSlice(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		size := len(s)
		buf = s[:size:size]
	}
	fmt.Fprint(io.Discard, buf)
}

func BenchmarkCacheByteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dup := fetch()
		copy(dup, s)
		put(dup)
	}
}

func BenchmarkLogs(b *testing.B) {
	b.ReportAllocs()

	log.SetOutput(io.Discard)
	for i := 0; i < b.N; i++ {
		Info(i)
	}
}

func fetch() []byte {
	select {
	case b := <-pool:
		return b
	default:
	}
	return make([]byte, 4096)
}

func getFileLine() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	return short, line
}

func put(b []byte) {
	select {
	case pool <- b:
	default:
	}
}

func doTestStructedLog(t *testing.T, level string, w *mockWriter, write func(...any)) {
	const message = "hello there"
	write(message)

	var entry map[string]any
	if err := json.Unmarshal([]byte(w.String()), &entry); err != nil {
		t.Error(err)
	}

	assert.Equal(t, level, entry[levelKey])
	val, ok := entry[contentKey]
	assert.True(t, ok)
	assert.True(t, strings.Contains(val.(string), message))
}

func doTestStructedLogConsole(t *testing.T, w *mockWriter, write func(...any)) {
	const message = "hello there"
	write(message)
	assert.True(t, strings.Contains(w.String(), message))
}

func doTestStructedLogEmpty(t *testing.T, w *mockWriter, level uint32, write func(...any)) {
	olevel := atomic.LoadUint32(&logLevel)
	SetLevel(level)
	defer SetLevel(olevel)

	const message = "hello there"
	write(message)
	assert.Empty(t, w.String())
}

func testSetLevelTwiceWithMode(t *testing.T, mode string, w *mockWriter) {
	writer.Store(nil)
	SetUp(LogConf{
		Mode:           mode,
		Level:          "debug",
		Path:           "/dev/null",
		Encoding:       plainEncoding,
		Stat:           false,
		TimeFormat:     time.RFC3339,
		FileTimeFormat: time.DateTime,
	})
	SetUp(LogConf{
		Mode:  mode,
		Level: "info",
		Path:  "/dev/null",
	})
	const message = "hello there"
	Info(message)
	assert.Equal(t, 0, w.builder.Len())
	Infof(message)
	assert.Equal(t, 0, w.builder.Len())
	ErrorStack(message)
	assert.Equal(t, 0, w.builder.Len())
	ErrorStackf(message)
	assert.Equal(t, 0, w.builder.Len())
}

type ValStringer struct {
	val string
}

func (v ValStringer) String() string {
	return v.val
}

func validateFields(t *testing.T, content string, fields map[string]any) {
	var m map[string]any
	if err := json.Unmarshal([]byte(content), &m); err != nil {
		t.Error(err)
	}

	for k, v := range fields {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			assert.EqualValues(t, v, m[k])
		} else {
			assert.Equal(t, v, m[k], content)
		}
	}
}

type nilError struct {
	Name string
}

func (e *nilError) Error() string {
	return e.Name
}

type nilStringer struct {
	Name string
}

func (s *nilStringer) String() string {
	return s.Name
}

type innerPanicStringer struct {
	Inner *struct {
		Name string
	}
}

func (s innerPanicStringer) String() string {
	return s.Inner.Name
}

type panicStringer struct {
}

func (s panicStringer) String() string {
	panic("panic")
}
