package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
)

func TestGenSqlWithCache(t *testing.T) {
	generator := gen.NewDefaultGenerator("./test.sql", "./modelwithcache")
	err := generator.Start(true)
	assert.Nil(t, err)
}

func TestGenSqlWithoutCache(t *testing.T) {
	generator := gen.NewDefaultGenerator("./test.sql", "./modelwithoutcache")
	err := generator.Start(false)
	assert.Nil(t, err)
}
