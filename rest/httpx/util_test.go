package httpx

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//@enhance

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

func TestParseAcceptLanguage(t *testing.T) {
	data := []struct {
		Str    string
		Target []string
	}{
		{
			"zh",
			[]string{"zh"},
		},
		{
			"zh,en;q=0.9,en-US;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6,la;q=0.5,ja;q=0.4,id;q=0.3,fr;q=0.2",
			[]string{"zh", "en", "en-US", "zh-CN", "zh-TW", "la", "ja", "id", "fr"},
		},
		{
			"zh-cn,zh;q=0.9",
			[]string{"zh-CN", "zh"},
		},
		{
			"en,zh;q=0.9",
			[]string{"en", "zh"},
		},
	}

	for _, v := range data {
		tmp, err := ParseAcceptLanguage(v.Str)
		if err != nil {
			t.Error(err)
		}
		for i := range tmp {
			if v.Target[i] != tmp[i] {
				t.Error("parse error: ", v.Str)
			}
		}
	}
}
