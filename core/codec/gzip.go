package codec

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Gzip(bs []byte) []byte {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)
	w.Write(bs)
	w.Close()

	return b.Bytes()
}

func Gunzip(bs []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var c bytes.Buffer
	_, err = io.Copy(&c, r)
	if err != nil {
		return nil, err
	}

	return c.Bytes(), nil
}
