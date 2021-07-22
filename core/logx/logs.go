package logx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/iox"
	"github.com/tal-tech/go-zero/core/sysx"
	"github.com/tal-tech/go-zero/core/timex"
)

const (
	// InfoLevel logs everything
	InfoLevel = iota
	// ErrorLevel includes errors, slows, stacks
	ErrorLevel
	// SevereLevel only log severe messages
	SevereLevel
)

const (
	accessFilename = "access.log"
	errorFilename  = "error.log"
	severeFilename = "severe.log"
	slowFilename   = "slow.log"
	statFilename   = "stat.log"

	consoleMode = "console"
	volumeMode  = "volume"

	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"

	backupFileDelimiter = "-"
	callerInnerDepth    = 5
	flags               = 0x0
)

var (
	// ErrLogPathNotSet is an error that indicates the log path is not set.
	ErrLogPathNotSet = errors.New("log path must be set")
	// ErrLogNotInitialized is an error that log is not initialized.
	ErrLogNotInitialized = errors.New("log not initialized")
	// ErrLogServiceNameNotSet is an error that indicates that the service name is not set.
	ErrLogServiceNameNotSet = errors.New("log service name must be set")

	timeFormat   = "2006-01-02T15:04:05.000Z07"
	writeConsole bool
	logLevel     uint32
	infoLog      io.WriteCloser
	errorLog     io.WriteCloser
	severeLog    io.WriteCloser
	slowLog      io.WriteCloser
	statLog      io.WriteCloser
	stackLog     io.Writer

	once        sync.Once
	initialized uint32
	options     logOptions
)

type (
	logEntry struct {
		Timestamp string `json:"@timestamp"`
		Level     string `json:"level"`
		Duration  string `json:"duration,omitempty"`
		Content   string `json:"content"`
	}

	logOptions struct {
		gzipEnabled           bool
		logStackCooldownMills int
		keepDays              int
	}

	// LogOption defines the method to customize the logging.
	LogOption func(options *logOptions)

	// A Logger represents a logger.
	Logger interface {
		Error(...interface{})
		Errorf(string, ...interface{})
		Info(...interface{})
		Infof(string, ...interface{})
		Slow(...interface{})
		Slowf(string, ...interface{})
		WithDuration(time.Duration) Logger
	}
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

func init() {


	Error = func(v ...interface{}) {
		errorSync(fmt.Sprint(v...), callerInnerDepth)
	}

	Errorf = func(format string, v ...interface{}) {
		errorSync(fmt.Sprintf(format, v...),  callerInnerDepth)
	}

	ErrorCaller = func(callDepth int, v ...interface{}) {
		errorSync(fmt.Sprint(v...), callDepth + callerInnerDepth)
	}

	ErrorCallerf = func(callDepth int, format string, v ...interface{}) {
		errorSync(fmt.Sprintf(format, v...), callDepth+callerInnerDepth)
	}

	ErrorStack = func(v ...interface{}) {
		stackSync(fmt.Sprint(v...))
	}

	ErrorStackf = func(format string, v ...interface{}) {
		stackSync(fmt.Sprintf(format, v...))
	}

	Info = func(v ...interface{}) {
		infoSync(fmt.Sprint(v...))
	}

	Infof = func(format string, v ...interface{}) {
		infoSync(fmt.Sprintf(format, v...))
	}

	Must = func(err error) {
		if err != nil {
			msg := formatWithCaller(err.Error(), 3)
			log.Print(msg)
			output(severeLog, levelFatal, msg)
			os.Exit(1)
		}
	}

	SetLevel = func (level uint32) {
		atomic.StoreUint32(&logLevel, level)
	}

	Severe = func (v ...interface{}) {
		severeSync(fmt.Sprint(v...))
	}

	Severef = func(format string, v ...interface{}) {
		severeSync(fmt.Sprintf(format, v...))
	}

	Slow = func(v ...interface{}) {
		slowSync(fmt.Sprint(v...))
	}

	Slowf = func(format string, v ...interface{}) {
		slowSync(fmt.Sprintf(format, v...))
	}

	Stat = func(v ...interface{}) {
		statSync(fmt.Sprint(v...))
	}

	Statf = func(format string, v ...interface{}) {
		statSync(fmt.Sprintf(format, v...))
	}

}

// Error writes v into error log.
var Error func(v ...interface{})

// Errorf writes v with format into error log.
var Errorf func (format string, v ...interface{})

// ErrorCaller writes v with context into error log.
var ErrorCaller func (callDepth int, v ...interface{})

// ErrorCallerf writes v with context in format into error log.
var ErrorCallerf func (callDepth int, format string, v ...interface{})

// ErrorStack writes v along with call stack into error log.
var  ErrorStack func (v ...interface{})

// ErrorStackf writes v along with call stack in format into error log.
var ErrorStackf func (format string, v ...interface{})

// Info writes v into access log.
var  Info func (v ...interface{})

// Infof writes v with format into access log.
var Infof func (format string, v ...interface{})

// Must checks if err is nil, otherwise logs the err and exits.
var Must func (err error)

// SetLevel sets the logging level. It can be used to suppress some logs.
var SetLevel func (level uint32)

// Severe writes v into severe log.
var  Severe func (v ...interface{})

// Severef writes v with format into severe log.
var Severef func (format string, v ...interface{})

// Slow writes v into slow log.
var Slow func  (v ...interface{})

// Slowf writes v with format into slow log.
var Slowf func (format string, v ...interface{})

// Stat writes v into stat log.
var Stat func (v ...interface{})

// Statf writes v with format into stat log.
var  Statf func (format string, v ...interface{})

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

func errorSync(msg string, callDepth int) {
	if shouldLog(ErrorLevel) {
		outputError(errorLog, msg, callDepth)
	}
}

func formatWithCaller(msg string, callDepth int) string {
	var buf strings.Builder

	caller := getCaller(callDepth)
	if len(caller) > 0 {
		buf.WriteString(caller)
		buf.WriteByte(' ')
	}

	buf.WriteString(msg)

	return buf.String()
}

func getCaller(callDepth int) string {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		buf.WriteString(short)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	return buf.String()
}

func getTimestamp() string {
	return timex.Time().Format(timeFormat)
}

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}

func infoSync(msg string) {
	if shouldLog(InfoLevel) {
		output(infoLog, levelInfo, msg)
	}
}

func output(writer io.Writer, level, msg string) {
	info := logEntry{
		Timestamp: getTimestamp(),
		Level:     level,
		Content:   msg,
	}
	outputJson(writer, info)
}

func outputError(writer io.Writer, msg string, callDepth int) {
	content := formatWithCaller(msg, callDepth)
	output(writer, levelError, content)
}

func outputJson(writer io.Writer, info interface{}) {
	if content, err := json.Marshal(info); err != nil {
		log.Println(err.Error())
	} else if atomic.LoadUint32(&initialized) == 0 || writer == nil {
		log.Println(string(content))
	} else {
		writer.Write(append(content, '\n'))
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
	if shouldLog(SevereLevel) {
		output(severeLog, levelSevere, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func shouldLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

func slowSync(msg string) {
	if shouldLog(ErrorLevel) {
		output(slowLog, levelSlow, msg)
	}
}

func stackSync(msg string) {
	if shouldLog(ErrorLevel) {
		output(stackLog, levelError, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func statSync(msg string) {
	if shouldLog(InfoLevel) {
		output(statLog, levelStat, msg)
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
