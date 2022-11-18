package format

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

type formatData struct {
	input    string
	expected string
}

func TestFormat_ImportLiteralStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    `import ""`,
			expected: ``,
		},
		{
			input:    `import"aa"`,
			expected: `import "aa"`,
		},
		{
			input: `/*aa*/import "aa"`,
			expected: `/*aa*/
import "aa"`,
		},
		{
			input: `/*aa*/import /*bb*/"aa"`,
			expected: `/*aa*/
import "aa"`,
		},
		{
			input: `/*aa*/import /*bb*/"aa"// cc`,
			expected: `/*aa*/
import "aa" // cc`,
		},
	})
}

func TestFormat_ImportGroupStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    `import()`,
			expected: ``,
		},
		{
			input:    `import("aa")`,
			expected: `import (
	"aa"
)`,
		},
		{
			input:    `import(
"aa")`,
			expected: `import (
	"aa"
)`,
		},
		{
			input:    `import(
"aa"
)`,
			expected: `import (
	"aa"
)`,
		},
		{
			input:    `import("aa""bb")`,
			expected: `import (
	"aa"
	"bb"
)`,
		},
		{
			input:    `/*aa*/import("aa""bb")`,
			expected: `/*aa*/
import (
	"aa"
	"bb"
)`,
		},
		{
			input:    `/*aa*/import("aa""bb")// bb`,
			expected: `/*aa*/
import (
	"aa"
	"bb"
) // bb`,
		},
		{
			input:    `/*aa*/import(// bb
"aa""bb")// cc`,
			expected: `/*aa*/
import ( // bb
	"aa"
	"bb"
) // cc`,
		},
		{
			input:    `import(// aa
"aa" // bb
"bb" // cc
)// dd`,
			expected: `import ( // aa
	"aa" // bb
	"bb" // cc
) // dd`,
		},
		{
			input:    `import (// aa
/*bb*/
	"aa" // cc
/*dd*/
	"bb" // ee
) // ff`,
			expected: `import ( // aa
	/*bb*/
	"aa" // cc
	/*dd*/
	"bb" // ee
) // ff`,
		},
	})
}

func testRun(t *testing.T, testData []formatData) {
	for _, v := range testData {
		buffer := bytes.NewBuffer(nil)
		err := Format([]byte(v.input), buffer)
		assert.NoError(t, err)
		assert.Equal(t, v.expected, buffer.String())
	}
}
