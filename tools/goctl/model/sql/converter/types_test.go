package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/ddl-parser/parser"
)

func TestConvertDataType(t *testing.T) {
	v, _, err := ConvertDataType(parser.TinyInt, false, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, _, err = ConvertDataType(parser.TinyInt, false, true, true)
	assert.Nil(t, err)
	assert.Equal(t, "uint64", v)

	v, _, err = ConvertDataType(parser.TinyInt, true, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullInt64", v)

	v, _, err = ConvertDataType(parser.Timestamp, false, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	v, _, err = ConvertDataType(parser.Timestamp, true, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullTime", v)

	v, _, err = ConvertDataType(parser.Decimal, false, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "float64", v)
}

func TestConvertStringDataType(t *testing.T) {
	type (
		input struct {
			dataType      string
			isDefaultNull bool
			unsigned      bool
			strict        bool
		}
		result struct {
			goType    string
			thirdPkg  string
			isPQArray bool
		}
	)
	var testData = []struct {
		input input
		want  result
	}{
		{
			input: input{
				dataType:      "bigint",
				isDefaultNull: false,
				unsigned:      false,
				strict:        false,
			},
			want: result{
				goType: "int64",
			},
		},
		{
			input: input{
				dataType:      "bigint",
				isDefaultNull: true,
				unsigned:      false,
				strict:        false,
			},
			want: result{
				goType: "sql.NullInt64",
			},
		},
		{
			input: input{
				dataType:      "bigint",
				isDefaultNull: false,
				unsigned:      true,
				strict:        false,
			},
			want: result{
				goType: "uint64",
			},
		},
		{
			input: input{
				dataType:      "_int2",
				isDefaultNull: false,
				unsigned:      false,
				strict:        false,
			},
			want: result{
				goType:    "pq.Int64Array",
				isPQArray: true,
			},
		},
	}
	for _, data := range testData {
		tp, thirdPkg, isPQArray, err := ConvertStringDataType(data.input.dataType, data.input.isDefaultNull, data.input.unsigned, data.input.strict)
		assert.NoError(t, err)
		assert.Equal(t, data.want, result{
			goType:    tp,
			thirdPkg:  thirdPkg,
			isPQArray: isPQArray,
		})
	}
}
