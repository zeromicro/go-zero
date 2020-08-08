package main

import (
	"fmt"
	"log"

	"github.com/tal-tech/go-zero/core/codec"
)

const (
	pubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD7bq4FLG0ctccbEFEsUBuRxkjE
eJ5U+0CAEjJk20V9/u2Fu76i1oKoShCs7GXtAFbDb5A/ImIXkPY62nAaxTGK4KVH
miYbRgh5Fy6336KepLCtCmV/r0PKZeCyJH9uYLs7EuE1z9Hgm5UUjmpHDhJtkAwR
my47YlhspwszKdRP+wIDAQAB
-----END PUBLIC KEY-----`
	body = "hello"
)

var key = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")

func main() {
	encrypter, err := codec.NewRsaEncrypter([]byte(pubKey))
	if err != nil {
		log.Fatal(err)
	}

	decrypter, err := codec.NewRsaDecrypter("private.pem")
	if err != nil {
		log.Fatal(err)
	}

	output, err := encrypter.Encrypt([]byte(body))
	if err != nil {
		log.Fatal(err)
	}

	actual, err := decrypter.Decrypt(output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(actual)

	out, err := codec.EcbEncrypt(key, []byte(body))
	if err != nil {
		log.Fatal(err)
	}

	ret, err := codec.EcbDecrypt(key, out)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(ret))
}
