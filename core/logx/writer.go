package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path"
	"runtime/debug"
	"sync"
	"sync/atomic"

	fatihcolor "github.com/fatih/color"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/errorx"
)

type (
	// Writer is the interface for writing logs.
	// It's designed to let users customize their own log writer,
	// such as writing logs to a kafka, a database, or using third-party loggers.
	Writer interface {
		// Alert sends an alert message, if your writer implemented alerting functionality.
		Alert(v any)
		// Close closes the writer.
		Close() error
		// Debug logs a message at debug level.
		Debug(v any, fields ...LogField)
		// Error logs a message at error level.
		Error(v any, fields ...LogField)
		// Info logs a message at info level.
		Info(v any, fields ...LogField)
		// Severe logs a message at severe level.
		Severe(v any)
		// Slow logs a message at slow level.
		Slow(v any, fields ...LogField)
		// Stack logs a message at error level.
		Stack(v any)
		// Stat logs a message at stat level.
		Stat(v any, fields ...LogField)
	}

	atomicWriter struct {
		writer Writer
		lock   sync.RWMutex
	}

	comboWriter struct {
		writers []Writer
	}

	concreteWriter struct {
		infoLog   io.WriteCloser
		errorLog  io.WriteCloser
		severeLog io.WriteCloser
		slowLog   io.WriteCloser
		statLog   io.WriteCloser
		stackLog  io.Writer
	}
)

// NewWriter creates a new Writer with the given io.Writer.
func NewWriter(w io.Writer) Writer {
	lw := newLogWriter(log.New(w, "", flags))

	return &concreteWriter{
		infoLog:   lw,
		errorLog:  lw,
		severeLog: lw,
		slowLog:   lw,
		statLog:   lw,
		stackLog:  lw,
	}
}

func (w *atomicWriter) Load() Writer {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.writer
}

func (w *atomicWriter) Store(v Writer) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.writer = v
}

func (w *atomicWriter) StoreIfNil(v Writer) Writer {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.writer == nil {
		w.writer = v
	}

	return w.writer
}

func (w *atomicWriter) Swap(v Writer) Writer {
	w.lock.Lock()
	defer w.lock.Unlock()
	old := w.writer
	w.writer = v
	return old
}

func (c comboWriter) Alert(v any) {
	for _, w := range c.writers {
		w.Alert(v)
	}
}

func (c comboWriter) Close() error {
	var be errorx.BatchError
	for _, w := range c.writers {
		be.Add(w.Close())
	}
	return be.Err()
}

func (c comboWriter) Debug(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Debug(v, fields...)
	}
}

func (c comboWriter) Error(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Error(v, fields...)
	}
}

func (c comboWriter) Info(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Info(v, fields...)
	}
}

func (c comboWriter) Severe(v any) {
	for _, w := range c.writers {
		w.Severe(v)
	}
}

func (c comboWriter) Slow(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Slow(v, fields...)
	}
}

func (c comboWriter) Stack(v any) {
	for _, w := range c.writers {
		w.Stack(v)
	}
}

func (c comboWriter) Stat(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Stat(v, fields...)
	}
}

func newConsoleWriter() Writer {
	outLog := newLogWriter(log.New(fatihcolor.Output, "", flags))
	errLog := newLogWriter(log.New(fatihcolor.Error, "", flags))
	return &concreteWriter{
		infoLog:   outLog,
		errorLog:  errLog,
		severeLog: errLog,
		slowLog:   errLog,
		stackLog:  newLessWriter(errLog, options.logStackCooldownMills),
		statLog:   outLog,
	}
}

func newFileWriter(c LogConf) (Writer, error) {
	var err error
	var opts []LogOption
	var infoLog io.WriteCloser
	var errorLog io.WriteCloser
	var severeLog io.WriteCloser
	var slowLog io.WriteCloser
	var statLog io.WriteCloser
	var stackLog io.Writer

	if len(c.Path) == 0 {
		return nil, ErrLogPathNotSet
	}

	opts = append(opts, WithCooldownMillis(c.StackCooldownMillis))
	if c.Compress {
		opts = append(opts, WithGzip())
	}
	if c.KeepDays > 0 {
		opts = append(opts, WithKeepDays(c.KeepDays))
	}
	if c.MaxBackups > 0 {
		opts = append(opts, WithMaxBackups(c.MaxBackups))
	}
	if c.MaxSize > 0 {
		opts = append(opts, WithMaxSize(c.MaxSize))
	}

	opts = append(opts, WithRotation(c.Rotation))

	accessFile := path.Join(c.Path, accessFilename)
	errorFile := path.Join(c.Path, errorFilename)
	severeFile := path.Join(c.Path, severeFilename)
	slowFile := path.Join(c.Path, slowFilename)
	statFile := path.Join(c.Path, statFilename)

	handleOptions(opts)
	setupLogLevel(c)

	if infoLog, err = createOutput(accessFile); err != nil {
		return nil, err
	}

	if errorLog, err = createOutput(errorFile); err != nil {
		return nil, err
	}

	if severeLog, err = createOutput(severeFile); err != nil {
		return nil, err
	}

	if slowLog, err = createOutput(slowFile); err != nil {
		return nil, err
	}

	if statLog, err = createOutput(statFile); err != nil {
		return nil, err
	}

	stackLog = newLessWriter(errorLog, options.logStackCooldownMills)

	return &concreteWriter{
		infoLog:   infoLog,
		errorLog:  errorLog,
		severeLog: severeLog,
		slowLog:   slowLog,
		statLog:   statLog,
		stackLog:  stackLog,
	}, nil
}

