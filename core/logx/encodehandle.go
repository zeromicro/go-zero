package logx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

const defaultEncoding = jsonEncodingName

var ErrFormatFailed = errors.New("log format error")

var (
	encodingStore map[string]LogEncoding

	_ LogEncoding = (*JsonLogEncoding)(nil)
	_ LogEncoding = (*PlainTextLogEncoding)(nil)
)

type (
	atomicEncoding struct {
		encoding LogEncoding
		lock     sync.RWMutex
	}

	LogEncoding interface {
		Output(l *LogData) ([]byte, error)
	}

	LogData struct {
		Level   string
		Content any
		Fields  []LogField
	}

	JsonLogEncoding struct {
		// Use the `context` field to collect the log context
		UseContextField bool
	}

	PlainTextLogEncoding struct{}
)

func init() {
	encodingStore = make(map[string]LogEncoding)

	EncodingRegister(jsonEncodingName, &JsonLogEncoding{})
	EncodingRegister(plainEncodingName, &JsonLogEncoding{})
}

func EncodingRegister(name string, format LogEncoding) {
	if f, ok := encodingStore[name]; ok {
		panic(fmt.Sprintf("logFormat number [%s] already exist [%T]", name, f))
	}

	encodingStore[name] = format
}

func getEncodingHandle(name string) LogEncoding {
	if f, ok := encodingStore[name]; ok {
		return f
	}

	return encodingStore[defaultEncoding]
}

// atomicEncoding start

func (e *atomicEncoding) Store(encoder LogEncoding) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.encoding = encoder
}

func (e *atomicEncoding) Swap(encoder LogEncoding) LogEncoding {
	e.lock.Lock()
	defer e.lock.Unlock()

	old := e.encoding
	e.encoding = encoder

	return old
}

func (e *atomicEncoding) Load() LogEncoding {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.encoding
}

// JsonLogEncoding start

func (j *JsonLogEncoding) Output(l *LogData) ([]byte, error) {
	entry := make(logEntry)
	if j.UseContextField {
		_context := make(map[string]any, len(l.Fields))

		for _, field := range l.Fields {
			_context[field.Key] = field.Value
		}

		entry["context"] = _context
	} else {
		for _, field := range l.Fields {
			entry[field.Key] = field.Value
		}
	}

	entry[timestampKey] = getTimestamp()
	entry[levelKey] = l.Level
	entry[contentKey] = l.Content

	marshal, err := json.Marshal(entry)
	if err != nil {
		return nil, ErrFormatFailed
	}

	return marshal, nil
}

// PlainTextLogEncoding start

func (p *PlainTextLogEncoding) Output(l *LogData) ([]byte, error) {
	level := wrapLevelWithColor(l.Level)
	fields := buildPlainFields(l.Fields...)

	switch v := l.Content.(type) {
	case string:
		return p.plainText(level, v, fields...)
	case error:
		return p.plainText(level, v.Error(), fields...)
	case fmt.Stringer:
		return p.plainText(level, v.String(), fields...)
	default:
		return p.plainValue(level, v, fields...)
	}
}

func (p *PlainTextLogEncoding) plainText(level string, msg string, fields ...string) ([]byte, error) {
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

	return buf.Bytes(), nil
}

func (p *PlainTextLogEncoding) plainValue(level string, val any, fields ...string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, ErrFormatFailed
	}

	for _, item := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(item)
	}

	return buf.Bytes(), nil
}
