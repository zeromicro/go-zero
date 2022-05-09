package token

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/timex"
)

const privateSecret = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDnfKw2I9iHe5JK
T7o3FLP7vY3jDTB+5Ey/T7k8Hp6IqEOBUcpc+ECuMKoTYzNLQvMMgcL4+qauZHEO
pKSx0n0CwFP8kkYpDH4Y0JvIueV7IE+0ZqHqysDH1lU1MhqrPPk5676kVpFSvEek
7fWaZG4NhpwYbiYuZoWjaOcJDlE2vmnicWVv5iYwhtUtpuXMhgMIVSoxzgUOlLSl
LC7eBAYXLL+IloZpdRjWv0gVxU5Odvlh4tHqaaqGLiNW6zqUprvVCSBVvjiBBY65
GuDvtQ2sv6bGReDts/hZAUqB9xTXQWK4YcIka9dXnJ9WZlAY4P5np3WiBBY6Z14h
pdekOMo3AgMBAAECggEAa+6gaRnrkrgWLJnh7F610LHAH1Z9/xw5gJYOey6XooY5
+2kEXrbNiapdEm8VcokDxBgYrXhJEVT5teckd1j6OrcsMb6OAgO2I6HYkQ3EJtWY
9DdKVaw1mLehwQzcjG0Ak3YMzJkkZxwsl4TwGA2tlpbl3yo0mTvqIZf+6SUIzum0
TUlDwnNRqPUU0worN9cgmVY8iwvFB9TFfjUcq0e8oKd2RTcK1y5PWhE16s2iqQqp
kbdFetT24+fNcI1LZBvrUVhYsGCSK9sNvaA9jKUJ3NHcX7iXc5eVBc2MB5x9Uoon
6U/UHIdFsvDu4u7Em3XmJLHsZA030IFlLaaAg1EQaQKBgQD9A/MOfcwxoz5F73F2
vo7HFUVcRVDQD8dowoqQfld+hNcZVYlMbBxs+JqNiXMox2aUNgJioYHNMFns6KQJ
1gleRWOW/XLzIiCX688POJyPJAXMM3MRido1oTZbPT3FQuE3s6vxOpAMzpWmd9vq
VoH0bMkW6QKzJb2voO6lpCVVVQKBgQDqN7ZR0iXZzsQxPkR2HfB2FDsTn+WRhd7J
WeaaNOkLCFW99AGW3ShcTyNp3I8GKfNcbiErEEjwuPloN4xZeVg2Lu2LFNegt+5+
NNZemvu4gEzdlawWJ9lAt3PsovHJO7Fad1zF/ogezBgwn+T5PWla+K507IyAVHMT
8M3/jyWhWwKBgFZW7q5XR0L5DdsXpoR66oYNQCoIjVcyyz14hYhhVMIb2rsOcVfe
3KRjAXqjGOUlhl+1PoMh0gWPJmCt0qx4maHN0/pGat+FGdI96d6r1uERzditBetK
O2hppv7jmxyhgfFcIqSi810rce3ooOcKtjYOmWB0CzPPATfZlxZ3OTYxAoGAeX36
sciLX8b0WALPqmFvWSC3YD+h6nGBlfpvNvBZLiLdrxHCPUps5C0c1o3VFsJt/TUX
OWpSG6Qno1qlD8h07G49Q9bE3xZpvMeVpy9HgXXz6UD5KejztbEzjb0cJGE1ZxLh
acbVPvxpU9etA2hKnSi//eCyJOMpal+Py4+qWl8CgYBSCrw+sae1x0tHYaO6Nrvs
sN682QT+oBvfU80+IN49CRcB6MS/8RyOqI29FsShPydauyCV4jl+EpXucf20m2Qb
lkeE4QCcWaCltaDvYFeWOCXtv7jQQwTSRafLJFt3lJLNn8zf7Mjyckp12DREkROd
0YDdOe3wXidBuP9PY4FJWQ==
-----END PRIVATE KEY-----
`

const publicSecret = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA53ysNiPYh3uSSk+6NxSz
+72N4w0wfuRMv0+5PB6eiKhDgVHKXPhArjCqE2MzS0LzDIHC+PqmrmRxDqSksdJ9
AsBT/JJGKQx+GNCbyLnleyBPtGah6srAx9ZVNTIaqzz5Oeu+pFaRUrxHpO31mmRu
DYacGG4mLmaFo2jnCQ5RNr5p4nFlb+YmMIbVLablzIYDCFUqMc4FDpS0pSwu3gQG
Fyy/iJaGaXUY1r9IFcVOTnb5YeLR6mmqhi4jVus6lKa71QkgVb44gQWOuRrg77UN
rL+mxkXg7bP4WQFKgfcU10FiuGHCJGvXV5yfVmZQGOD+Z6d1ogQWOmdeIaXXpDjK
NwIDAQAB
-----END PUBLIC KEY-----
`

const publicPreSecret = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1UVZIUsajbnEOfad37X8
GdpfkNZNxwZ4aqySb1AwyOUiJv4/tMzUxXb+ZxGtk6hj+35MDAvRJbhFWGc3cyui
1kwTDIHzXzx0I4uadknlJ0dJx+rupkfJZvvP9ECpLdgAPR1+HAzJsx4BM5cJXMNp
9fK2dGdYqEdPmUbxENYNUQdzlWCTo1Md3ETPBONS0Q2pBDip0upSczzR2qLJ9O2V
Fyq4Y59FDAXWP8BGSvv/IRT3AOiWVa0SUIcax14iFRZEJ9uvl8/e1isErsHXpA5X
NTjld9jgf+8mCynWhr2xuG8PHsAUP0eIain/zEj8xQVfYK1/7D1K4qNa3ugzltJ0
OwIDAQAB
-----END PUBLIC KEY-----
`

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
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		token, err := buildToken(key, map[string]interface{}{
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

func TestTokenParser_Expired(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	token, err := buildToken(key, map[string]interface{}{
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

func TestTokenParser_WithPemFile(t *testing.T) {
	const (
		key     = publicSecret
		prevKey = publicPreSecret
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
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		encodeKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateSecret))
		assert.Nil(t, err)
		token, err := buildPemToken(encodeKey, map[string]interface{}{
			"key": "value",
		}, 3600)
		assert.Nil(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		parser := NewTokenParser(WithResetDuration(time.Minute))
		decodeKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pair.key))
		assert.Nil(t, err)
		var decodePreKey interface{}
		decodePreKey = ""
		if len(pair.prevKey) > 0 {
			decodePreKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pair.prevKey))
			assert.Nil(t, err)
		}
		tok, err := parser.ParseToken(req, decodeKey, decodePreKey)
		assert.Nil(t, err)
		assert.Equal(t, "value", tok.Claims.(jwt.MapClaims)["key"])
	}
}

// generate public key
func buildPemToken(pem interface{}, payloads map[string]interface{}, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = claims
	// 继续添加秘钥值，生成最后一部分
	return token.SignedString(pem)
}

func buildToken(secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
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
