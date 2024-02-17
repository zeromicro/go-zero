package mapping

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/io"
)

func TestUnmarshalYamlBytes(t *testing.T) {
	var c struct {
		Name string
	}
	content := []byte(`Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalYamlBytesErrorInput(t *testing.T) {
	var c struct {
		Name string
	}
	content := []byte(`liao`)
	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesEmptyInput(t *testing.T) {
	var c struct {
		Name string
	}
	content := []byte(``)
	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",optional"`
	}
	content := []byte(`Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalYamlBytesOptionalDefault(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",optional,default=1"`
	}
	content := []byte(`Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalYamlBytesDefaultOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",default=1,optional"`
	}
	content := []byte(`Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalYamlBytesDefault(t *testing.T) {
	var c struct {
		Name string `json:",default=liao"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalYamlBytesBool(t *testing.T) {
	var c struct {
		Great bool
	}
	content := []byte(`Great: true`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.True(t, c.Great)
}

func TestUnmarshalYamlBytesInt(t *testing.T) {
	var c struct {
		Age int
	}
	content := []byte(`Age: 1`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalYamlBytesUint(t *testing.T) {
	var c struct {
		Age uint
	}
	content := []byte(`Age: 1`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, uint(1), c.Age)
}

func TestUnmarshalYamlBytesFloat(t *testing.T) {
	var c struct {
		Age float32
	}
	content := []byte(`Age: 1.5`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, float32(1.5), c.Age)
}

func TestUnmarshalYamlBytesMustInOptional(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesMustInOptionalMissedPart(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`Inner:
  There: sure`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesMustInOptionalOnlyOptionalFilled(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`Inner:
  Optional: sure`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesPartial(t *testing.T) {
	var c struct {
		Name string
		Age  float32
	}
	content := []byte(`Age: 1.5`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesStruct(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
		}
	}
	content := []byte(`Inner:
  Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalYamlBytesStructOptional(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
			Age  int `json:",optional"`
		}
	}
	content := []byte(`Inner:
  Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalYamlBytesStructPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
		}
	}
	content := []byte(`Inner:
  Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalYamlBytesStructPtrOptional(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `json:",optional"`
		}
	}
	content := []byte(`Inner:
  Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesStructPtrDefault(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `json:",default=4"`
		}
	}
	content := []byte(`Inner:
  Name: liao`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
	assert.Equal(t, 4, c.Inner.Age)
}

func TestUnmarshalYamlBytesSliceString(t *testing.T) {
	var c struct {
		Names []string
	}
	content := []byte(`Names:
- liao
- chaoxin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalYamlBytesSliceStringOptional(t *testing.T) {
	var c struct {
		Names []string
		Age   []int `json:",optional"`
	}
	content := []byte(`Names:
- liao
- chaoxin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalYamlBytesSliceStruct(t *testing.T) {
	var c struct {
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []struct {
		Name string
		Age  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %q, got %q", c.People, want)
	}
}

func TestUnmarshalYamlBytesSliceStructOptional(t *testing.T) {
	var c struct {
		People []struct {
			Name   string
			Age    int
			Emails []string `json:",optional"`
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []struct {
		Name   string
		Age    int
		Emails []string `json:",optional"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %q, got %q", c.People, want)
	}
}

func TestUnmarshalYamlBytesSliceStructPtr(t *testing.T) {
	var c struct {
		People []*struct {
			Name string
			Age  int
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []*struct {
		Name string
		Age  int
	}{
		{"liao", 1},
		{"chaoxin", 2},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %v, got %v", c.People, want)
	}
}

func TestUnmarshalYamlBytesSliceStructPtrOptional(t *testing.T) {
	var c struct {
		People []*struct {
			Name   string
			Age    int
			Emails []string `json:",optional"`
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []*struct {
		Name   string
		Age    int
		Emails []string `json:",optional"`
	}{
		{"liao", 1, nil},
		{"chaoxin", 2, nil},
	}
	if !reflect.DeepEqual(c.People, want) {
		t.Fatalf("want %v, got %v", c.People, want)
	}
}

func TestUnmarshalYamlBytesSliceStructPtrPartial(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesSliceStructPtrDefault(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string `json:",default=chaoxin@liao.com"`
		}
	}
	content := []byte(`People:
- Name: liao
  Age: 1
- Name: chaoxin
  Age: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))

	want := []*struct {
		Name  string
		Age   int
		Email string
	}{
		{"liao", 1, "chaoxin@liao.com"},
		{"chaoxin", 2, "chaoxin@liao.com"},
	}

	for i := range c.People {
		actual := c.People[i]
		expect := want[i]
		assert.Equal(t, expect.Age, actual.Age)
		assert.Equal(t, expect.Email, actual.Email)
		assert.Equal(t, expect.Name, actual.Name)
	}
}

func TestUnmarshalYamlBytesSliceStringPartial(t *testing.T) {
	var c struct {
		Names []string
		Age   int
	}
	content := []byte(`Age: 1`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesSliceStructPartial(t *testing.T) {
	var c struct {
		Group  string
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`Group: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesInnerAnonymousPartial(t *testing.T) {
	type (
		Deep struct {
			A string
			B string `json:",optional"`
		}
		Inner struct {
			Deep
			InnerV string `json:",optional"`
		}
	)

	var c struct {
		Value Inner `json:",optional"`
	}
	content := []byte(`Value:
  InnerV: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesStructPartial(t *testing.T) {
	var c struct {
		Group  string
		Person struct {
			Name string
			Age  int
		}
	}
	content := []byte(`Group: chaoxin`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesEmptyMap(t *testing.T) {
	var c struct {
		Persons map[string]int `json:",optional"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Empty(t, c.Persons)
}

func TestUnmarshalYamlBytesMap(t *testing.T) {
	var c struct {
		Persons map[string]int
	}
	content := []byte(`Persons:
  first: 1
  second: 2`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 2, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"])
	assert.Equal(t, 2, c.Persons["second"])
}

func TestUnmarshalYamlBytesMapStruct(t *testing.T) {
	var c struct {
		Persons map[string]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
    ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalYamlBytesMapStructPtr(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
    ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalYamlBytesMapStructMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string
		}
	}
	content := []byte(`Persons:
  first:
    ID: 1`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesMapStructOptional(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
    ID: 1`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
}

func TestUnmarshalYamlBytesMapStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
  - ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalYamlBytesMapEmptyStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first: []`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalYamlBytesMapStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
  - ID: 1
    name: kevin`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalYamlBytesMapEmptyStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first: []`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalYamlBytesMapStructPtrSliceMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string
		}
	}
	content := []byte(`Persons:
  first:
  - ID: 1`)

	assert.NotNil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlBytesMapStructPtrSliceOptional(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`Persons:
  first:
  - ID: 1`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
}

func TestUnmarshalYamlStructOptional(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Hosts []string
			Key   string
		} `json:",optional"`
	}
	content := []byte(`Name: kevin`)

	err := UnmarshalYamlBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
}

func TestUnmarshalYamlStructLowerCase(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Key string
		} `json:"etcd"`
	}
	content := []byte(`Name: kevin
etcd:
  Key: the key`)

	err := UnmarshalYamlBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
	assert.Equal(t, "the key", c.Etcd.Key)
}

func TestUnmarshalYamlWithStructAllOptionalWithEmpty(t *testing.T) {
	var c struct {
		Inner struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlWithStructAllOptionalPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlWithStructOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   Inner `json:",optional"`
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.In.Must)
}

func TestUnmarshalYamlWithStructPtrOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   *Inner `json:",optional"`
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.In)
}

func TestUnmarshalYamlWithStructAllOptionalAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlWithStructAllOptionalAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
}

func TestUnmarshalYamlWithStructAllOptionalProvoidedAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`Else: sure
Optional: optional`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalYamlWithStructAllOptionalProvoidedAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`Else: sure
Optional: optional`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalYamlWithStructAnonymous(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`Else: sure
Must: must`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalYamlWithStructAnonymousPtr(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`Else: sure
Must: must`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalYamlWithStructAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner `json:",optional"`
		Else  string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.Must)
}

func TestUnmarshalYamlWithStructPtrAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner `json:",optional"`
		Else   string
	}
	content := []byte(`Else: sure`)

	assert.Nil(t, UnmarshalYamlBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.Inner)
}

func TestUnmarshalYamlWithZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `json:"negative"`
		Int    int    `json:"int"`
		String string `json:"string"`
	}
	content := []byte(`negative: false
int: 0
string: ""`)

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalYamlBytes(content, &in))
	ast.False(in.False)
	ast.Equal(0, in.Int)
	ast.Equal("", in.String)
}

func TestUnmarshalYamlBytesError(t *testing.T) {
	payload := `abcd:
- cdef`
	var v struct {
		Any []string `json:"abcd"`
	}

	err := UnmarshalYamlBytes([]byte(payload), &v)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(v.Any))
	assert.Equal(t, "cdef", v.Any[0])
}

func TestUnmarshalYamlReaderError(t *testing.T) {
	var v struct {
		Any string
	}

	reader := strings.NewReader(`abcd: cdef`)
	err := UnmarshalYamlReader(reader, &v)
	assert.NotNil(t, err)

	reader = strings.NewReader("foo")
	assert.Error(t, UnmarshalYamlReader(reader, &v))
}

func TestUnmarshalYamlBadReader(t *testing.T) {
	var v struct {
		Any string
	}

	err := UnmarshalYamlReader(new(badReader), &v)
	assert.NotNil(t, err)
}

func TestUnmarshalYamlMapBool(t *testing.T) {
	text := `machine:
  node1: true
  node2: true
  node3: true
`
	var v struct {
		Machine map[string]bool `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.True(t, v.Machine["node1"])
	assert.True(t, v.Machine["node2"])
	assert.True(t, v.Machine["node3"])
}

func TestUnmarshalYamlMapInt(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]int `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, 1, v.Machine["node1"])
	assert.Equal(t, 2, v.Machine["node2"])
	assert.Equal(t, 3, v.Machine["node3"])
}

func TestUnmarshalYamlMapByte(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]byte `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, byte(1), v.Machine["node1"])
	assert.Equal(t, byte(2), v.Machine["node2"])
	assert.Equal(t, byte(3), v.Machine["node3"])
}

func TestUnmarshalYamlMapRune(t *testing.T) {
	text := `machine:
  node1: 1
  node2: 2
  node3: 3
`
	var v struct {
		Machine map[string]rune `json:"machine,optional"`
	}
	reader := strings.NewReader(text)
	assert.Nil(t, UnmarshalYamlReader(reader, &v))
	assert.Equal(t, rune(1), v.Machine["node1"])
	assert.Equal(t, rune(2), v.Machine["node2"])
	assert.Equal(t, rune(3), v.Machine["node3"])
}

func TestUnmarshalYamlStringOfInt(t *testing.T) {
	text := `password: 123456`
	var v struct {
		Password string `json:"password"`
	}
	reader := strings.NewReader(text)
	assert.Error(t, UnmarshalYamlReader(reader, &v))
}

func TestUnmarshalYamlBadInput(t *testing.T) {
	var v struct {
		Any string
	}
	assert.Error(t, UnmarshalYamlBytes([]byte("':foo"), &v))
}

type badReader struct{}

func (b *badReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrLimitReached
}
