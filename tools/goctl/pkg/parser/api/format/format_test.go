package format

import (
	"bytes"
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

type formatData struct {
	input     string
	expected  string
	converter formatResultConvert
}

type formatResultConvert func(s string) string

// EXPERIMENTAL: just for view format code.
func TestFormat(t *testing.T) {
	assert.NoError(t, File("testdata/test_format.api"))
}

//go:embed testdata/test_type_struct_lit.api
var testStructLitData string

//go:embed testdata/expected_type_struct_lit.api
var expectedStructLitData string

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
			input: `import("aa")`,
			expected: `import (
	"aa"
)`,
		},
		{
			input: `import(
"aa")`,
			expected: `import (
	"aa"
)`,
		},
		{
			input: `import(
"aa"
)`,
			expected: `import (
	"aa"
)`,
		},
		{
			input: `import("aa""bb")`,
			expected: `import (
	"aa"
	"bb"
)`,
		},
		{
			input: `/*aa*/import("aa""bb")`,
			expected: `/*aa*/
import (
	"aa"
	"bb"
)`,
		},
		{
			input: `/*aa*/import("aa""bb")// bb`,
			expected: `/*aa*/
import (
	"aa"
	"bb"
) // bb`,
		},
		{
			input: `/*aa*/import(// bb
"aa""bb")// cc`,
			expected: `/*aa*/
import ( // bb
	"aa"
	"bb"
) // cc`,
		},
		{
			input: `import(// aa
"aa" // bb
"bb" // cc
)// dd`,
			expected: `import ( // aa
	"aa" // bb
	"bb" // cc
) // dd`,
		},
		{
			input: `import (// aa
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

func TestFormat_InfoStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    `info()`,
			expected: ``,
		},
		{
			input: `info(foo:"foo")`,
			expected: `info (
	foo: "foo"
)`,
		},
		{
			input: `info(foo:"foo" bar:"bar")`,
			expected: `info (
	foo: "foo"
	bar: "bar"
)`,
		},
		{
			input: `info(foo:"foo" bar:"bar" quux:"quux")`,
			expected: `info (
	foo:  "foo"
	bar:  "bar"
	quux: "quux"
)`,
		},
		{
			input: `info(foo:"foo"
bar: "bar")`,
			expected: `info (
	foo: "foo"
	bar: "bar"
)`,
		},
		{
			input: `info(foo:"foo"// aa
bar: "bar"// bb
)`,
			expected: `info (
	foo: "foo" // aa
	bar: "bar" // bb
)`,
		},
		{
			input: `info(// aa
foo:"foo"// bb
bar: "bar"// cc
)`,
			expected: `info ( // aa
	foo: "foo" // bb
	bar: "bar" // cc
)`,
		},
		{
			input: `/*aa*/info(// bb
foo:"foo"// cc
bar: "bar"// dd
)`,
			expected: `/*aa*/
info ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
)`,
		},
		{
			input: `/*aa*/
info(// bb
foo:"foo"// cc
bar: "bar"// dd
)// ee`,
			expected: `/*aa*/
info ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
) // ee`,
		},
		{
			input: `/*aa*/
info ( // bb
	/*cc*/foo: "foo" // dd
	/*ee*/bar: "bar" // ff
) // gg`,
			expected: `/*aa*/
info ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
		},
		{
			input: `/*aa*/
info/*xx*/( // bb
	/*cc*/foo:/*xx*/ "foo" // dd
	/*ee*/bar:/*xx*/ "bar" // ff
) // gg`,
			expected: `/*aa*/
info ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
		},
	})
}

func TestFormat_SyntaxStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    `syntax="v1"`,
			expected: `syntax = "v1"`,
		},
		{
			input:    `syntax="v1"// aa`,
			expected: `syntax = "v1" // aa`,
		},
		{
			input: `syntax
="v1"// aa`,
			expected: `syntax = "v1" // aa`,
		},
		{
			input: `syntax=
"v1"// aa`,
			expected: `syntax = "v1" // aa`,
		},
		{
			input: `/*aa*/syntax="v1"// bb`,
			expected: `/*aa*/
syntax = "v1" // bb`,
		},
		{
			input: `/*aa*/
syntax="v1"// bb`,
			expected: `/*aa*/
syntax = "v1" // bb`,
		},
		{
			input:    `syntax/*xx*/=/*xx*/"v1"// bb`,
			expected: `syntax = "v1" // bb`,
		},
	})
}

