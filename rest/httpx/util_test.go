package httpx

import (
	"golang.org/x/text/language"
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

	result2 := v.Validate(u, "zh,en;q=0.9,en-US;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6,la;q=0.5,ja;q=0.4,id;q=0.3,fr;q=0.2")
	if result2 != "Password长度必须至少为6个字符 " {
		t.Error(result2)
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

func TestGetLangFromHeader(t *testing.T) {
	type GetLangTest struct {
		get      string
		expected language.Tag
	}
	tests := []GetLangTest{
		{"en", language.English},
		{"en_US", enUSParseDash},
		{"en-US", enUSParseDash},
		{"zh", language.Chinese},
		{"zh-CN", zhCNParseDash},
		{"zh_CN", zhCNParseDash},
		{"fr", language.English},
		{"zh,en;q=0.9,en-US;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6,la;q=0.5,ja;q=0.4,id;q=0.3,fr;q=0.2", language.Chinese},
	}
	for _, test := range tests {
		tag, _ := language.MatchStrings(matcher, test.get)
		if tag != test.expected {
			t.Errorf("input: %v, got=%v, expected=%v\n", test.get, tag, test.expected)
		}
	}
}
