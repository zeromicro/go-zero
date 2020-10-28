package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBaseNameWithoutExt(t *testing.T) {
	var filePath = "cmd/api/bookstore.api"
	fileName := FileBaseNameWithoutExt(filePath)
	assert.Equal(t, "bookstore", fileName)
}
