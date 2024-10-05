package gogen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"path/filepath"
	"testing"
)

func TestGenRoutesString(t *testing.T) {
	parse, err := parser.Parse(filepath.Join("testdata", "example.api"), nil)
	assert.Nil(t, err)

	routesString, err := GenRoutesString("example", &config.Config{
		NamingFormat: "gozero",
	}, parse)
	assert.NotNil(t, routesString)
	assert.Nil(t, err)
	fmt.Println(routesString)
}