func TestFormat_TypeLiteralStmt(t *testing.T) {
	t.Run("any", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type Any any`,
				expected: `type Any any`,
			},
			{
				input: `type
Any
any
`,
				expected: `type Any any`,
			},
			{
				input:    `type Any=any`,
				expected: `type Any = any`,
			},
			{
				input: `
type
Any
=
any
`,
				expected: `type Any = any`,
			},
			{
				input: `type // aa
Any  // bb
any // cc
`,
				expected: `type // aa
Any // bb
any // cc`,
			},
			{
				input: `
type
Any
=
any`,
				expected: `type Any = any`,
			},
			{
				input: `
type
Any
=
any
`,
				expected: `type Any = any`,
			},
			{
				input:    `type Any any// aa`,
				expected: `type Any any // aa`,
			},
			{
				input:    `type Any=any// aa`,
				expected: `type Any = any // aa`,
			},
			{
				input:    `type Any any/*aa*/// bb`,
				expected: `type Any any /*aa*/ // bb`,
			},
			{
				input:    `type Any = any/*aa*/// bb`,
				expected: `type Any = any /*aa*/ // bb`,
			},
			{
				input:    `type Any/*aa*/ =/*bb*/ any/*cc*/// dd`,
				expected: `type Any /*aa*/ = /*bb*/ any /*cc*/ // dd`,
			},
			{
				input: `/*aa*/type Any any/*bb*/// cc`,
				expected: `/*aa*/
type Any any /*bb*/ // cc`,
			},
			{
				input: `/*aa*/
type
/*bb*/
Any
/*cc*/
any/*dd*/// ee`,
				expected: `/*aa*/
type
/*bb*/
Any
/*cc*/
any /*dd*/ // ee`,
			},
		})
	})
	t.Run("array", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type A [2]int`,
				expected: `type A [2]int`,
			},
			{
				input: `type
A
[2]int
`,
				expected: `type A [2]int`,
			},
			{
				input:    `type A=[2]int`,
				expected: `type A = [2]int`,
			},
			{
				input: `type
A
=
[2]int
`,
				expected: `type A = [2]int`,
			},
			{
				input:    `type A [/*xx*/2/*xx*/]/*xx*/int// aa`,
				expected: `type A [2]int // aa`,
			},
			{
				input: `/*aa*/type/*bb*/A/*cc*/[/*xx*/2/*xx*/]/*xx*/int// dd`,
				expected: `/*aa*/
type /*bb*/ A /*cc*/ [2]int // dd`,
			},
			{
				input: `/*aa*/type
/*bb*/A
/*cc*/[/*xx*/2/*xx*/]/*xx*/int// dd`,
				expected: `/*aa*/
type
/*bb*/
A
/*cc*/
[2]int // dd`,
			},
			{
				input:    `type A [ 2 ] int`,
				expected: `type A [2]int`,
			},
			{
				input: `type A [
2
]
int`,
				expected: `type A [2]int`,
			},
			{
				input: `type A [// aa
2 // bb
] // cc
int`,
				expected: `type A [2]int`,
			},
			{
				input: `type A [// aa
/*xx*/
2 // bb
/*xx*/
] // cc
/*xx*/
int`,
				expected: `type A [2]int`,
			},
			{
				input:    `type A [...]int`,
				expected: `type A [...]int`,
			},
			{
				input:    `type A=[...]int`,
				expected: `type A = [...]int`,
			},
			{
				input:    `type A/*aa*/[/*xx*/.../*xx*/]/*xx*/int// bb`,
				expected: `type A /*aa*/ [...]int // bb`,
			},
			{
				input: `/*aa*/
// bb
type /*cc*/
// dd
A /*ee*/
// ff
[/*xx*/.../*xx*/]/*xx*/int// bb`,
				expected: `/*aa*/
// bb
type /*cc*/
// dd
A /*ee*/
// ff
[...]int // bb`,
			},
			{
				input:    `type A [2][2]int`,
				expected: `type A [2][2]int`,
			},
			{
				input:    `type A=[2][2]int`,
				expected: `type A = [2][2]int`,
			},
			{
				input:    `type A [2][]int`,
				expected: `type A [2][]int`,
			},
			{
				input:    `type A=[2][]int`,
				expected: `type A = [2][]int`,
			},
		})
	})
	t.Run("base", func(t *testing.T) {
		testRun(t, []formatData{
			// base
			{
				input:    `type A int`,
				expected: `type A int`,
			},
			{
				input:    `type A =int`,
				expected: `type A = int`,
			},
			{
				input:    `type/*aa*/A/*bb*/ int// cc`,
				expected: `type /*aa*/ A /*bb*/ int // cc`,
			},
			{
				input:    `type/*aa*/A/*bb*/ =int// cc`,
				expected: `type /*aa*/ A /*bb*/ = int // cc`,
			},
			{
				input:    `type A int// aa`,
				expected: `type A int // aa`,
			},
			{
				input:    `type A=int// aa`,
				expected: `type A = int // aa`,
			},
			{
				input: `/*aa*/type A int`,
				expected: `/*aa*/
type A int`,
			},
			{
				input: `/*aa*/type A = int`,
				expected: `/*aa*/
type A = int`,
			},
			{
				input: `/*aa*/type/*bb*/ A/*cc*/ int// dd`,
				expected: `/*aa*/
type /*bb*/ A /*cc*/ int // dd`,
			},
			{
				input: `/*aa*/type/*bb*/ A/*cc*/ = /*dd*/int// ee`,
				expected: `/*aa*/
type /*bb*/ A /*cc*/ = /*dd*/ int // ee`,
			},
			{
				input: `/*aa*/
type 
/*bb*/
A 
/*cc*/
int`,
				expected: `/*aa*/
type
/*bb*/
A
/*cc*/
int`,
			},
		})
	})
	t.Run("interface", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type any interface{}`,
				expected: `type any interface{}`,
			},
			{
				input:    `type any=interface{}`,
				expected: `type any = interface{}`,
			},
			{
				input: `type
any
interface{}
`,
				expected: `type any interface{}`,
			},
			{
				input: `/*aa*/type /*bb*/any /*cc*/interface{} // dd`,
				expected: `/*aa*/
type /*bb*/ any /*cc*/ interface{} // dd`,
			},
			{
				input: `/*aa*/type 
/*bb*/any 
/*cc*/interface{} // dd`,
				expected: `/*aa*/
type
/*bb*/
any
/*cc*/
interface{} // dd`,
			},
			{
				input: `/*aa*/type 
// bb
any 
// cc
interface{} // dd`,
				expected: `/*aa*/
type
// bb
any
// cc
interface{} // dd`,
			},
		})
	})
	t.Run("map", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type M map[int]int`,
				expected: `type M map[int]int`,
			},
			{
				input:    `type M map [ int ] int`,
				expected: `type M map[int]int`,
			},
			{
				input:    `type M map [/*xx*/int/*xx*/]/*xx*/int // aa`,
				expected: `type M map[int]int // aa`,
			},
			{
				input: `/*aa*/type /*bb*/ M/*cc*/map[int]int // dd`,
				expected: `/*aa*/
type /*bb*/ M /*cc*/ map[int]int // dd`,
			},
			{
				input: `/*aa*/type// bb
// cc
M // dd
// ee
map // ff
[int]// gg
// hh
int // dd`,
				expected: `/*aa*/
type // bb
// cc
M // dd
// ee
map[int]int // dd`,
			},
			{
				input:    `type M map[string][2]int // aa`,
				expected: `type M map[string][2]int // aa`,
			},
			{
				input:    `type M map[string]any`,
				expected: `type M map[string]any`,
			},
			{
				input:    `type M /*aa*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/[/*xx*/2/*xx*/]/*xx*/int// bb`,
				expected: `type M /*aa*/ map[string][2]int // bb`,
			},
			{
				input: `type M /*aa*/
// bb
map/*xx*/
//
[/*xx*/
//
string/*xx*/
//
]/*xx*/
//
[/*xx*/
//
2/*xx*/
//
]/*xx*/
//
int// bb`,
				expected: `type M /*aa*/
// bb
map[string][2]int // bb`,
			},
			{
				input:    `type M map[int]map[string]int`,
				expected: `type M map[int]map[string]int`,
			},
			{
				input:    `type M map/*xx*/[/*xx*/int/*xx*/]/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int// aa`,
				expected: `type M map[int]map[string]int // aa`,
			},
			{
				input:    `type M map/*xx*/[/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int/*xx*/]/*xx*/string // aa`,
				expected: `type M map[map[string]int]string // aa`,
			},
			{
				input:    `type M map[[2]int]int`,
				expected: `type M map[[2]int]int`,
			},
			{
				input:    `type M map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int// aa`,
				expected: `type M map[[2]int]int // aa`,
			},
		})
	})
	t.Run("pointer", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type P *int`,
				expected: `type P *int`,
			},
			{
				input:    `type P=*int`,
				expected: `type P = *int`,
			},
			{
				input: `type 
P 
*int
`,
				expected: `type P *int`,
			},
			{
				input: `/*aa*/type // bb
/*cc*/
P // dd
/*ee*/
*/*ff*/int // gg
`,
				expected: `/*aa*/
type // bb
/*cc*/
P // dd
/*ee*/
*int // gg`,
			},
			{
				input:    `type P *bool`,
				expected: `type P *bool`,
			},
			{
				input:    `type P *[2]int`,
				expected: `type P *[2]int`,
			},
			{
				input:    `type P=*[2]int`,
				expected: `type P = *[2]int`,
			},
			{
				input: `/*aa*/type /*bb*/P /*cc*/*/*xx*/[/*xx*/2/*xx*/]/*xx*/int // dd`,
				expected: `/*aa*/
type /*bb*/ P /*cc*/ *[2]int // dd`,
			},
			{
				input:    `type P *[...]int`,
				expected: `type P *[...]int`,
			},
			{
				input:    `type P=*[...]int`,
				expected: `type P = *[...]int`,
			},
			{
				input: `/*aa*/type /*bb*/P /*cc*/*/*xx*/[/*xx*/.../*xx*/]/*xx*/int // dd`,
				expected: `/*aa*/
type /*bb*/ P /*cc*/ *[...]int // dd`,
			},
			{
				input:    `type P *map[string]int`,
				expected: `type P *map[string]int`,
			},
			{
				input:    `type P=*map[string]int`,
				expected: `type P = *map[string]int`,
			},
			{
				input:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int// bb`,
				expected: `type P /*aa*/ *map[string]int // bb`,
			},
			{
				input:    `type P *interface{}`,
				expected: `type P *interface{}`,
			},
			{
				input:    `type P=*interface{}`,
				expected: `type P = *interface{}`,
			},
			{
				input:    `type P /*aa*/*/*xx*/interface{}// bb`,
				expected: `type P /*aa*/ *interface{} // bb`,
			},
			{
				input:    `type P *any`,
				expected: `type P *any`,
			},
			{
				input:    `type P=*any`,
				expected: `type P = *any`,
			},
			{
				input:    `type P *map[int][2]int`,
				expected: `type P *map[int][2]int`,
			},
			{
				input:    `type P=*map[int][2]int`,
				expected: `type P = *map[int][2]int`,
			},
			{
				input:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/int/*xx*/]/*xx*/[/*xx*/2/*xx*/]/*xx*/int// bb`,
				expected: `type P /*aa*/ *map[int][2]int // bb`,
			},
			{
				input:    `type P *map[[2]int]int`,
				expected: `type P *map[[2]int]int`,
			},
			{
				input:    `type P=*map[[2]int]int`,
				expected: `type P = *map[[2]int]int`,
			},
			{
				input:    `type P /*aa*/*/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int// bb`,
				expected: `type P /*aa*/ *map[[2]int]int // bb`,
			},
		})

	})

	t.Run("slice", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type S []int`,
				expected: `type S []int`,
			},
			{
				input:    `type S=[]int`,
				expected: `type S = []int`,
			},
			{
				input:    `type S	[	]	int	`,
				expected: `type S []int`,
			},
			{
				input:    `type S	[ /*xx*/	]	/*xx*/ int	`,
				expected: `type S []int`,
			},
			{
				input:    `type S [][]int`,
				expected: `type S [][]int`,
			},
			{
				input:    `type S=[][]int`,
				expected: `type S = [][]int`,
			},
			{
				input:    `type S	[	]	[	]	int`,
				expected: `type S [][]int`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/[/*xx*/]/*xx*/int`,
				expected: `type S [][]int`,
			},
			{
				input: `type S [//
]//
[//
]//
int`,
				expected: `type S [][]int`,
			},
			{
				input:    `type S []map[string]int`,
				expected: `type S []map[string]int`,
			},
			{
				input:    `type S=[]map[string]int`,
				expected: `type S = []map[string]int`,
			},
			{
				input: `type S [	]	
map	[	string	]	
int`,
				expected: `type S []map[string]int`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/string/*xx*/]/*xx*/int`,
				expected: `type S []map[string]int`,
			},
			{
				input: `/*aa*/type// bb
