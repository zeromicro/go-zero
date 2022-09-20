package httpx

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	host := "8.8.8.8"
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	r.Header.Set(xForwardedFor, host)
	assert.Equal(t, host, GetRemoteAddr(r))
}

func TestGetRemoteAddrNoHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	assert.True(t, len(GetRemoteAddr(r)) == 0)
}

func TestValidator(t *testing.T) {
	v := NewValidator()
	type User struct {
		Username string `validate:"required,alphanum,max=20"`
		Password string `validate:"required,min=6,max=30"`
	}
	u := User{
		Username: "admin",
		Password: "1",
	}
	result := v.Validate(u, "en")
	if result != "Password must be at least 6 characters in length " {
		t.Error(result)
	}

	u = User{
		Username: "admin",
		Password: "123456",
	}
	result = v.Validate(u, "en")
	if result != "" {
		t.Error(result)
	}
}
