package internal

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fullstorydev/grpcurl"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/hash"
)

const b64pb = `CpgBCgtoZWxsby5wcm90bxIFaGVsbG8iHQoHUmVxdWVzdBISCgRwaW5nGAEgASgJUgRwaW5nIh4KCFJlc3BvbnNlEhIKBHBvbmcYASABKAlSBHBvbmcyMAoFSGVsbG8SJwoEUGluZxIOLmhlbGxvLlJlcXVlc3QaDy5oZWxsby5SZXNwb25zZUIJWgcuL2hlbGxvYgZwcm90bzM=`

func TestGetMethods(t *testing.T) {
	tmpfile, err := ioutil.TempFile(os.TempDir(), hash.Md5Hex([]byte(b64pb)))
	assert.Nil(t, err)
	b, err := base64.StdEncoding.DecodeString(b64pb)
	assert.Nil(t, err)
	assert.Nil(t, ioutil.WriteFile(tmpfile.Name(), b, os.ModeTemporary))
	defer os.Remove(tmpfile.Name())

	source, err := grpcurl.DescriptorSourceFromProtoSets(tmpfile.Name())
	assert.Nil(t, err)
	methods, err := GetMethods(source)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"hello.Hello/Ping"}, methods)
}
