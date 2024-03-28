package httpx

import (
	"net/http"
	"net/url"
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

type User struct {
	UserId   []int  `form:"userId"`
	Username string `form:"username"`
	Password string `form:"password,optional"`
}

func TestQueryURLArrayParameterFormValues(t *testing.T) {
	u := User{
		UserId:   []int{1, 2},
		Username: "admin",
	}

	// Unit 1
	r, err := http.NewRequest(http.MethodGet, "/test/?userId=1&userId=2&username=admin&password=", nil)
	assert.Nil(t, err)

	params, err := GetFormValues(r)
	assert.Nil(t, err)

	var u1 User
	assert.Nil(t, formUnmarshaler.Unmarshal(params, &u1))
	assert.Equal(t, u, u1)

	// Unit 2
	r, err = http.NewRequest(http.MethodGet, "/test/?userId=[1,2]&username=admin", nil)
	assert.Nil(t, err)

	params, err = GetFormValues(r)
	assert.Nil(t, err)

	var u2 User
	assert.Nil(t, formUnmarshaler.Unmarshal(params, &u2))
	assert.Equal(t, u, u2)
}

func TestPostFormValues(t *testing.T) {
	u := User{
		UserId:   []int{1, 2},
		Username: "admin",
	}

	formData := url.Values{
		"userId":   {"1", "2"},
		"username": {"admin"},
	}

	r, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(formData.Encode()))
	assert.Nil(t, err)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	params, err := GetFormValues(r)
	assert.Nil(t, err)

	var u1 User
	assert.Nil(t, formUnmarshaler.Unmarshal(params, &u1))
	assert.Equal(t, u, u1)

}
