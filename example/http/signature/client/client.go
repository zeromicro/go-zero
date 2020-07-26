package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"zero/core/codec"
)

const pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD7bq4FLG0ctccbEFEsUBuRxkjE
eJ5U+0CAEjJk20V9/u2Fu76i1oKoShCs7GXtAFbDb5A/ImIXkPY62nAaxTGK4KVH
miYbRgh5Fy6336KepLCtCmV/r0PKZeCyJH9uYLs7EuE1z9Hgm5UUjmpHDhJtkAwR
my47YlhspwszKdRP+wIDAQAB
-----END PUBLIC KEY-----`

var (
	crypt = flag.Bool("crypt", false, "encrypt body or not")
	key   = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")
)

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

func main() {
	flag.Parse()

	var err error
	body := "hello world!"
	if *crypt {
		bodyBytes, err := codec.EcbEncrypt(key, []byte(body))
		if err != nil {
			log.Fatal(err)
		}
		body = base64.StdEncoding.EncodeToString(bodyBytes)
	}

	r, err := http.NewRequest(http.MethodPost, "http://localhost:3333/a/b?c=first&d=second", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

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
	var mode string
	if *crypt {
		mode = "1"
	} else {
		mode = "0"
	}
	content := strings.Join([]string{
		"version=v1",
		"type=" + mode,
		fmt.Sprintf("key=%s", base64.StdEncoding.EncodeToString(key)),
		"time=" + strconv.FormatInt(timestamp, 10),
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
		"secret=" + encryptedContent,
		"signature=" + sign,
	}, "; "))
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	io.Copy(os.Stdout, resp.Body)
}