// cc
S// dd
// ff
/*gg*/[ // hh
/*xx*/] // ii
/*xx*/map// jj
/*xx*/[/*xx*/string/*xx*/]/*xx*/int// mm`,
				expected: `/*aa*/
type // bb
// cc
S // dd
// ff
/*gg*/
[]map[string]int // mm`,
			},
			{
				input:    `type S []map[[2]int]int`,
				expected: `type S []map[[2]int]int`,
			},
			{
				input:    `type S=[]map[[2]int]int`,
				expected: `type S = []map[[2]int]int`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/int`,
				expected: `type S []map[[2]int]int`,
			},
			{
				input: `/*aa*/type// bb
// cc
/*dd*/S// ee
// ff
/*gg*/[//
/*xx*/]//
/*xx*/map//
/*xx*/[//
/*xx*/[//
/*xx*/2//
/*xx*/]//
/*xx*/int//
/*xx*/]//
/*xx*/int // hh`,
				expected: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
// ff
/*gg*/
[]map[[2]int]int // hh`,
			},
			{
				input:    `type S []map[[2]int]map[int]string`,
				expected: `type S []map[[2]int]map[int]string`,
			},
			{
				input:    `type S=[]map[[2]int]map[int]string`,
				expected: `type S = []map[[2]int]map[int]string`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/map/*xx*/[/*xx*/[/*xx*/2/*xx*/]/*xx*/int/*xx*/]/*xx*/map/*xx*/[/*xx*/int/*xx*/]/*xx*/string`,
				expected: `type S []map[[2]int]map[int]string`,
			},
			{
				input: `/*aa*/type// bb
// cc
/*dd*/S// ee
/*ff*/[//
/*xx*/]//
/*xx*/map
/*xx*/[//
/*xx*/[//
/*xx*/2//
/*xx*/]//
/*xx*/int//
/*xx*/]//
/*xx*/map//
/*xx*/[//
/*xx*/int//
/*xx*/]//
/*xx*/string// gg`,
				expected: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]map[[2]int]map[int]string // gg`,
			},
			{
				input:    `type S []*P`,
				expected: `type S []*P`,
			},
			{
				input:    `type S=[]*P`,
				expected: `type S = []*P`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/*/*xx*/P`,
				expected: `type S []*P`,
			},
			{
				input: `/*aa*/type// bb
// cc
/*dd*/S// ee 
/*ff*/[//
/*xx*/]//
/*xx*/*//
/*xx*/P // gg`,
				expected: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]*P // gg`,
			},
			{
				input:    `type S []*[]int`,
				expected: `type S []*[]int`,
			},
			{
				input:    `type S=[]*[]int`,
				expected: `type S = []*[]int`,
			},
			{
				input:    `type S [/*xx*/]/*xx*/*/*xx*/[/*xx*/]/*xx*/int`,
				expected: `type S []*[]int`,
			},
			{
				input: `/*aa*/
type // bb
// cc
/*dd*/S// ee
/*ff*/[//
/*xx*/]//
/*xx*/*//
/*xx*/[//
/*xx*/]//
/*xx*/int // gg`,
				expected: `/*aa*/
type // bb
// cc
/*dd*/
S // ee
/*ff*/
[]*[]int // gg`,
			},
		})
	})

	t.Run("struct", func(t *testing.T) {
		testRun(t, []formatData{
			{
				input:    `type T {}`,
				expected: `type T {}`,
			},
			{
				input: `type T 	{
			}	`,
				expected: `type T {}`,
			},
			{
				input:    `type T={}`,
				expected: `type T = {}`,
			},
			{
				input:    `type T /*aa*/{/*xx*/}// cc`,
				expected: `type T /*aa*/ {} // cc`,
			},
			{
				input: `/*aa*/type// bb
// cc
/*dd*/T // ee
/*ff*/{//
/*xx*/}// cc`,
				expected: `/*aa*/
type // bb
// cc
/*dd*/
T // ee
/*ff*/
{} // cc`,
			},
			{
				input: `type T {
			Name string
			}`,
				expected: `type T {
	Name string
}`,
			},
			{
				input: `type T {
			Foo
			}`,
				expected: `type T {
	Foo
}`,
			},
			{
				input: `type T {
			*Foo
			}`,
				expected: `type T {
	*Foo
}`,
			},
			{
				input:    testStructLitData,
				expected: expectedStructLitData,
				converter: func(s string) string {
					return strings.ReplaceAll(s, "\t", "    ")
				},
			},
		})
	})
}

