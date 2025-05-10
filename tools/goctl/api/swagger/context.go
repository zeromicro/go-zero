package swagger

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"testing"
)

type Context struct {
	UseDefinitions         bool
	WrapCodeMsg            bool
	BizCodeEnumDescription string
}

func testingContext(_ *testing.T) Context {
	return Context{}
}

func contextFromApi(info spec.Info) Context {
	if len(info.Properties) == 0 {
		return Context{}
	}
	return Context{
		UseDefinitions:         getBoolFromKVOrDefault(info.Properties, propertyKeyUseDefinitions, defaultValueOfPropertyUseDefinition),
		WrapCodeMsg:            getBoolFromKVOrDefault(info.Properties, propertyKeyWrapCodeMsg, false),
		BizCodeEnumDescription: getStringFromKVOrDefault(info.Properties, propertyKeyBizCodeEnumDescription, "business code"),
	}
}
