package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/fs"
)

func TestProperties(t *testing.T) {
	text := `app.name = test

    app.program=app

    # this is comment
    app.threads = 5`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	props, err := LoadProperties(tmpfile)
	assert.Nil(t, err)
	assert.Equal(t, "test", props.GetString("app.name"))
	assert.Equal(t, "app", props.GetString("app.program"))
	assert.Equal(t, 5, props.GetInt("app.threads"))

	val := props.ToString()
	assert.Contains(t, val, "app.name")
	assert.Contains(t, val, "app.program")
	assert.Contains(t, val, "app.threads")
}

func TestLoadProperties_badContent(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	_, err = LoadProperties(filename)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
}

func TestSetString(t *testing.T) {
	key := "a"
	value := "the value of a"
	props := NewProperties()
	props.SetString(key, value)
	assert.Equal(t, value, props.GetString(key))
}

func TestSetInt(t *testing.T) {
	key := "a"
	value := 101
	props := NewProperties()
	props.SetInt(key, value)
	assert.Equal(t, value, props.GetInt(key))
}

func TestLoadBadFile(t *testing.T) {
	_, err := LoadProperties("nosuchfile")
	assert.NotNil(t, err)
}
