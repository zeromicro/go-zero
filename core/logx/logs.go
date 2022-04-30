package logx

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/sysx"
)

const callerDepth = 4

var (
	timeFormat   = "2006-01-02T15:04:05.000Z07:00"
	writeConsole bool
	logLevel     uint32
	encoding     uint32 = jsonEncodingType
	// use uint32 for atomic operations
	disableStat uint32
	infoLog     io.WriteCloser
	errorLog    io.WriteCloser
	severeLog   io.WriteCloser
	slowLog     io.WriteCloser
	statLog     io.WriteCloser
	stackLog    io.Writer

	once        sync.Once
	initialized uint32
	options     logOptions
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

// MustSetup sets up logging with given config c. It exits on error.
func MustSetup(c LogConf) {
	Must(SetUp(c))
}

// SetUp sets up the logx. If already set up, just return nil.
// we allow SetUp to be called multiple times, because for example
// we need to allow different service frameworks to initialize logx respectively.
// the same logic for SetUp
func SetUp(c LogConf) error {
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
	case consoleMode:
		setupWithConsole(c)
		return nil
	case volumeMode:
		return setupWithVolume(c)
	default:
		return setupWithFiles(c)
	}
}

// Alert alerts v in alert level, and the message is written to error log.
func Alert(v string) {
	output(errorLog, levelAlert, v)
}

// Close closes the logging.
func Close() error {
	if writeConsole {
		return nil
	}

	if atomic.LoadUint32(&initialized) == 0 {
		return ErrLogNotInitialized
	}

	atomic.StoreUint32(&initialized, 0)

	if infoLog != nil {
		if err := infoLog.Close(); err != nil {
			return err
		}
	}

	if errorLog != nil {
		if err := errorLog.Close(); err != nil {
			return err
		}
	}

	if severeLog != nil {
		if err := severeLog.Close(); err != nil {
			return err
		}
	}

	if slowLog != nil {
		if err := slowLog.Close(); err != nil {
			return err
		}
	}

	if statLog != nil {
		if err := statLog.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Disable disables the logging.
func Disable() {
	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)

		infoLog = iox.NopCloser(ioutil.Discard)
		errorLog = iox.NopCloser(ioutil.Discard)
		severeLog = iox.NopCloser(ioutil.Discard)
		slowLog = iox.NopCloser(ioutil.Discard)
		statLog = iox.NopCloser(ioutil.Discard)
		stackLog = ioutil.Discard
	})
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
	case time.Duration:
		return LogField{Key: key, Value: fmt.Sprint(val)}
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
	output(severeLog, levelFatal, msg)
	os.Exit(1)
}

// SetLevel sets the logging level. It can be used to suppress some logs.
func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
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

func buildFields(fields ...LogField) []string {
	var items []string

	for _, field := range fields {
		items = append(items, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	return items
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
		output(errorLog, levelError, v)
	}
}

func errorFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		output(errorLog, levelError, content, fields...)
	}
}

func errorTextSync(msg string) {
	if shallLog(ErrorLevel) {
		output(errorLog, levelError, msg)
	}
}

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}

func infoAnySync(val interface{}) {
	if shallLog(InfoLevel) {
		output(infoLog, levelInfo, val)
	}
}

func infoFieldsSync(content string, fields ...LogField) {
	if shallLog(InfoLevel) {
		output(infoLog, levelInfo, content, fields...)
	}
}

func infoTextSync(msg string) {
	if shallLog(InfoLevel) {
		output(infoLog, levelInfo, msg)
	}
}

func output(writer io.Writer, level string, val interface{}, fields ...LogField) {
	fields = append(fields, Field(callerKey, getCaller(callerDepth)))

	switch atomic.LoadUint32(&encoding) {
	case plainEncodingType:
		writePlainAny(writer, level, val, buildFields(fields...)...)
	default:
		entry := make(logEntryWithFields)
		for _, field := range fields {
			entry[field.Key] = field.Value
		}
		entry[timestampKey] = getTimestamp()
		entry[levelKey] = level
		entry[contentKey] = val
		writeJson(writer, entry)
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

func setupWithConsole(c LogConf) {
	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)
		writeConsole = true
		setupLogLevel(c)

		infoLog = newLogWriter(log.New(os.Stdout, "", flags))
		errorLog = newLogWriter(log.New(os.Stderr, "", flags))
		severeLog = newLogWriter(log.New(os.Stderr, "", flags))
		slowLog = newLogWriter(log.New(os.Stderr, "", flags))
		stackLog = newLessWriter(errorLog, options.logStackCooldownMills)
		statLog = infoLog
	})
}

