package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var files = []string{
	"example",
	"empty",
	"syntax",
	"info",
	"types",
	"service",
}

func TestGrammar(t *testing.T) {
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			_, err := parser.Parse("./apis/" + file + ".api")
			assert.Nil(t, err)
		})
	}
}
