package codec

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < 10000; i++ {
		fmt.Fprint(&buf, i)
	}

	bs := Gzip(buf.Bytes())
	actual, err := Gunzip(bs)

	assert.Nil(t, err)
	assert.True(t, len(bs) < buf.Len())
	assert.Equal(t, buf.Bytes(), actual)
}
