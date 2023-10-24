package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/sysx"
)

const callerDepth = 4

var (
	timeFormat = "2006-01-02T15:04:05.000Z07:00"
	logLevel   uint32
	encoding   uint32 = jsonEncodingType
	// maxContentLength is used to truncate the log content, 0 for not truncating.
	maxContentLength uint32
	// use uint32 for atomic operations
	disableLog  uint32
	disableStat uint32
	options     logOptions
	writer      = new(atomicWriter)
	setupOnce   sync.Once
)

type (
	// LogField is a key-value pair that will be added to the log entry.
	LogField struct {
		Key   string
		Value any
	}

	// LogOption defines the method to customize the logging.
	LogOption func(options *logOptions)

	logEntry map[string]any

	logOptions struct {
		gzipEnabled           bool
		logStackCooldownMills int
		keepDays              int
		maxBackups            int
		maxSize               int
		rotationRule          string
	}
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

// Debug writes v into access log.
func Debug(v ...any) {
	if shallLog(DebugLevel) {
		writeDebug(fmt.Sprint(v...))
	}
}

// Debugf writes v with format into access log.
func Debugf(format string, v ...any) {
	if shallLog(DebugLevel) {
		writeDebug(fmt.Sprintf(format, v...))
	}
}

// Debugv writes v into access log with json content.
func Debugv(v any) {
	if shallLog(DebugLevel) {
		writeDebug(v)
	}
}

// Debugw writes msg along with fields into access log.
func Debugw(msg string, fields ...LogField) {
	if shallLog(DebugLevel) {
		writeDebug(msg, fields...)
	}
}

// Disable disables the logging.
func Disable() {
	atomic.StoreUint32(&disableLog, 1)
	writer.Store(nopWriter{})
}

// DisableStat disables the stat logs.
func DisableStat() {
	atomic.StoreUint32(&disableStat, 1)
}

// Error writes v into error log.
func Error(v ...any) {
	if shallLog(ErrorLevel) {
		writeError(fmt.Sprint(v...))
	}
}

// Errorf writes v with format into error log.
func Errorf(format string, v ...any) {
	if shallLog(ErrorLevel) {
		writeError(fmt.Errorf(format, v...).Error())
	}
}

// ErrorStack writes v along with call stack into error log.
func ErrorStack(v ...any) {
	if shallLog(ErrorLevel) {
		// there is newline in stack string
		writeStack(fmt.Sprint(v...))
	}
}

// ErrorStackf writes v along with call stack in format into error log.
func ErrorStackf(format string, v ...any) {
	if shallLog(ErrorLevel) {
		// there is newline in stack string
		writeStack(fmt.Sprintf(format, v...))
	}
}

// Errorv writes v into error log with json content.
// No call stack attached, because not elegant to pack the messages.
func Errorv(v any) {
	if shallLog(ErrorLevel) {
		writeError(v)
	}
}

// Errorw writes msg along with fields into error log.
func Errorw(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		writeError(msg, fields...)
	}
}

// Field returns a LogField for the given key and value.
func Field(key string, value any) LogField {
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
func Info(v ...any) {
	if shallLog(InfoLevel) {
		writeInfo(fmt.Sprint(v...))
	}
}

// Infof writes v with format into access log.
func Infof(format string, v ...any) {
	if shallLog(InfoLevel) {
		writeInfo(fmt.Sprintf(format, v...))
	}
}

// Infov writes v into access log with json content.
func Infov(v any) {
	if shallLog(InfoLevel) {
		writeInfo(v)
	}
}

// Infow writes msg along with fields into access log.
func Infow(msg string, fields ...LogField) {
	if shallLog(InfoLevel) {
		writeInfo(msg, fields...)
	}
}

// Must checks if err is nil, otherwise logs the error and exits.
func Must(err error) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf("%+v\n\n%s", err.Error(), debug.Stack())
	log.Print(msg)
	getWriter().Severe(msg)

	if ExitOnFatal.True() {
		os.Exit(1)
	} else {
		panic(msg)
	}
}

// MustSetup sets up logging with given config c. It exits on error.
func MustSetup(c LogConf) {
	Must(SetUp(c))
}

// Reset clears the writer and resets the log level.
func Reset() Writer {
	return writer.Swap(nil)
}

// SetLevel sets the logging level. It can be used to suppress some logs.
func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
}

// SetWriter sets the logging writer. It can be used to customize the logging.
func SetWriter(w Writer) {
	if atomic.LoadUint32(&disableLog) == 0 {
		writer.Store(w)
	}
}

