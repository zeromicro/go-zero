//go:build linux || darwin || freebsd

package proc

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestDumpGoroutines(t *testing.T) {
	t.Run("real file", func(t *testing.T) {
		buf := logtest.NewCollector(t)
		dumpGoroutines(fileCreator{})
		assert.True(t, strings.Contains(buf.String(), ".dump"))
	})

	t.Run("fake file", func(t *testing.T) {
		const msg = "any message"
		buf := logtest.NewCollector(t)
		err := errors.New(msg)
		dumpGoroutines(fakeCreator{
			file: &os.File{},
			err:  err,
		})
		assert.True(t, strings.Contains(buf.String(), msg))
	})
}

type fakeCreator struct {
	file *os.File
	err  error
}

func (fc fakeCreator) Create(_ string) (file *os.File, err error) {
	return fc.file, fc.err
}
