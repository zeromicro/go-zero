package parser

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultProtoParse(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test.proto", nil)
	assert.Nil(t, err)
	assert.Equal(t, "base.proto", func() string {
		ip := data.Import[0]
		return ip.Filename
	}())
	assert.Equal(t, "test", data.Package.Name)
	assert.Equal(t, true, data.GoPackage == "go")
	assert.Equal(t, true, data.PbPackage == "_go")
	assert.Equal(t, []string{"Inline", "Inner", "TestMessage", "TestReply", "TestReq"},
		func() []string {
			var list []string
			for _, item := range data.Message {
				list = append(list, item.Name)
			}
			sort.Strings(list)
			return list
		}())

	assert.Equal(t, true, func() bool {
		if len(data.Service) != 1 {
			return false
		}

		s := data.Service[0]
		if s.Name != "TestService" {
			return false
		}
		rpcOne := s.RPC[0]

		return rpcOne.Name == "TestRpcOne" && rpcOne.RequestType == "TestReq" && rpcOne.ReturnsType == "TestReply"
	}())
}

func TestDefaultProtoParseCaseInvalidRequestType(t *testing.T) {
	p := NewDefaultProtoParser()
	_, err := p.Parse("./test_invalid_request.proto", nil)
	assert.True(t, true, func() bool {
		return strings.Contains(err.Error(), "request type must defined in")
	}())
}

func TestDefaultProtoParseCaseInvalidResponseType(t *testing.T) {
	p := NewDefaultProtoParser()
	_, err := p.Parse("./test_invalid_response.proto", nil)
	assert.True(t, true, func() bool {
		return strings.Contains(err.Error(), "response type must defined in")
	}())
}

func TestDefaultProtoParseError(t *testing.T) {
	p := NewDefaultProtoParser()
	_, err := p.Parse("./nil.proto", nil)
	assert.NotNil(t, err)
}

func TestDefaultProtoParse_Option(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test_option.proto", nil)
	assert.Nil(t, err)
	assert.Equal(t, "github.com/zeromicro/go-zero", data.GoPackage)
	assert.Equal(t, "go_zero", data.PbPackage)
}

func TestDefaultProtoParse_Option2(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test_option2.proto", nil)
	assert.Nil(t, err)
	assert.Equal(t, "stream", data.GoPackage)
	assert.Equal(t, "stream", data.PbPackage)
}

func TestDefaultProtoParse_import(t *testing.T) {
	p := NewDefaultProtoParser()
	data, err := p.Parse("./test_import.proto", []string{"./"})
	assert.Nil(t, err)
	assert.NotNil(t, data.ImportMessageMap)
	assert.Equal(t, 2, len(data.ImportMessageMap))
}