// SetUp sets up the logx. If already set up, just return nil.
// we allow SetUp to be called multiple times, because for example
// we need to allow different service frameworks to initialize logx respectively.
func SetUp(c LogConf) (err error) {
	// Just ignore the subsequent SetUp calls.
	// Because multiple services in one process might call SetUp respectively.
	// Need to wait for the first caller to complete the execution.
	setupOnce.Do(func() {
		setupLogLevel(c)

		if !c.Stat {
			DisableStat()
		}

		if len(c.TimeFormat) > 0 {
			timeFormat = c.TimeFormat
		}

		atomic.StoreUint32(&maxContentLength, c.MaxContentLength)

		switch c.Encoding {
		case plainEncoding:
			atomic.StoreUint32(&encoding, plainEncodingType)
		default:
			atomic.StoreUint32(&encoding, jsonEncodingType)
		}

		switch c.Mode {
		case fileMode:
			err = setupWithFiles(c)
		case volumeMode:
			err = setupWithVolume(c)
		default:
			setupWithConsole()
		}
	})

	return
}

// Severe writes v into severe log.
func Severe(v ...any) {
	if shallLog(SevereLevel) {
		writeSevere(fmt.Sprint(v...))
	}
}

// Severef writes v with format into severe log.
func Severef(format string, v ...any) {
	if shallLog(SevereLevel) {
		writeSevere(fmt.Sprintf(format, v...))
	}
}

// Slow writes v into slow log.
func Slow(v ...any) {
	if shallLog(ErrorLevel) {
		writeSlow(fmt.Sprint(v...))
	}
}

// Slowf writes v with format into slow log.
func Slowf(format string, v ...any) {
	if shallLog(ErrorLevel) {
		writeSlow(fmt.Sprintf(format, v...))
	}
}

// Slowv writes v into slow log with json content.
func Slowv(v any) {
	if shallLog(ErrorLevel) {
		writeSlow(v)
	}
}

// Sloww writes msg along with fields into slow log.
func Sloww(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		writeSlow(msg, fields...)
	}
}

// Stat writes v into stat log.
func Stat(v ...any) {
	if shallLogStat() && shallLog(InfoLevel) {
		writeStat(fmt.Sprint(v...))
	}
}

// Statf writes v with format into stat log.
func Statf(format string, v ...any) {
	if shallLogStat() && shallLog(InfoLevel) {
		writeStat(fmt.Sprintf(format, v...))
	}
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

// WithMaxBackups customizes how many log files backups will be kept.
func WithMaxBackups(count int) LogOption {
	return func(opts *logOptions) {
		opts.maxBackups = count
	}
}

// WithMaxSize customizes how much space the writing log file can take up.
func WithMaxSize(size int) LogOption {
	return func(opts *logOptions) {
		opts.maxSize = size
	}
}

// WithRotation customizes which log rotation rule to use.
func WithRotation(r string) LogOption {
	return func(opts *logOptions) {
		opts.rotationRule = r
	}
}

func addCaller(fields ...LogField) []LogField {
	return append(fields, Field(callerKey, getCaller(callerDepth)))
}

func createOutput(path string) (io.WriteCloser, error) {
	if len(path) == 0 {
		return nil, ErrLogPathNotSet
	}

	var rule RotateRule
	switch options.rotationRule {
	case sizeRotationRule:
		rule = NewSizeLimitRotateRule(path, backupFileDelimiter, options.keepDays, options.maxSize,
			options.maxBackups, options.gzipEnabled)
	default:
		rule = DefaultRotateRule(path, backupFileDelimiter, options.keepDays, options.gzipEnabled)
	}

	return NewLogger(path, rule, options.gzipEnabled)
}

func getWriter() Writer {
	w := writer.Load()
	if w == nil {
		w = writer.StoreIfNil(newConsoleWriter())
	}

	return w
}

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}

func setupLogLevel(c LogConf) {
	switch c.Level {
	case levelDebug:
		SetLevel(DebugLevel)
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

func shallLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

func shallLogStat() bool {
	return atomic.LoadUint32(&disableStat) == 0
}

// writeDebug writes v into debug log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeDebug(val any, fields ...LogField) {
	getWriter().Debug(val, addCaller(fields...)...)
}

// writeError writes v into error log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeError(val any, fields ...LogField) {
	getWriter().Error(val, addCaller(fields...)...)
}

// writeInfo writes v into info log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeInfo(val any, fields ...LogField) {
	getWriter().Info(val, addCaller(fields...)...)
}

// writeSevere writes v into severe log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeSevere(msg string) {
	getWriter().Severe(fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
}

// writeSlow writes v into slow log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeSlow(val any, fields ...LogField) {
	getWriter().Slow(val, addCaller(fields...)...)
}

// writeStack writes v into stack log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeStack(msg string) {
	getWriter().Stack(fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
}

// writeStat writes v into stat log.
// Not checking shallLog here is for performance consideration.
// If we check shallLog here, the fmt.Sprint might be called even if the log level is not enabled.
// The caller should check shallLog before calling this function.
func writeStat(msg string) {
	getWriter().Stat(msg, addCaller()...)
}
