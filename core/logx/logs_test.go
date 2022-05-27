package logx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

type mockWriter struct {
	lock    sync.Mutex
	builder strings.Builder
}

func (mw *mockWriter) Alert(v interface{}) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelAlert, v)
}

func (mw *mockWriter) Error(v interface{}, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v, fields...)
}

func (mw *mockWriter) Info(v interface{}, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelInfo, v, fields...)
}

func (mw *mockWriter) Severe(v interface{}) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSevere, v)
}

func (mw *mockWriter) Slow(v interface{}, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSlow, v, fields...)
}

func (mw *mockWriter) Stack(v interface{}) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v)
}

func (mw *mockWriter) Stat(v interface{}, fields ...LogField) {
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
		want map[string]interface{}
	}{
		{
			name: "error",
			f:    Field("foo", errors.New("bar")),
			want: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "errors",
			f:    Field("foo", []error{errors.New("bar"), errors.New("baz")}),
			want: map[string]interface{}{
				"foo": []interface{}{"bar", "baz"},
			},
		},
		{
			name: "strings",
			f:    Field("foo", []string{"bar", "baz"}),
			want: map[string]interface{}{
				"foo": []interface{}{"bar", "baz"},
			},
		},
		{
			name: "duration",
			f:    Field("foo", time.Second),
			want: map[string]interface{}{
				"foo": "1s",
			},
		},
		{
			name: "durations",
			f:    Field("foo", []time.Duration{time.Second, 2 * time.Second}),
			want: map[string]interface{}{
				"foo": []interface{}{"1s", "2s"},
			},
		},
		{
			name: "times",
			f: Field("foo", []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			}),
			want: map[string]interface{}{
				"foo": []interface{}{"2020-01-01 00:00:00 +0000 UTC", "2020-01-02 00:00:00 +0000 UTC"},
			},
		},
		{
			name: "stringer",
			f:    Field("foo", ValStringer{val: "bar"}),
			want: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name: "stringers",
			f:    Field("foo", []fmt.Stringer{ValStringer{val: "bar"}, ValStringer{val: "baz"}}),
			want: map[string]interface{}{
				"foo": []interface{}{"bar", "baz"},
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

func TestStructedLogAlert(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelAlert, w, func(v ...interface{}) {
		Alert(fmt.Sprint(v...))
	})
}

func TestStructedLogError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...interface{}) {
		Error(v...)
	})
}

func TestStructedLogErrorf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...interface{}) {
		Errorf("%s", fmt.Sprint(v...))
	})
}

func TestStructedLogErrorv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...interface{}) {
		Errorv(fmt.Sprint(v...))
	})
}

func TestStructedLogErrorw(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelError, w, func(v ...interface{}) {
		Errorw(fmt.Sprint(v...), Field("foo", "bar"))
	})
}

func TestStructedLogInfo(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...interface{}) {
		Info(v...)
	})
}

func TestStructedLogInfof(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...interface{}) {
		Infof("%s", fmt.Sprint(v...))
	})
}

func TestStructedLogInfov(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...interface{}) {
		Infov(fmt.Sprint(v...))
	})
}

func TestStructedLogInfow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelInfo, w, func(v ...interface{}) {
		Infow(fmt.Sprint(v...), Field("foo", "bar"))
	})
}

func TestStructedLogInfoConsoleAny(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLogConsole(t, w, func(v ...interface{}) {
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

	doTestStructedLogConsole(t, w, func(v ...interface{}) {
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

	doTestStructedLogConsole(t, w, func(v ...interface{}) {
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

	doTestStructedLogConsole(t, w, func(v ...interface{}) {
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

	doTestStructedLogConsole(t, w, func(v ...interface{}) {
		old := atomic.LoadUint32(&encoding)
		atomic.StoreUint32(&encoding, plainEncodingType)
		defer func() {
			atomic.StoreUint32(&encoding, old)
		}()

		Info(fmt.Sprint(v...))
	})
}

func TestStructedLogSlow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...interface{}) {
		Slow(v...)
	})
}

func TestStructedLogSlowf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...interface{}) {
		Slowf(fmt.Sprint(v...))
	})
}

func TestStructedLogSlowv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...interface{}) {
		Slowv(fmt.Sprint(v...))
	})
}

