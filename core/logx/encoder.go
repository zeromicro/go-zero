package logx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

const defaultEncoding = jsonEncodingName

var ErrLogEncodeFailed = errors.New("log encode error")

var (
	_ LogEncoder = (*JsonLogEncoder)(nil)
	_ LogEncoder = (*PlainTextLogEncoder)(nil)
)

type (
	// atomic switch encoder, but not use
	atomicEncoding struct {
		encoding LogEncoder
		lock     sync.RWMutex
	}

	LogEncoder interface {
		Name() string
		Output(l *LogData) ([]byte, error)
	}

	LogData struct {
		Level   string
		Content any
		Fields  []LogField
	}

	JsonLogEncoder struct {
		// Use the `context` field to collect the log context
		UseContextField bool
	}

	PlainTextLogEncoder struct {
		WithColor bool
	}
)

func EncoderRegister(encoder LogEncoder) {
	encoders[encoder.Name()] = encoder
}

func getEncodingHandle(name string) LogEncoder {
	if f, ok := encoders[name]; ok {
		return f
	}

	return encoders[defaultEncoding]
}

// atomicEncoding start

func (e *atomicEncoding) Store(encoder LogEncoder) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.encoding = encoder
}

func (e *atomicEncoding) Swap(encoder LogEncoder) LogEncoder {
	e.lock.Lock()
	defer e.lock.Unlock()

	old := e.encoding
	e.encoding = encoder

	return old
}

func (e *atomicEncoding) Load() LogEncoder {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.encoding
}

// JsonLogEncoder start

func (j *JsonLogEncoder) Name() string {
	return jsonEncodingName
}

func (j *JsonLogEncoder) Output(l *LogData) ([]byte, error) {
	entry := make(logEntry)

	fieldLen := len(l.Fields)
	if j.UseContextField && fieldLen > 0 {
		_context := make(map[string]any, fieldLen)

		for _, field := range l.Fields {
			_context[field.Key] = field.Value
		}

		entry[contextField] = _context
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
		return nil, ErrLogEncodeFailed
	}

	return marshal, nil
}

// PlainTextLogEncoder start

func (p *PlainTextLogEncoder) Name() string {
	return plainEncodingName
}

func (p *PlainTextLogEncoder) Color() bool {
	return p.WithColor
}

func (p *PlainTextLogEncoder) Output(l *LogData) ([]byte, error) {
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

func (p *PlainTextLogEncoder) plainText(level string, msg string, fields ...string) ([]byte, error) {
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

func (p *PlainTextLogEncoder) plainValue(level string, val any, fields ...string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		return nil, ErrLogEncodeFailed
	}

	for _, item := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(item)
	}

	return buf.Bytes(), nil
}
