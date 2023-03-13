package logx

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEncoding(e LogEncoder) {
	if e != nil {
		//encoding.Store(e)
		encoding = e
	}
}

func TestJsonLogEncoding_Output(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	oldE := encoding
	setEncoding(&JsonLogEncoder{UseContextField: true})
	defer func() {
		setEncoding(oldE)
	}()

	doTestStructedLog(t, "info", w, func(v ...any) {
		Info(v)
	})
}

func TestJsonLogEncoding_Output_W(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	oldE := encoding
	setEncoding(&JsonLogEncoder{UseContextField: true})
	defer func() {
		setEncoding(oldE)
	}()

	doTestContextLog(t, "info", w, func(content string, v ...LogField) {
		Infow(content, v...)
	})
}

func doTestContextLog(t *testing.T, level string, w *mockWriter, write func(string, ...LogField)) {
	const message = "hello there"
	write(message, Field("1", "1"), Field("2", "2"))

	var entry map[string]any
	if err := json.Unmarshal([]byte(w.String()), &entry); err != nil {
		t.Error(err)
	}

	fmt.Println(w.String())

	assert.Equal(t, level, entry[levelKey])
	val, ok := entry[contentKey]
	assert.True(t, ok)
	assert.True(t, strings.Contains(val.(string), message))

	ctx, ok := entry[contextField]
	ctxVal := ctx.(map[string]any)
	f1, ok := ctxVal["1"]
	f2, ok := ctxVal["2"]

	assert.True(t, strings.Contains(f1.(string), "1"))
	assert.True(t, strings.Contains(f2.(string), "2"))

}