func setupWithFiles(c LogConf) error {
	var opts []LogOption
	var err error

	if len(c.Path) == 0 {
		return ErrLogPathNotSet
	}

	opts = append(opts, WithCooldownMillis(c.StackCooldownMillis))
	if c.Compress {
		opts = append(opts, WithGzip())
	}
	if c.KeepDays > 0 {
		opts = append(opts, WithKeepDays(c.KeepDays))
	}

	accessFile := path.Join(c.Path, accessFilename)
	errorFile := path.Join(c.Path, errorFilename)
	severeFile := path.Join(c.Path, severeFilename)
	slowFile := path.Join(c.Path, slowFilename)
	statFile := path.Join(c.Path, statFilename)

	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)
		handleOptions(opts)
		setupLogLevel(c)

		if infoLog, err = createOutput(accessFile); err != nil {
			return
		}

		if errorLog, err = createOutput(errorFile); err != nil {
			return
		}

		if severeLog, err = createOutput(severeFile); err != nil {
			return
		}

		if slowLog, err = createOutput(slowFile); err != nil {
			return
		}

		if statLog, err = createOutput(statFile); err != nil {
			return
		}

		stackLog = newLessWriter(errorLog, options.logStackCooldownMills)
	})

	return err
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
		output(severeLog, levelSevere, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
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
		output(slowLog, levelSlow, v)
	}
}

func slowFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		output(slowLog, levelSlow, content, fields...)
	}
}

func slowTextSync(msg string) {
	if shallLog(ErrorLevel) {
		output(slowLog, levelSlow, msg)
	}
}

func stackSync(msg string) {
	if shallLog(ErrorLevel) {
		output(stackLog, levelError, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func statSync(msg string) {
	if shallLogStat() && shallLog(InfoLevel) {
		output(statLog, levelStat, msg)
	}
}

func writeJson(writer io.Writer, info interface{}) {
	if content, err := json.Marshal(info); err != nil {
		log.Println(err.Error())
	} else if atomic.LoadUint32(&initialized) == 0 || writer == nil {
		log.Println(string(content))
	} else {
		writer.Write(append(content, '\n'))
	}
}

func writePlainAny(writer io.Writer, level string, val interface{}, fields ...string) {
	switch v := val.(type) {
	case string:
		writePlainText(writer, level, v, fields...)
	case error:
		writePlainText(writer, level, v.Error(), fields...)
	case fmt.Stringer:
		writePlainText(writer, level, v.String(), fields...)
	default:
		var buf strings.Builder
		buf.WriteString(getTimestamp())
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(level)
		buf.WriteByte(plainEncodingSep)
		if err := json.NewEncoder(&buf).Encode(val); err != nil {
			log.Println(err.Error())
			return
		}

		for _, item := range fields {
			buf.WriteByte(plainEncodingSep)
			buf.WriteString(item)
		}
		buf.WriteByte('\n')
		if atomic.LoadUint32(&initialized) == 0 || writer == nil {
			log.Println(buf.String())
			return
		}

		if _, err := fmt.Fprint(writer, buf.String()); err != nil {
			log.Println(err.Error())
		}
	}
}

func writePlainText(writer io.Writer, level, msg string, fields ...string) {
	var buf strings.Builder
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(msg)
	for _, item := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(item)
	}
	buf.WriteByte('\n')
	if atomic.LoadUint32(&initialized) == 0 || writer == nil {
		log.Println(buf.String())
		return
	}

	if _, err := fmt.Fprint(writer, buf.String()); err != nil {
		log.Println(err.Error())
	}
}

type logWriter struct {
	logger *log.Logger
}

func newLogWriter(logger *log.Logger) logWriter {
	return logWriter{
		logger: logger,
	}
}

func (lw logWriter) Close() error {
	return nil
}

func (lw logWriter) Write(data []byte) (int, error) {
	lw.logger.Print(string(data))
	return len(data), nil
}