//go:embed testdata/test_type_struct_group.api
var testStructGroupData string

//go:embed testdata/expected_type_struct_group.api
var expectedStructgroupData string

func TestFormat_TypeGroupStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    testStructGroupData,
			expected: expectedStructgroupData,
			converter: func(s string) string {
				return strings.ReplaceAll(s, "\t", "    ")
			},
		},
	})
}

func TestFormat_AtServerStmt(t *testing.T) {
	testRunStmt(t, []formatData{
		{
			input:    `@server()`,
			expected: ``,
		},
		{
			input: `@server(foo:foo)`,
			expected: `@server (
	foo: foo
)`,
		},
		{
			input: `@server(foo:foo quux:quux)`,
			expected: `@server (
	foo:  foo
	quux: quux
)`,
		},
		{
			input: `@server(
foo:
foo
quux:
quux
)`,
			expected: `@server (
	foo:  foo
	quux: quux
)`,
		},
		{
			input: `/*aa*/@server/*bb*/(/*cc*/foo:/**/foo /*dd*/quux:/**/quux/*ee*/)`,
			expected: `/*aa*/
@server ( /*cc*/
	foo:  foo /*dd*/
	quux: quux /*ee*/
)`,
		},
		{
			input: `/*aa*/
@server
/*bb*/(// cc
/*dd*/foo:/**/foo// ee
/*ff*/quux:/**/quux// gg
)`,
			expected: `/*aa*/
@server
/*bb*/
( // cc
	/*dd*/
	foo: foo // ee
	/*ff*/
	quux: quux // gg
)`,
		},
	})
}