func (w *concreteWriter) Alert(v any) {
	output(w.errorLog, levelAlert, v)
}

func (w *concreteWriter) Close() error {
	if err := w.infoLog.Close(); err != nil {
		return err
	}

	if err := w.errorLog.Close(); err != nil {
		return err
	}

	if err := w.severeLog.Close(); err != nil {
		return err
	}

	if err := w.slowLog.Close(); err != nil {
		return err
	}

	return w.statLog.Close()
}

func (w *concreteWriter) Debug(v any, fields ...LogField) {
	output(w.infoLog, levelDebug, v, fields...)
}

func (w *concreteWriter) Error(v any, fields ...LogField) {
	output(w.errorLog, levelError, v, fields...)
}

func (w *concreteWriter) Info(v any, fields ...LogField) {
	output(w.infoLog, levelInfo, v, fields...)
}

func (w *concreteWriter) Severe(v any) {
	output(w.severeLog, levelFatal, v)
}

func (w *concreteWriter) Slow(v any, fields ...LogField) {
	output(w.slowLog, levelSlow, v, fields...)
}

func (w *concreteWriter) Stack(v any) {
	output(w.stackLog, levelError, v)
}

func (w *concreteWriter) Stat(v any, fields ...LogField) {
	output(w.statLog, levelStat, v, fields...)
}

type nopWriter struct{}

func (n nopWriter) Alert(_ any) {
}

func (n nopWriter) Close() error {
	return nil
}

func (n nopWriter) Debug(_ any, _ ...LogField) {
}

func (n nopWriter) Error(_ any, _ ...LogField) {
}

func (n nopWriter) Info(_ any, _ ...LogField) {
}

func (n nopWriter) Severe(_ any) {
}

func (n nopWriter) Slow(_ any, _ ...LogField) {
}

func (n nopWriter) Stack(_ any) {
}

func (n nopWriter) Stat(_ any, _ ...LogField) {
}

func buildPlainFields(fields logEntry) []string {
	items := make([]string, 0, len(fields))
	for k, v := range fields {
		items = append(items, fmt.Sprintf("%s=%+v", k, v))
	}

	return items
}

func marshalJson(t interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	// go 1.5+ will append a newline to the end of the json string
	// https://github.com/golang/go/issues/13520
	if l := buf.Len(); l > 0 && buf.Bytes()[l-1] == '\n' {
		buf.Truncate(l - 1)
	}

	return buf.Bytes(), err
}

func mergeGlobalFields(fields []LogField) []LogField {
	globals := globalFields.Load()
	if globals == nil {
		return fields
	}

	gf := globals.([]LogField)
	ret := make([]LogField, 0, len(gf)+len(fields))
	ret = append(ret, gf...)
	ret = append(ret, fields...)

	return ret
}

func output(writer io.Writer, level string, val any, fields ...LogField) {
	// only truncate string content, don't know how to truncate the values of other types.
	if v, ok := val.(string); ok {
		maxLen := atomic.LoadUint32(&maxContentLength)
		if maxLen > 0 && len(v) > int(maxLen) {
			val = v[:maxLen]
			fields = append(fields, truncatedField)
		}
	}

	// +3 for timestamp, level and content
	entry := make(logEntry, len(fields)+3)
	for _, field := range fields {
		entry[field.Key] = field.Value
	}

	switch atomic.LoadUint32(&encoding) {
	case plainEncodingType:
		plainFields := buildPlainFields(entry)
		writePlainAny(writer, level, val, plainFields...)
	default:
		entry[timestampKey] = getTimestamp()
		entry[levelKey] = level
		entry[contentKey] = val
		writeJson(writer, entry)
	}
}

func wrapLevelWithColor(level string) string {
	var colour color.Color
	switch level {
	case levelAlert:
		colour = color.FgRed
	case levelError:
		colour = color.FgRed
	case levelFatal:
		colour = color.FgRed
	case levelInfo:
		colour = color.FgBlue
	case levelSlow:
		colour = color.FgYellow
	case levelDebug:
		colour = color.FgYellow
	case levelStat:
		colour = color.FgGreen
	}

	if colour == color.NoColor {
		return level
	}

	return color.WithColorPadding(level, colour)
}

func writeJson(writer io.Writer, info any) {
	if content, err := marshalJson(info); err != nil {
		log.Printf("err: %s\n\n%s", err.Error(), debug.Stack())
	} else if writer == nil {
		log.Println(string(content))
	} else {
		if _, err := writer.Write(append(content, '\n')); err != nil {
			log.Println(err.Error())
		}
	}
}

func writePlainAny(writer io.Writer, level string, val any, fields ...string) {
	level = wrapLevelWithColor(level)

	switch v := val.(type) {
	case string:
		writePlainText(writer, level, v, fields...)
	case error:
		writePlainText(writer, level, v.Error(), fields...)
	case fmt.Stringer:
		writePlainText(writer, level, v.String(), fields...)
	default:
		writePlainValue(writer, level, v, fields...)
	}
}

func writePlainText(writer io.Writer, level, msg string, fields ...string) {
	var buf bytes.Buffer
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
	if writer == nil {
		log.Println(buf.String())
		return
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		log.Println(err.Error())
	}
}

func writePlainValue(writer io.Writer, level string, val any, fields ...string) {
	var buf bytes.Buffer
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		log.Printf("err: %s\n\n%s", err.Error(), debug.Stack())
		return
	}

	for _, item := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(item)
	}
	buf.WriteByte('\n')
	if writer == nil {
		log.Println(buf.String())
		return
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		log.Println(err.Error())
	}
}
