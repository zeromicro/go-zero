package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/sysx"
)

const callerDepth = 5

var (
	timeFormat = "2006-01-02T15:04:05.000Z07:00"
	logLevel   uint32
	encoding   uint32 = jsonEncodingType
	// use uint32 for atomic operations
	disableStat uint32

	options logOptions
	writer  = new(atomicWriter)
)

type (
	logEntry struct {
		Timestamp string      `json:"@timestamp"`
		Level     string      `json:"level"`
		Duration  string      `json:"duration,omitempty"`
		Caller    string      `json:"caller,omitempty"`
		Content   interface{} `json:"content"`
	}

	logEntryWithFields map[string]interface{}

	logOptions struct {
		gzipEnabled           bool
		logStackCooldownMills int
		keepDays              int
	}

	// LogField is a key-value pair that will be added to the log entry.
	LogField struct {
		Key   string
		Value interface{}
	}

	// LogOption defines the method to customize the logging.
	LogOption func(options *logOptions)
)

// Alert alerts v in alert level, and the message is written to error log.
func Alert(v string) {
	getWriter().Alert(v)
}

// Close closes the logging.
func Close() error {
	if w := writer.Swap(nil); w != nil {
		return w.(io.Closer).Close()
	}

	return nil
}

// Disable disables the logging.
func Disable() {
	writer.Store(nopWriter{})
}

// DisableStat disables the stat logs.
func DisableStat() {
	atomic.StoreUint32(&disableStat, 1)
}

// Error writes v into error log.
func Error(v ...interface{}) {
	errorTextSync(fmt.Sprint(v...))
}

// Errorf writes v with format into error log.
func Errorf(format string, v ...interface{}) {
	errorTextSync(fmt.Errorf(format, v...).Error())
}

// ErrorStack writes v along with call stack into error log.
func ErrorStack(v ...interface{}) {
	// there is newline in stack string
	stackSync(fmt.Sprint(v...))
}

// ErrorStackf writes v along with call stack in format into error log.
func ErrorStackf(format string, v ...interface{}) {
	// there is newline in stack string
	stackSync(fmt.Sprintf(format, v...))
}

// Errorv writes v into error log with json content.
// No call stack attached, because not elegant to pack the messages.
func Errorv(v interface{}) {
	errorAnySync(v)
}

// Errorw writes msg along with fields into error log.
func Errorw(msg string, fields ...LogField) {
	errorFieldsSync(msg, fields...)
}

// Field returns a LogField for the given key and value.
func Field(key string, value interface{}) LogField {
	switch val := value.(type) {
	case error:
		return LogField{Key: key, Value: val.Error()}
	case []error:
		var errs []string
		for _, err := range val {
			errs = append(errs, err.Error())
		}
		return LogField{Key: key, Value: errs}
	case time.Duration:
		return LogField{Key: key, Value: fmt.Sprint(val)}
	case []time.Duration:
		var durs []string
		for _, dur := range val {
			durs = append(durs, fmt.Sprint(dur))
		}
		return LogField{Key: key, Value: durs}
	case []time.Time:
		var times []string
		for _, t := range val {
			times = append(times, fmt.Sprint(t))
		}
		return LogField{Key: key, Value: times}
	case fmt.Stringer:
		return LogField{Key: key, Value: val.String()}
	case []fmt.Stringer:
		var strs []string
		for _, str := range val {
			strs = append(strs, str.String())
		}
		return LogField{Key: key, Value: strs}
	default:
		return LogField{Key: key, Value: val}
	}
}

// Info writes v into access log.
func Info(v ...interface{}) {
	infoTextSync(fmt.Sprint(v...))
}

// Infof writes v with format into access log.
func Infof(format string, v ...interface{}) {
	infoTextSync(fmt.Sprintf(format, v...))
}

// Infov writes v into access log with json content.
func Infov(v interface{}) {
	infoAnySync(v)
}

// Infow writes msg along with fields into access log.
func Infow(msg string, fields ...LogField) {
	infoFieldsSync(msg, fields...)
}

// Must checks if err is nil, otherwise logs the error and exits.
func Must(err error) {
	if err == nil {
		return
	}

	msg := err.Error()
	log.Print(msg)
	getWriter().Severe(msg)
	os.Exit(1)
}

// MustSetup sets up logging with given config c. It exits on error.
func MustSetup(c LogConf) {
	Must(SetUp(c))
}

// Reset clears the writer and resets the log level.
func Reset() Writer {
	SetLevel(InfoLevel)
	return writer.Swap(nil)
}

// SetLevel sets the logging level. It can be used to suppress some logs.
func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
}

// SetWriter sets the logging writer. It can be used to customize the logging.
// Call Reset before calling SetWriter again.
func SetWriter(w Writer) {
	if writer.Load() == nil {
		writer.Store(w)
	}
}

