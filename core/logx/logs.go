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
	timeFormat = "2006-01-02T15:04:05.000Z07"

	accessFilename = "access.log"
	errorFilename  = "error.log"
	severeFilename = "severe.log"
	slowFilename   = "slow.log"
	statFilename   = "stat.log"

	consoleMode = "console"
	volumeMode  = "volume"

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
	ErrLogPathNotSet        = errors.New("log path must be set")
	ErrLogNotInitialized    = errors.New("log not initialized")
	ErrLogServiceNameNotSet = errors.New("log service name must be set")

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

	LogOption func(options *logOptions)

	Logger interface {
		Error(...interface{})
		Errorf(string, ...interface{})
		Info(...interface{})
		Infof(string, ...interface{})
		Slow(...interface{})
		Slowf(string, ...interface{})
	}
)

func MustSetup(c LogConf) {
	Must(SetUp(c))
}

// SetUp sets up the logx. If already set up, just return nil.
// we allow SetUp to be called multiple times, because for example
// we need to allow different service frameworks to initialize logx respectively.
// the same logic for SetUp
func SetUp(c LogConf) error {
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

func Error(v ...interface{}) {
	ErrorCaller(1, v...)
}

func Errorf(format string, v ...interface{}) {
	ErrorCallerf(1, format, v...)
}

func ErrorCaller(callDepth int, v ...interface{}) {
	errorSync(fmt.Sprint(v...), callDepth+callerInnerDepth)
}

func ErrorCallerf(callDepth int, format string, v ...interface{}) {
	errorSync(fmt.Sprintf(format, v...), callDepth+callerInnerDepth)
}

func ErrorStack(v ...interface{}) {
	// there is newline in stack string
	stackSync(fmt.Sprint(v...))
}

func ErrorStackf(format string, v ...interface{}) {
	// there is newline in stack string
	stackSync(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	infoSync(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	infoSync(fmt.Sprintf(format, v...))
}

func Must(err error) {
	if err != nil {
		msg := formatWithCaller(err.Error(), 3)
		log.Print(msg)
		output(severeLog, levelFatal, msg)
		os.Exit(1)
	}
}

func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
}

func Severe(v ...interface{}) {
	severeSync(fmt.Sprint(v...))
}

func Severef(format string, v ...interface{}) {
	severeSync(fmt.Sprintf(format, v...))
}

func Slow(v ...interface{}) {
	slowSync(fmt.Sprint(v...))
}

func Slowf(format string, v ...interface{}) {
	slowSync(fmt.Sprintf(format, v...))
}

func Stat(v ...interface{}) {
	statSync(fmt.Sprint(v...))
}

func Statf(format string, v ...interface{}) {
	statSync(fmt.Sprintf(format, v...))
}

func WithCooldownMillis(millis int) LogOption {
	return func(opts *logOptions) {
		opts.logStackCooldownMills = millis
	}
}

func WithKeepDays(days int) LogOption {
	return func(opts *logOptions) {
		opts.keepDays = days
	}
}

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
		stackLog = NewLessWriter(errorLog, options.logStackCooldownMills)
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

		stackLog = NewLessWriter(errorLog, options.logStackCooldownMills)
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
