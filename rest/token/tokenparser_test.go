package token

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/timex"
)

func TestTokenParser(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	keys := []struct {
		key     string
		prevKey string
	}{
		{
			key,
			prevKey,
		},
		{
			key,
			"",
		},
	}

	for _, pair := range keys {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
		token, err := buildToken(key, map[string]any{
			"key": "value",
		}, 3600)
		assert.Nil(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		parser := NewTokenParser(WithResetDuration(time.Minute))
		tok, err := parser.ParseToken(req, pair.key, pair.prevKey)
		assert.Nil(t, err)
		assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	}
}

func TestTokenParser_CustomHeader(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	token, err := buildToken(key, map[string]any{"key": "value"}, 3600)
	assert.Nil(t, err)
	req.Header.Set("Token", token)

	parser := NewTokenParser(WithExtractor([]string{"header:Token"}))
	tok, err := parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	parser.resetTime = timex.Now() - time.Hour
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
}

func TestTokenParser_URLArgument(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	token, err := buildToken(key, map[string]any{"key": "value"}, 3600)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodGet, "http://localhost?token="+token, http.NoBody)

	parser := NewTokenParser(WithExtractor([]string{"query:token"}))
	tok, err := parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	parser.resetTime = timex.Now() - time.Hour
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
}

func TestTokenParser_FormArgument(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	token, err := buildToken(key, map[string]any{"key": "value"}, 3600)
	assert.Nil(t, err)

	// create form data
	form := url.Values{}
	form.Add("form_token", token)

	// Using httptest.NewRequest to create a fake POST request
	req := httptest.NewRequest(http.MethodPost, "http://localhost", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	parser := NewTokenParser(WithExtractor([]string{"form:form_token"}))
	tok, err := parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	parser.resetTime = timex.Now() - time.Hour
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])

}

func TestTokenParser_Expired(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	token, err := buildToken(key, map[string]any{
		"key": "value",
	}, 3600)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	parser := NewTokenParser(WithResetDuration(time.Second))
	tok, err := parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	parser.resetTime = timex.Now() - time.Hour
	tok, err = parser.ParseToken(req, key, prevKey)
	assert.Nil(t, err)
	assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
}

func buildToken(secretKey string, payloads map[string]any, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
