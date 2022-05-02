package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempFileWithText(t *testing.T) {
	f, err := TempFileWithText("test")
	if err != nil {
		t.Error(err)
	}
	if f == nil {
		t.Error("TempFileWithText returned nil")
	}
	if f.Name() == "" {
		t.Error("TempFileWithText returned empty file name")
	}
	defer os.Remove(f.Name())

	bs, err := ioutil.ReadAll(f)
	assert.Nil(t, err)
	if len(bs) != 4 {
		t.Error("TempFileWithText returned wrong file size")
	}
	if f.Close() != nil {
		t.Error("TempFileWithText returned error on close")
	}
}

func TestTempFilenameWithText(t *testing.T) {
	f, err := TempFilenameWithText("test")
	if err != nil {
		t.Error(err)
	}
	if f == "" {
		t.Error("TempFilenameWithText returned empty file name")
	}
	defer os.Remove(f)

	bs, err := ioutil.ReadFile(f)
	assert.Nil(t, err)
	if len(bs) != 4 {
		t.Error("TempFilenameWithText returned wrong file size")
	}
}
