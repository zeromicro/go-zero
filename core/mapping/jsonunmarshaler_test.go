package mapping

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBytes(t *testing.T) {
	var c struct {
		Name string
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",optional"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesOptionalDefault(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",optional,default=1"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesDefaultOptional(t *testing.T) {
	var c struct {
		Name string
		Age  int `json:",default=1,optional"`
	}
	content := []byte(`{"Name": "liao"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesDefault(t *testing.T) {
	var c struct {
		Name string `json:",default=liao"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Name)
}

func TestUnmarshalBytesBool(t *testing.T) {
	var c struct {
		Great bool
	}
	content := []byte(`{"Great": true}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.True(t, c.Great)
}

func TestUnmarshalBytesInt(t *testing.T) {
	var c struct {
		Age int
	}
	content := []byte(`{"Age": 1}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, c.Age)
}

func TestUnmarshalBytesUint(t *testing.T) {
	var c struct {
		Age uint
	}
	content := []byte(`{"Age": 1}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, uint(1), c.Age)
}

func TestUnmarshalBytesFloat(t *testing.T) {
	var c struct {
		Age float32
	}
	content := []byte(`{"Age": 1.5}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, float32(1.5), c.Age)
}

func TestUnmarshalBytesMustInOptional(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`{}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMustInOptionalMissedPart(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`{"Inner": {"There": "sure"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMustInOptionalOnlyOptionalFilled(t *testing.T) {
	var c struct {
		Inner struct {
			There    string
			Must     string
			Optional string `json:",optional"`
		} `json:",optional"`
	}
	content := []byte(`{"Inner": {"Optional": "sure"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesNil(t *testing.T) {
	var c struct {
		Int int64 `json:"int,optional"`
	}
	content := []byte(`{"int":null}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, int64(0), c.Int)
}

func TestUnmarshalBytesNilSlice(t *testing.T) {
	var c struct {
		Ints []int64 `json:"ints"`
	}
	content := []byte(`{"ints":[null]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 0, len(c.Ints))
}

func TestUnmarshalBytesPartial(t *testing.T) {
	var c struct {
		Name string
		Age  float32
	}
	content := []byte(`{"Age": 1.5}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStruct(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructOptional(t *testing.T) {
	var c struct {
		Inner struct {
			Name string
			Age  int `json:",optional"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
}

func TestUnmarshalBytesStructPtrOptional(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `json:",optional"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStructPtrDefault(t *testing.T) {
	var c struct {
		Inner *struct {
			Name string
			Age  int `json:",default=4"`
		}
	}
	content := []byte(`{"Inner": {"Name": "liao"}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "liao", c.Inner.Name)
	assert.Equal(t, 4, c.Inner.Age)
}

func TestUnmarshalBytesSliceString(t *testing.T) {
	var c struct {
		Names []string
	}
	content := []byte(`{"Names": ["liao", "chaoxin"]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalBytesSliceStringOptional(t *testing.T) {
	var c struct {
		Names []string
		Age   []int `json:",optional"`
	}
	content := []byte(`{"Names": ["liao", "chaoxin"]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

	want := []string{"liao", "chaoxin"}
	if !reflect.DeepEqual(c.Names, want) {
		t.Fatalf("want %q, got %q", c.Names, want)
	}
}

func TestUnmarshalBytesSliceStruct(t *testing.T) {
	var c struct {
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

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

func TestUnmarshalBytesSliceStructOptional(t *testing.T) {
	var c struct {
		People []struct {
			Name   string
			Age    int
			Emails []string `json:",optional"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

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

func TestUnmarshalBytesSliceStructPtr(t *testing.T) {
	var c struct {
		People []*struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

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

func TestUnmarshalBytesSliceStructPtrOptional(t *testing.T) {
	var c struct {
		People []*struct {
			Name   string
			Age    int
			Emails []string `json:",optional"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

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

func TestUnmarshalBytesSliceStructPtrPartial(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesSliceStructPtrDefault(t *testing.T) {
	var c struct {
		People []*struct {
			Name  string
			Age   int
			Email string `json:",default=chaoxin@liao.com"`
		}
	}
	content := []byte(`{"People": [{"Name": "liao", "Age": 1}, {"Name": "chaoxin", "Age": 2}]}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))

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

func TestUnmarshalBytesSliceStringPartial(t *testing.T) {
	var c struct {
		Names []string
		Age   int
	}
	content := []byte(`{"Age": 1}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesSliceStructPartial(t *testing.T) {
	var c struct {
		Group  string
		People []struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"Group": "chaoxin"}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesInnerAnonymousPartial(t *testing.T) {
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
	content := []byte(`{"Value": {"InnerV": "chaoxin"}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesStructPartial(t *testing.T) {
	var c struct {
		Group  string
		Person struct {
			Name string
			Age  int
		}
	}
	content := []byte(`{"Group": "chaoxin"}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesEmptyMap(t *testing.T) {
	var c struct {
		Persons map[string]int `json:",optional"`
	}
	content := []byte(`{"Persons": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Empty(t, c.Persons)
}

func TestUnmarshalBytesMap(t *testing.T) {
	var c struct {
		Persons map[string]int
	}
	content := []byte(`{"Persons": {"first": 1, "second": 2}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 2, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"])
	assert.Equal(t, 2, c.Persons["second"])
}

func TestUnmarshalBytesMapStruct(t *testing.T) {
	var c struct {
		Persons map[string]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": {"ID": 1, "name": "kevin"}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalBytesMapStructPtr(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": {"ID": 1, "name": "kevin"}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
	assert.Equal(t, "kevin", c.Persons["first"].Name)
}

func TestUnmarshalBytesMapStructMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string
		}
	}
	content := []byte(`{"Persons": {"first": {"ID": 1}}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMapStructOptional(t *testing.T) {
	var c struct {
		Persons map[string]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": {"ID": 1}}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"].ID)
}

func TestUnmarshalBytesMapEmptyStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": []}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalBytesMapStructSlice(t *testing.T) {
	var c struct {
		Persons map[string][]struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"ID": 1, "name": "kevin"}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalBytesMapEmptyStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": []}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Empty(t, c.Persons["first"])
}

func TestUnmarshalBytesMapStructPtrSlice(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"ID": 1, "name": "kevin"}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
	assert.Equal(t, "kevin", c.Persons["first"][0].Name)
}

func TestUnmarshalBytesMapStructPtrSliceMissingPartial(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string
		}
	}
	content := []byte(`{"Persons": {"first": [{"ID": 1}]}}`)

	assert.NotNil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalBytesMapStructPtrSliceOptional(t *testing.T) {
	var c struct {
		Persons map[string][]*struct {
			ID   int
			Name string `json:"name,optional"`
		}
	}
	content := []byte(`{"Persons": {"first": [{"ID": 1}]}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, 1, len(c.Persons))
	assert.Equal(t, 1, c.Persons["first"][0].ID)
}

func TestUnmarshalStructOptional(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Hosts []string
			Key   string
		} `json:",optional"`
	}
	content := []byte(`{"Name": "kevin"}`)

	err := UnmarshalJsonBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
}

func TestUnmarshalStructLowerCase(t *testing.T) {
	var c struct {
		Name string
		Etcd struct {
			Key string
		} `json:"etcd"`
	}
	content := []byte(`{"Name": "kevin", "etcd": {"Key": "the key"}}`)

	err := UnmarshalJsonBytes(content, &c)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", c.Name)
	assert.Equal(t, "the key", c.Etcd.Key)
}

func TestUnmarshalWithStructAllOptionalWithEmpty(t *testing.T) {
	var c struct {
		Inner struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	content := []byte(`{"Else": "sure", "Inner": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalPtr(t *testing.T) {
	var c struct {
		Inner *struct {
			Optional string `json:",optional"`
		}
		Else string
	}
	content := []byte(`{"Else": "sure", "Inner": {}}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   Inner `json:",optional"`
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.In.Must)
}

func TestUnmarshalWithStructPtrOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		In   *Inner `json:",optional"`
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.In)
}

func TestUnmarshalWithStructAllOptionalAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
}

func TestUnmarshalWithStructAllOptionalProvoidedAnonymous(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Optional": "optional"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalWithStructAllOptionalProvoidedAnonymousPtr(t *testing.T) {
	type Inner struct {
		Optional string `json:",optional"`
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Optional": "optional"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "optional", c.Optional)
}

func TestUnmarshalWithStructAnonymous(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Must": "must"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalWithStructAnonymousPtr(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner
		Else string
	}
	content := []byte(`{"Else": "sure", "Must": "must"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "must", c.Must)
}

func TestUnmarshalWithStructAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		Inner `json:",optional"`
		Else  string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Equal(t, "", c.Must)
}

func TestUnmarshalWithStructPtrAnonymousOptional(t *testing.T) {
	type Inner struct {
		Must string
	}

	var c struct {
		*Inner `json:",optional"`
		Else   string
	}
	content := []byte(`{"Else": "sure"}`)

	assert.Nil(t, UnmarshalJsonBytes(content, &c))
	assert.Equal(t, "sure", c.Else)
	assert.Nil(t, c.Inner)
}

func TestUnmarshalWithZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `json:"no"`
		Int    int    `json:"int"`
		String string `json:"string"`
	}
	content := []byte(`{"no": false, "int": 0, "string": ""}`)
	reader := bytes.NewReader(content)

	var in inner
	ast := assert.New(t)
	ast.Nil(UnmarshalJsonReader(reader, &in))
	ast.False(in.False)
	ast.Equal(0, in.Int)
	ast.Equal("", in.String)
}

func TestUnmarshalBytesError(t *testing.T) {
	payload := `[{"abcd": "cdef"}]`
	var v struct {
		Any string
	}

	err := UnmarshalJsonBytes([]byte(payload), &v)
	assert.Equal(t, errTypeMismatch, err)
}

func TestUnmarshalReaderError(t *testing.T) {
	payload := `[{"abcd": "cdef"}]`
	reader := strings.NewReader(payload)
	var v struct {
		Any string
	}

	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(reader, &v))
}

func TestUnmarshalMap(t *testing.T) {
	t.Run("nil map and valid", func(t *testing.T) {
		var m map[string]any
		var v struct {
			Any string `json:",optional"`
		}

		err := UnmarshalJsonMap(m, &v)
		assert.Nil(t, err)
		assert.True(t, len(v.Any) == 0)
	})

	t.Run("empty map but not valid", func(t *testing.T) {
		m := map[string]any{}
		var v struct {
			Any string
		}

		err := UnmarshalJsonMap(m, &v)
		assert.NotNil(t, err)
	})

	t.Run("empty map and valid", func(t *testing.T) {
		m := map[string]any{}
		var v struct {
			Any string `json:",optional"`
		}

		err := UnmarshalJsonMap(m, &v, WithCanonicalKeyFunc(func(s string) string {
			return s
		}))
		assert.Nil(t, err)
		assert.True(t, len(v.Any) == 0)
	})

	t.Run("valid map", func(t *testing.T) {
		m := map[string]any{
			"Any": "foo",
		}
		var v struct {
			Any string
		}

		err := UnmarshalJsonMap(m, &v)
		assert.Nil(t, err)
		assert.Equal(t, "foo", v.Any)
	})
}

func TestUnmarshalJsonArray(t *testing.T) {
	var v []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	body := `[{"name":"kevin", "age": 18}]`
	assert.NoError(t, UnmarshalJsonBytes([]byte(body), &v))
	assert.Equal(t, 1, len(v))
	assert.Equal(t, "kevin", v[0].Name)
	assert.Equal(t, 18, v[0].Age)
}

func TestUnmarshalJsonBytesError(t *testing.T) {
	var v []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	assert.Error(t, UnmarshalJsonBytes([]byte((``)), &v))
	assert.Error(t, UnmarshalJsonReader(strings.NewReader(``), &v))
}