func TestFormat_AtDocStmt(t *testing.T) {
	t.Run("AtDocLiteralStmt", func(t *testing.T) {
		testRunStmt(t, []formatData{
			{
				input:    `@doc ""`,
				expected: ``,
			},
			{
				input:    `@doc "foo"`,
				expected: `@doc "foo"`,
			},
			{
				input:    `@doc 		"foo"`,
				expected: `@doc "foo"`,
			},
			{
				input:    `@doc"foo"`,
				expected: `@doc "foo"`,
			},
			{
				input: `/*aa*/@doc/**/"foo"// bb`,
				expected: `/*aa*/
@doc "foo" // bb`,
			},
			{
				input: `/*aa*/
/*bb*/@doc // cc
"foo"// ee`,
				expected: `/*aa*/
/*bb*/
@doc "foo" // ee`,
			},
		})
	})
	t.Run("AtDocGroupStmt", func(t *testing.T) {
		testRunStmt(t, []formatData{
			{
				input:    `@doc()`,
				expected: ``,
			},
			{
				input: `@doc(foo:"foo")`,
				expected: `@doc (
	foo: "foo"
)`,
			},
			{
				input: `@doc(foo:"foo" bar:"bar")`,
				expected: `@doc (
	foo: "foo"
	bar: "bar"
)`,
			},
			{
				input: `@doc(foo:"foo" bar:"bar" quux:"quux")`,
				expected: `@doc (
	foo:  "foo"
	bar:  "bar"
	quux: "quux"
)`,
			},
			{
				input: `@doc(foo:"foo"
bar: "bar")`,
				expected: `@doc (
	foo: "foo"
	bar: "bar"
)`,
			},
			{
				input: `@doc(foo:"foo"// aa
bar: "bar"// bb
)`,
				expected: `@doc (
	foo: "foo" // aa
	bar: "bar" // bb
)`,
			},
			{
				input: `@doc(// aa
foo:"foo"// bb
bar: "bar"// cc
)`,
				expected: `@doc ( // aa
	foo: "foo" // bb
	bar: "bar" // cc
)`,
			},
			{
				input: `/*aa*/@doc(// bb
foo:"foo"// cc
bar: "bar"// dd
)`,
				expected: `/*aa*/
@doc ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
)`,
			},
			{
				input: `/*aa*/
@doc(// bb
foo:"foo"// cc
bar: "bar"// dd
)// ee`,
				expected: `/*aa*/
@doc ( // bb
	foo: "foo" // cc
	bar: "bar" // dd
) // ee`,
			},
			{
				input: `/*aa*/
@doc ( // bb
	/*cc*/foo: "foo" // dd
	/*ee*/bar: "bar" // ff
) // gg`,
				expected: `/*aa*/
@doc ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
			},
			{
				input: `/*aa*/
@doc/*xx*/( // bb
	/*cc*/foo:/*xx*/ "foo" // dd
	/*ee*/bar:/*xx*/ "bar" // ff
) // gg`,
				expected: `/*aa*/
@doc ( // bb
	/*cc*/
	foo: "foo" // dd
	/*ee*/
	bar: "bar" // ff
) // gg`,
			},
		})
	})
}

