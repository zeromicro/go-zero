package codec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHmac(t *testing.T) {
	ret := Hmac([]byte("foo"), "bar")
	assert.Equal(t, "f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
		fmt.Sprintf("%x", ret))
}

func TestHmacBase64(t *testing.T) {
	ret := HmacBase64([]byte("foo"), "bar")
	assert.Equal(t, "+TILrwJJFp5zhQzWFW3tAQbiu2rYyrAbe7vr5tEGUxc=", ret)
}
