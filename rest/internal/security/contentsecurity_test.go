package security

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/codec"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCyeDYV2ieOtNDi6tuNtAbmUjN9
pTHluAU5yiKEz8826QohcxqUKP3hybZBcm60p+rUxMAJFBJ8Dt+UJ6sEMzrf1rOF
YOImVvORkXjpFU7sCJkhnLMs/kxtRzcZJG6ADUlG4GDCNcZpY/qELEvwgm2kCcHi
tGC2mO8opFFFHTR0aQIDAQAB
-----END PUBLIC KEY-----`
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCyeDYV2ieOtNDi6tuNtAbmUjN9pTHluAU5yiKEz8826QohcxqU
KP3hybZBcm60p+rUxMAJFBJ8Dt+UJ6sEMzrf1rOFYOImVvORkXjpFU7sCJkhnLMs
/kxtRzcZJG6ADUlG4GDCNcZpY/qELEvwgm2kCcHitGC2mO8opFFFHTR0aQIDAQAB
AoGAcENv+jT9VyZkk6karLuG75DbtPiaN5+XIfAF4Ld76FWVOs9V88cJVON20xpx
ixBphqexCMToj8MnXuHJEN5M9H15XXx/9IuiMm3FOw0i6o0+4V8XwHr47siT6T+r
HuZEyXER/2qrm0nxyC17TXtd/+TtpfQWSbivl6xcAEo9RRECQQDj6OR6AbMQAIDn
v+AhP/y7duDZimWJIuMwhigA1T2qDbtOoAEcjv3DB1dAswJ7clcnkxI9a6/0RDF9
0IEHUcX9AkEAyHdcegWiayEnbatxWcNWm1/5jFnCN+GTRRFrOhBCyFr2ZdjFV4T+
acGtG6omXWaZJy1GZz6pybOGy93NwLB93QJARKMJ0/iZDbOpHqI5hKn5mhd2Je25
IHDCTQXKHF4cAQ+7njUvwIMLx2V5kIGYuMa5mrB/KMI6rmyvHv3hLewhnQJBAMMb
cPUOENMllINnzk2oEd3tXiscnSvYL4aUeoErnGP2LERZ40/YD+mMZ9g6FVboaX04
0oHf+k5mnXZD7WJyJD0CQQDJ2HyFbNaUUHK+lcifCibfzKTgmnNh9ZpePFumgJzI
EfFE5H+nzsbbry2XgJbWzRNvuFTOLWn4zM+aFyy9WvbO
-----END RSA PRIVATE KEY-----`
	body = "hello world!"
)

var key = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")

func TestContentSecurity(t *testing.T) {
	tests := []struct {
		name        string
		mode        string
		extraKey    string
		extraSecret string
		extraTime   string
		err         error
		code        int
	}{
		{
			name: "encrypted",
			mode: "1",
		},
		{
			name: "unencrypted",
			mode: "0",
		},
		{
			name: "bad content type",
			mode: "a",
			err:  ErrInvalidContentType,
		},
		{
			name:        "bad secret",
			mode:        "1",
			extraSecret: "any",
			err:         ErrInvalidSecret,
		},
		{
			name:     "bad key",
			mode:     "1",
			extraKey: "any",
			err:      ErrInvalidKey,
		},
		{
			name:      "bad time",
			mode:      "1",
			extraTime: "any",
			code:      httpx.CodeSignatureInvalidHeader,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			r, err := http.NewRequest(http.MethodPost, "http://localhost:3333/a/b?c=first&d=second",
				strings.NewReader(body))
			assert.Nil(t, err)

			timestamp := time.Now().Unix()
			sha := sha256.New()
			sha.Write([]byte(body))
			bodySign := fmt.Sprintf("%x", sha.Sum(nil))
			contentOfSign := strings.Join([]string{
				strconv.FormatInt(timestamp, 10),
				http.MethodPost,
				r.URL.Path,
				r.URL.RawQuery,
				bodySign,
			}, "\n")
			sign := hs256(key, contentOfSign)
			content := strings.Join([]string{
				"version=v1",
				"type=" + test.mode,
				fmt.Sprintf("key=%s", base64.StdEncoding.EncodeToString(key)) + test.extraKey,
				"time=" + strconv.FormatInt(timestamp, 10) + test.extraTime,
			}, "; ")

			encrypter, err := codec.NewRsaEncrypter([]byte(pubKey))
			if err != nil {
				log.Fatal(err)
			}

			output, err := encrypter.Encrypt([]byte(content))
			if err != nil {
				log.Fatal(err)
			}

			encryptedContent := base64.StdEncoding.EncodeToString(output)
			r.Header.Set("X-Content-Security", strings.Join([]string{
				fmt.Sprintf("key=%s", fingerprint(pubKey)),
				"secret=" + encryptedContent + test.extraSecret,
				"signature=" + sign,
			}, "; "))

			file, err := fs.TempFilenameWithText(priKey)
			assert.Nil(t, err)
			defer os.Remove(file)

			dec, err := codec.NewRsaDecrypter(file)
			assert.Nil(t, err)

			header, err := ParseContentSecurity(map[string]codec.RsaDecrypter{
				fingerprint(pubKey): dec,
			}, r)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			encrypted := test.mode != "0"
			assert.Equal(t, encrypted, header.Encrypted())
			assert.Equal(t, test.code, VerifySignature(r, header, time.Minute))
		})
	}
}

func fingerprint(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func hs256(key []byte, body string) string {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
