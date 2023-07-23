package mapping

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Bar struct {
	Val string `json:"val"`
}

func TestFieldOptionOptionalDep(t *testing.T) {
	var bar Bar
	rt := reflect.TypeOf(bar)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		val, opt, err := parseKeyAndOptions(jsonTagKey, field)
		assert.Equal(t, "val", val)
		assert.Nil(t, opt)
		assert.Nil(t, err)
	}

	// check nil working
	var o *fieldOptions
	check := func(o *fieldOptions) {
		assert.Equal(t, 0, len(o.optionalDep()))
	}
	check(o)
}
