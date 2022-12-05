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
