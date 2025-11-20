package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*******************    💫 Codegeex Suggestion    *******************/
// TestProto_GetImportMessage tests the functionality of retrieving imported messages from a proto file.
// It verifies that:
// 1. The proto parser can successfully parse a proto file with imports
// 2. The GetImportMessage method can correctly locate and return imported message information
// 3. The returned message contains accurate metadata including name, package, and go package path
// The test uses a test proto file located at "./test_import.proto" with import paths set to "./"
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

/****************  4a7a2f3b23d54ad48e80b9ccd4e34a8c  ****************/

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
