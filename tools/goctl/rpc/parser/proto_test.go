package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProto_GetImportMessage(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test_import.proto", []string{"./"})
	assert.Nil(t, err)
	pbMsg, existed := data.GetImportMessage("common.CommonMessage")
	assert.Equal(t, true, existed)
	assert.Equal(t, "CommonMessage", pbMsg.Name)
	assert.Equal(t, "common", pbMsg.Package)
	assert.Equal(t, "common", pbMsg.PbPackage)
	assert.Equal(t, "github/go-zero/common", pbMsg.GoPackage)
}

func TestProto_HasGrpcService(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test_import.proto", []string{"./"})
	assert.Nil(t, err)
	assert.Equal(t, true, data.HasGrpcService())
}

func TestProto_HasGrpcService2(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./common/common.proto", nil)
	assert.Nil(t, err)
	assert.Equal(t, false, data.HasGrpcService())
}
