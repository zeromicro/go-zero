package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTomlToJson(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"\n",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a = \"foo\"\nb = 1\nc = \"${FOO}\"\nd = \"abcd!@#$112\"\n",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			got, err := TomlToJson([]byte(test.input))
			assert.NoError(t, err)
			assert.Equal(t, test.expect, string(got))
		})
	}
}

func TestTomlToJsonError(t *testing.T) {
	_, err := TomlToJson([]byte("foo"))
	assert.Error(t, err)
}

func TestYamlToJson(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112\n",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			input:  "a: foo\nb: 1\nc: ${FOO}\nd: abcd!@#$112\n",
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			got, err := YamlToJson([]byte(test.input))
			assert.NoError(t, err)
			assert.Equal(t, test.expect, string(got))
		})
	}
}

func TestYamlToJsonError(t *testing.T) {
	_, err := YamlToJson([]byte("':foo"))
	assert.Error(t, err)
}

func TestYamlToJsonSlice(t *testing.T) {
	b, err := YamlToJson([]byte(`foo:
- bar
- baz`))
	assert.NoError(t, err)
	assert.Equal(t, `{"foo":["bar","baz"]}
`, string(b))
}

func TestJson5ToJson(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "standard json",
			input:  `{"a":"foo","b":1,"c":"${FOO}","d":"abcd!@#$112"}`,
			expect: "{\"a\":\"foo\",\"b\":1,\"c\":\"${FOO}\",\"d\":\"abcd!@#$112\"}\n",
		},
		{
			name:   "json5 with comments",
			input:  `{/*comment*/"a":"foo","b":1}`,
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
		{
			name:   "json5 with trailing commas",
			input:  `{"a":"foo","b":1,}`,
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
		{
			name:   "json5 with unquoted keys",
			input:  `{a:"foo",b:1}`,
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
		{
			name:   "json5 with single quotes",
			input:  `{"a":'foo',"b":1}`,
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
		{
			name:   "json5 with line comments",
			input:  "{\n// This is a comment\n\"a\":\"foo\",\n\"b\":1\n}",
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
		{
			name:   "json5 all features combined",
			input:  "{\n// comment\na: 'foo', // trailing comma\nb: 1,\n}",
			expect: "{\"a\":\"foo\",\"b\":1}\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := Json5ToJson([]byte(test.input))
			assert.NoError(t, err)
			assert.Equal(t, test.expect, string(got))
		})
	}
}

func TestJson5ToJsonError(t *testing.T) {
	// Invalid JSON5: unquoted string value
	_, err := Json5ToJson([]byte("{a: foo}"))
	assert.Error(t, err)
}

func TestJson5ToJsonInfinity(t *testing.T) {
	// JSON5 allows Infinity but standard JSON does not
	_, err := Json5ToJson([]byte(`{value: Infinity}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Infinity")

	// Negative infinity
	_, err = Json5ToJson([]byte(`{value: -Infinity}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Infinity")

	// Infinity in array
	_, err = Json5ToJson([]byte(`{values: [1, Infinity, 3]}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Infinity")
}

func TestJson5ToJsonNaN(t *testing.T) {
	// JSON5 allows NaN but standard JSON does not
	_, err := Json5ToJson([]byte(`{value: NaN}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NaN")

	// NaN in nested structure
	_, err = Json5ToJson([]byte(`{nested: {value: NaN}}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NaN")
}

func TestJson5ToJsonSlice(t *testing.T) {
	b, err := Json5ToJson([]byte(`{
		// comment
		foo: [
			'bar',
			"baz",  // trailing comma
		],
	}`))
	assert.NoError(t, err)
	assert.Equal(t, `{"foo":["bar","baz"]}
`, string(b))
}