// SetUp sets up the logx. If already set up, just return nil.
// we allow SetUp to be called multiple times, because for example
// we need to allow different service frameworks to initialize logx respectively.
// the same logic for SetUp
func SetUp(c LogConf) error {
	setupLogLevel(c)

	if len(c.TimeFormat) > 0 {
		timeFormat = c.TimeFormat
	}

	switch c.Encoding {
	case plainEncoding:
		atomic.StoreUint32(&encoding, plainEncodingType)
	default:
		atomic.StoreUint32(&encoding, jsonEncodingType)
	}

	switch c.Mode {
	case fileMode:
		return setupWithFiles(c)
	case volumeMode:
		return setupWithVolume(c)
	default:
		setupWithConsole()
		return nil
	}
}

// Severe writes v into severe log.
func Severe(v ...interface{}) {
	severeSync(fmt.Sprint(v...))
}

// Severef writes v with format into severe log.
func Severef(format string, v ...interface{}) {
	severeSync(fmt.Sprintf(format, v...))
}

// Slow writes v into slow log.
func Slow(v ...interface{}) {
	slowTextSync(fmt.Sprint(v...))
}

// Slowf writes v with format into slow log.
func Slowf(format string, v ...interface{}) {
	slowTextSync(fmt.Sprintf(format, v...))
}

// Slowv writes v into slow log with json content.
func Slowv(v interface{}) {
	slowAnySync(v)
}

// Sloww writes msg along with fields into slow log.
func Sloww(msg string, fields ...LogField) {
	slowFieldsSync(msg, fields...)
}

// Stat writes v into stat log.
func Stat(v ...interface{}) {
	statSync(fmt.Sprint(v...))
}

// Statf writes v with format into stat log.
func Statf(format string, v ...interface{}) {
	statSync(fmt.Sprintf(format, v...))
}

// WithCooldownMillis customizes logging on writing call stack interval.
func WithCooldownMillis(millis int) LogOption {
	return func(opts *logOptions) {
		opts.logStackCooldownMills = millis
	}
}

// WithKeepDays customizes logging to keep logs with days.
func WithKeepDays(days int) LogOption {
	return func(opts *logOptions) {
		opts.keepDays = days
	}
}

// WithGzip customizes logging to automatically gzip the log files.
func WithGzip() LogOption {
	return func(opts *logOptions) {
		opts.gzipEnabled = true
	}
}

func createOutput(path string) (io.WriteCloser, error) {
	if len(path) == 0 {
		return nil, ErrLogPathNotSet
	}

	return NewLogger(path, DefaultRotateRule(path, backupFileDelimiter, options.keepDays,
		options.gzipEnabled), options.gzipEnabled)
}

func errorAnySync(v interface{}) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v)
	}
}

func errorFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(content, fields...)
	}
}

func errorTextSync(msg string) {
	if shallLog(ErrorLevel) {
		getWriter().Error(msg)
	}
}

func getWriter() Writer {
	w := writer.Load()
	if w == nil {
		w = newConsoleWriter()
		writer.Store(w)
	}

	return w
}

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}

func infoAnySync(val interface{}) {
	if shallLog(InfoLevel) {
		getWriter().Info(val)
	}
}

func infoFieldsSync(content string, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(content, fields...)
	}
}

func infoTextSync(msg string) {
	if shallLog(InfoLevel) {
		getWriter().Info(msg)
	}
}

func setupLogLevel(c LogConf) {
	switch c.Level {
	case levelInfo:
		SetLevel(InfoLevel)
	case levelError:
		SetLevel(ErrorLevel)
	case levelSevere:
		SetLevel(SevereLevel)
	}
}

func setupWithConsole() {
	SetWriter(newConsoleWriter())
}

func setupWithFiles(c LogConf) error {
	w, err := newFileWriter(c)
	if err != nil {
		return err
	}

	SetWriter(w)
	return nil
}

func setupWithVolume(c LogConf) error {
	if len(c.ServiceName) == 0 {
		return ErrLogServiceNameNotSet
	}

	c.Path = path.Join(c.Path, c.ServiceName, sysx.Hostname())
	return setupWithFiles(c)
}

func severeSync(msg string) {
	if shallLog(SevereLevel) {
		getWriter().Severe(fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func shallLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

func shallLogStat() bool {
	return atomic.LoadUint32(&disableStat) == 0
}

func slowAnySync(v interface{}) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(v)
	}
}

func slowFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(content, fields...)
	}
}

func slowTextSync(msg string) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(msg)
	}
}

func stackSync(msg string) {
	if shallLog(ErrorLevel) {
		getWriter().Stack(fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func statSync(msg string) {
	if shallLogStat() && shallLog(InfoLevel) {
		getWriter().Stat(msg)
	}
}