func TestStructedLogSloww(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSlow, w, func(v ...interface{}) {
		Sloww(fmt.Sprint(v...), Field("foo", time.Second))
	})
}

func TestStructedLogStat(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelStat, w, func(v ...interface{}) {
		Stat(v...)
	})
}

func TestStructedLogStatf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelStat, w, func(v ...interface{}) {
		Statf(fmt.Sprint(v...))
	})
}

func TestStructedLogSevere(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSevere, w, func(v ...interface{}) {
		Severe(v...)
	})
}

func TestStructedLogSeveref(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	doTestStructedLog(t, levelSevere, w, func(v ...interface{}) {
		Severef(fmt.Sprint(v...))
	})
}

func TestStructedLogWithDuration(t *testing.T) {
	const message = "hello there"
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Info(message)
	var entry logEntry
	if err := json.Unmarshal([]byte(w.String()), &entry); err != nil {
		t.Error(err)
	}
	assert.Equal(t, levelInfo, entry.Level)
	assert.Equal(t, message, entry.Content)
	assert.Equal(t, "1000.0ms", entry.Duration)
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
		"mode",
		"console",
		"volumn",
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

	Errorf("hello %w", errors.New(message))
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

	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
	})
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "file",
		Path:        os.TempDir(),
	})
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "volume",
		Path:        os.TempDir(),
	})
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		TimeFormat:  timeFormat,
	})
	MustSetup(LogConf{
		ServiceName: "any",
		Mode:        "console",
		Encoding:    plainEncoding,
	})

	assert.NotNil(t, setupWithVolume(LogConf{}))
	assert.NotNil(t, setupWithFiles(LogConf{}))
	assert.Nil(t, setupWithFiles(LogConf{
		ServiceName: "any",
		Path:        os.TempDir(),
		Compress:    true,
		KeepDays:    1,
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

	var opt logOptions
	WithKeepDays(1)(&opt)
	WithGzip()(&opt)
	assert.Nil(t, Close())
	assert.Nil(t, Close())
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

func TestSetWriter(t *testing.T) {
	Reset()
	SetWriter(nopWriter{})
	assert.NotNil(t, writer.Load())
	assert.True(t, writer.Load() == nopWriter{})
	SetWriter(new(mockWriter))
	assert.True(t, writer.Load() == nopWriter{})
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
	fmt.Fprint(ioutil.Discard, buf)
}

func BenchmarkCopyOnWriteByteSlice(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		size := len(s)
		buf = s[:size:size]
	}
	fmt.Fprint(ioutil.Discard, buf)
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

	log.SetOutput(ioutil.Discard)
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

func doTestStructedLog(t *testing.T, level string, w *mockWriter, write func(...interface{})) {
	const message = "hello there"
	write(message)
	fmt.Println(w.String())
	var entry logEntry
	if err := json.Unmarshal([]byte(w.String()), &entry); err != nil {
		t.Error(err)
	}
	assert.Equal(t, level, entry.Level)
	val, ok := entry.Content.(string)
	assert.True(t, ok)
	assert.True(t, strings.Contains(val, message))
}

func doTestStructedLogConsole(t *testing.T, w *mockWriter, write func(...interface{})) {
	const message = "hello there"
	write(message)
	assert.True(t, strings.Contains(w.String(), message))
}

func testSetLevelTwiceWithMode(t *testing.T, mode string, w *mockWriter) {
	writer.Store(nil)
	SetUp(LogConf{
		Mode:  mode,
		Level: "error",
		Path:  "/dev/null",
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

func validateFields(t *testing.T, content string, fields map[string]interface{}) {
	var m map[string]interface{}
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