func TestFormat_AtHandlerStmt(t *testing.T) {
	testRunStmt(t, []formatData{
		{
			input:    `@handler foo`,
			expected: `@handler foo`,
		},
		{
			input:    `@handler 		foo`,
			expected: `@handler foo`,
		},
		{
			input: `/*aa*/@handler/**/foo// bb`,
			expected: `/*aa*/
@handler foo // bb`,
		},
		{
			input: `/*aa*/
/*bb*/@handler // cc
foo// ee`,
			expected: `/*aa*/
/*bb*/
@handler foo // ee`,
		},
	})
}

//go:embed testdata/test_service.api
var testServiceData string

//go:embed testdata/expected_service.api
var expectedServiceData string

func TestFormat_ServiceStmt(t *testing.T) {
	testRun(t, []formatData{
		{
			input:    `service foo{}`,
			expected: `service foo {}`,
		},
		{
			input:    `service foo	{	}`,
			expected: `service foo {}`,
		},
		{
			input:    `@server()service foo	{	}`,
			expected: `service foo {}`,
		},
		{
			input: `@server(foo:foo quux:quux)service foo	{	}`,
			expected: `@server (
	foo:  foo
	quux: quux
)
service foo {}`,
		},
		{
			input:    `service foo-api	{	}`,
			expected: `service foo-api {}`,
		},
		{
			input: `service foo-api	{
@doc "foo"
@handler foo
post /ping
}`,
			expected: `service foo-api {
	@doc "foo"
	@handler foo
	post /ping
}`,
		},
		{
			input: `service foo-api	{
@doc(foo: "foo" bar: "bar")
@handler foo
post /ping
}`,
			expected: `service foo-api {
	@doc (
		foo: "foo"
		bar: "bar"
	)
	@handler foo
	post /ping
}`,
		},
		{
			input: `service foo-api	{
@doc(foo: "foo" bar: "bar"
quux: "quux"
)
@handler 	foo
post 	/ping
}`,
			expected: `service foo-api {
	@doc (
		foo:  "foo"
		bar:  "bar"
		quux: "quux"
	)
	@handler foo
	post /ping
}`,
		},
		{
			input: `service
foo-api
{
@doc
(foo: "foo" bar: "bar"
quux: "quux"
)
@handler
foo
post
/aa/:bb/cc-dd/ee

@handler bar
get /bar () returns (Bar);

@handler baz
get /bar (Baz) returns ();
}`,
			expected: `service foo-api {
	@doc (
		foo:  "foo"
		bar:  "bar"
		quux: "quux"
	)
	@handler foo
	post /aa/:bb/cc-dd/ee

	@handler bar
	get /bar returns (Bar)

	@handler baz
	get /bar (Baz)
}`,
		},
		{
			input:    testServiceData,
			expected: expectedServiceData,
			converter: func(s string) string {
				return strings.ReplaceAll(s, "\t", "    ")
			},
		},
	})
}

func TestFormat_error(t *testing.T) {
	err := Source([]byte("aaa"), os.Stdout)
	assertx.Error(t, err)
}

func testRun(t *testing.T, testData []formatData) {
	for _, v := range testData {
		buffer := bytes.NewBuffer(nil)
		err := formatForUnitTest([]byte(v.input), buffer)
		assert.NoError(t, err)
		var result = buffer.String()
		if v.converter != nil {
			result = v.converter(result)
		}
		assert.Equal(t, v.expected, result)
	}
}

func testRunStmt(t *testing.T, testData []formatData) {
	for _, v := range testData {
		p := parser.New("foo.api", v.input)
		ast := p.ParseForUintTest()
		assert.NoError(t, p.CheckErrors())
		assert.True(t, len(ast.Stmts) > 0)
		one := ast.Stmts[0]
		actual := one.Format()
		if v.converter != nil {
			actual = v.converter(actual)
		}
		assert.Equal(t, v.expected, actual)
	}
}
