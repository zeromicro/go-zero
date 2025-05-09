package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestParseForm(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		var v struct {
			Name    string  `form:"name"`
			Age     int     `form:"age"`
			Percent float64 `form:"percent,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=hello&age=18&percent=3.4",
			http.NoBody)
		assert.Nil(t, err)
		assert.Nil(t, Parse(r, &v))
		assert.Equal(t, "hello", v.Name)
		assert.Equal(t, 18, v.Age)
		assert.Equal(t, 3.4, v.Percent)
	})

	t.Run("no value", func(t *testing.T) {
		var v struct {
			NoValue string `form:"noValue,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=hello&age=18&percent=3.4&statuses=try&statuses=done&singleValue=one",
			http.NoBody)
		assert.Nil(t, err)
		assert.Nil(t, Parse(r, &v))
		assert.Equal(t, 0, len(v.NoValue))
	})

	t.Run("slice with one value on array format", func(t *testing.T) {
		var v struct {
			Names string `form:"names"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?names=1,2,3",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.Equal(t, "1,2,3", v.Names)
		}
	})
}

func TestParseFormArray(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		var v struct {
			Name    []string  `form:"name"`
			Age     []int     `form:"age"`
			Percent []float64 `form:"percent,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=hello&name=world&age=18&age=19&percent=3.4&percent=4.5",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"hello", "world"}, v.Name)
			assert.ElementsMatch(t, []int{18, 19}, v.Age)
			assert.ElementsMatch(t, []float64{3.4, 4.5}, v.Percent)
		}
	})

	t.Run("slice with single value", func(t *testing.T) {
		var v struct {
			Name    []string  `form:"name"`
			Age     []int     `form:"age"`
			Percent []float64 `form:"percent,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=hello&age=18&percent=3.4",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"hello"}, v.Name)
			assert.ElementsMatch(t, []int{18}, v.Age)
			assert.ElementsMatch(t, []float64{3.4}, v.Percent)
		}
	})

	t.Run("slice with empty", func(t *testing.T) {
		var v struct {
			Name []string `form:"name,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{}, v.Name)
		}
	})

	t.Run("slice with empty", func(t *testing.T) {
		var v struct {
			Name []string `form:"name,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.Empty(t, v.Name)
		}
	})

	t.Run("slice with empty and non-empty", func(t *testing.T) {
		var v struct {
			Name []string `form:"name"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?name=&name=1",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"1"}, v.Name)
		}
	})

	t.Run("slice with one value on array format", func(t *testing.T) {
		var v struct {
			Names []string `form:"names"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?names=1,2,3",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"1,2,3"}, v.Names)
		}
	})

	t.Run("slice with one value on combined array format", func(t *testing.T) {
		var v struct {
			Names []string `form:"names"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?names=[1,2,3]&names=4",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"[1,2,3]", "4"}, v.Names)
		}
	})

	t.Run("slice with one value on integer array format", func(t *testing.T) {
		var v struct {
			Numbers []int `form:"numbers"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?numbers=1,2,3",
			http.NoBody)
		assert.NoError(t, err)
		assert.Error(t, Parse(r, &v))
	})

	t.Run("slice with one value on array format brackets", func(t *testing.T) {
		var v struct {
			Names []string `form:"names"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?names[]=1&names[]=2&names[]=3",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"1", "2", "3"}, v.Names)
		}
	})

	t.Run("slice with one empty value on integer array format", func(t *testing.T) {
		var v struct {
			Numbers []int `form:"numbers,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?numbers=",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.Empty(t, v.Numbers)
		}
	})

	t.Run("slice with one value on integer array format", func(t *testing.T) {
		var v struct {
			Numbers []int `form:"numbers,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?numbers=&numbers=2",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []int{2}, v.Numbers)
		}
	})

	t.Run("slice with one empty value on float64 array format", func(t *testing.T) {
		var v struct {
			Numbers []float64 `form:"numbers,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?numbers=",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.Empty(t, v.Numbers)
		}
	})

	t.Run("slice with one value on float64 array format", func(t *testing.T) {
		var v struct {
			Numbers []float64 `form:"numbers,optional"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?numbers=&numbers=2",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []float64{2}, v.Numbers)
		}
	})

	t.Run("slice with one value", func(t *testing.T) {
		var v struct {
			Codes []string `form:"codes"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?codes=aaa,bbb,ccc",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"aaa,bbb,ccc"}, v.Codes)
		}
	})

	t.Run("slice with multiple values", func(t *testing.T) {
		var v struct {
			Codes []string `form:"codes,arrayComma=false"`
		}

		r, err := http.NewRequest(
			http.MethodGet,
			"/a?codes=aaa,bbb,ccc&codes=ccc,ddd,eee",
			http.NoBody)
		assert.NoError(t, err)
		if assert.NoError(t, Parse(r, &v)) {
			assert.ElementsMatch(t, []string{"aaa,bbb,ccc", "ccc,ddd,eee"}, v.Codes)
		}
	})
}

func TestParseForm_Error(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	r := httptest.NewRequest(http.MethodGet, "/a?name=hello;", http.NoBody)
	assert.NotNil(t, ParseForm(r, &v))
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		expect map[string]string
	}{
		{
			name:   "empty",
			value:  "",
			expect: map[string]string{},
		},
		{
			name:   "regular",
			value:  "key=value",
			expect: map[string]string{"key": "value"},
		},
		{
			name:   "next empty",
			value:  "key=value;",
			expect: map[string]string{"key": "value"},
		},
		{
			name:   "regular",
			value:  "key=value;foo",
			expect: map[string]string{"key": "value"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			m := ParseHeader(test.value)
			assert.EqualValues(t, test.expect, m)
		})
	}
}

func TestParsePath(t *testing.T) {
	var v struct {
		Name string `path:"name"`
		Age  int    `path:"age"`
	}

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r = pathvar.WithVars(r, map[string]string{
		"name": "foo",
		"age":  "18",
	})
	err := Parse(r, &v)
	assert.Nil(t, err)
	assert.Equal(t, "foo", v.Name)
	assert.Equal(t, 18, v.Age)
}

func TestParsePath_Error(t *testing.T) {
	var v struct {
		Name string `path:"name"`
		Age  int    `path:"age"`
	}

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r = pathvar.WithVars(r, map[string]string{
		"name": "foo",
	})
	assert.NotNil(t, Parse(r, &v))
}

func TestParseFormOutOfRange(t *testing.T) {
	var v struct {
		Age int `form:"age,range=[10:20)"`
	}

	tests := []struct {
		url  string
		pass bool
	}{
		{
			url:  "/a?age=5",
			pass: false,
		},
		{
			url:  "/a?age=10",
			pass: true,
		},
		{
			url:  "/a?age=15",
			pass: true,
		},
		{
			url:  "/a?age=20",
			pass: false,
		},
		{
			url:  "/a?age=28",
			pass: false,
		},
	}

	for _, test := range tests {
		r, err := http.NewRequest(http.MethodGet, test.url, http.NoBody)
		assert.Nil(t, err)

		err = Parse(r, &v)
		if test.pass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestParseMultipartForm(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	body := strings.Replace(`----------------------------220477612388154780019383
Content-Disposition: form-data; name="name"

kevin
----------------------------220477612388154780019383
Content-Disposition: form-data; name="age"

18
----------------------------220477612388154780019383--`, "\n", "\r\n", -1)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set(ContentType, "multipart/form-data; boundary=--------------------------220477612388154780019383")

	assert.Nil(t, Parse(r, &v))
	assert.Equal(t, "kevin", v.Name)
	assert.Equal(t, 18, v.Age)
}

func TestParseMultipartFormWrongBoundary(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	body := strings.Replace(`----------------------------22047761238815478001938
Content-Disposition: form-data; name="name"

kevin
----------------------------22047761238815478001938
Content-Disposition: form-data; name="age"

18
----------------------------22047761238815478001938--`, "\n", "\r\n", -1)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set(ContentType, "multipart/form-data; boundary=--------------------------220477612388154780019383")

	assert.NotNil(t, Parse(r, &v))
}

func TestParseJsonBody(t *testing.T) {
	t.Run("has body", func(t *testing.T) {
		var v struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		body := `{"name":"kevin", "age": 18}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(ContentType, header.ContentTypeJson)

		if assert.NoError(t, Parse(r, &v)) {
			assert.Equal(t, "kevin", v.Name)
			assert.Equal(t, 18, v.Age)
		}
	})

	t.Run("bad body", func(t *testing.T) {
		var v struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		body := `{"name":"kevin", "ag": 18}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(ContentType, header.ContentTypeJson)

		assert.Error(t, Parse(r, &v))
	})

	t.Run("hasn't body", func(t *testing.T) {
		var v struct {
			Name string `json:"name,optional"`
			Age  int    `json:"age,optional"`
		}

		r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		assert.Nil(t, Parse(r, &v))
		assert.Equal(t, "", v.Name)
		assert.Equal(t, 0, v.Age)
	})

	t.Run("array body", func(t *testing.T) {
		var v []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		body := `[{"name":"kevin", "age": 18}]`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(ContentType, header.ContentTypeJson)

		assert.NoError(t, Parse(r, &v))
		assert.Equal(t, 1, len(v))
		assert.Equal(t, "kevin", v[0].Name)
		assert.Equal(t, 18, v[0].Age)
	})

	t.Run("form and array body", func(t *testing.T) {
		var v []struct {
			// we can only ignore the form tag,
			// because if the value is a slice, it should be in the body,
			// but it's hard to detect when we treat it as a json body.
			Product string `form:"product"`
			Name    string `json:"name"`
			Age     int    `json:"age"`
		}

		body := `[{"name":"apple", "age": 18}]`
		r := httptest.NewRequest(http.MethodPost, "/a?product=tree", strings.NewReader(body))
		r.Header.Set(ContentType, header.ContentTypeJson)

		assert.NoError(t, Parse(r, &v))
		assert.Equal(t, 1, len(v))
		assert.Equal(t, "apple", v[0].Name)
		assert.Equal(t, 18, v[0].Age)
	})

	t.Run("bytes field", func(t *testing.T) {
		type v struct {
			Signature []byte `json:"signature,optional"`
		}
		v1 := v{
			Signature: []byte{0x01, 0xff, 0x00},
		}
		body, _ := json.Marshal(v1)
		t.Logf("body:%s", string(body))
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		r.Header.Set(ContentType, header.ContentTypeJson)
		var v2 v
		err := ParseJsonBody(r, &v2)
		if assert.NoError(t, err) {
			assert.Greater(t, len(v2.Signature), 0)
		}
		t.Logf("%x", v2.Signature)
		assert.EqualValues(t, v1.Signature, v2.Signature)
	})
}

func TestParseRequired(t *testing.T) {
	v := struct {
		Name    string  `form:"name"`
		Percent float64 `form:"percent"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "/a?name=hello", http.NoBody)
	assert.Nil(t, err)
	assert.NotNil(t, Parse(r, &v))
}

func TestParseOptions(t *testing.T) {
	v := struct {
		Position int8 `form:"pos,options=1|2"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "/a?pos=4", http.NoBody)
	assert.Nil(t, err)
	assert.NotNil(t, Parse(r, &v))
}

func TestParseHeaders(t *testing.T) {
	type AnonymousStruct struct {
		XRealIP string `header:"x-real-ip"`
		Accept  string `header:"Accept,optional"`
	}
	v := struct {
		Name          string   `header:"name,optional"`
		Percent       string   `header:"percent"`
		Addrs         []string `header:"addrs"`
		XForwardedFor string   `header:"X-Forwarded-For,optional"`
		AnonymousStruct
	}{}
	request, err := http.NewRequest("POST", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("name", "chenquan")
	request.Header.Set("percent", "1")
	request.Header.Add("addrs", "addr1")
	request.Header.Add("addrs", "addr2")
	request.Header.Add("X-Forwarded-For", "10.0.10.11")
	request.Header.Add("x-real-ip", "10.0.11.10")
	request.Header.Add("Accept", header.ContentTypeJson)
	err = ParseHeaders(request, &v)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "chenquan", v.Name)
	assert.Equal(t, "1", v.Percent)
	assert.Equal(t, []string{"addr1", "addr2"}, v.Addrs)
	assert.Equal(t, "10.0.10.11", v.XForwardedFor)
	assert.Equal(t, "10.0.11.10", v.XRealIP)
	assert.Equal(t, header.ContentTypeJson, v.Accept)
}

func TestParseHeaders_Error(t *testing.T) {
	v := struct {
		Name string `header:"name"`
		Age  int    `header:"age"`
	}{}

	r := httptest.NewRequest("POST", "/", http.NoBody)
	r.Header.Set("name", "foo")
	assert.NotNil(t, Parse(r, &v))
}

func TestParseWithValidator(t *testing.T) {
	SetValidator(mockValidator{})
	defer SetValidator(mockValidator{nop: true})

	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "/a?name=hello&age=18&percent=3.4", http.NoBody)
	assert.Nil(t, err)
	if assert.NoError(t, Parse(r, &v)) {
		assert.Equal(t, "hello", v.Name)
		assert.Equal(t, 18, v.Age)
		assert.Equal(t, 3.4, v.Percent)
	}
}

func TestParseWithValidatorWithError(t *testing.T) {
	SetValidator(mockValidator{})
	defer SetValidator(mockValidator{nop: true})

	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "/a?name=world&age=18&percent=3.4", http.NoBody)
	assert.Nil(t, err)
	assert.Error(t, Parse(r, &v))
}

func TestParseWithValidatorRequest(t *testing.T) {
	SetValidator(mockValidator{})
	defer SetValidator(mockValidator{nop: true})

	var v mockRequest
	r, err := http.NewRequest(http.MethodGet, "/a?&age=18", http.NoBody)
	assert.Nil(t, err)
	assert.Error(t, Parse(r, &v))
}

func TestParseFormWithDot(t *testing.T) {
	var v struct {
		Age int `form:"user.age"`
	}
	r, err := http.NewRequest(http.MethodGet, "/a?user.age=18", http.NoBody)
	assert.Nil(t, err)
	assert.NoError(t, Parse(r, &v))
	assert.Equal(t, 18, v.Age)
}

func TestParsePathWithDot(t *testing.T) {
	var v struct {
		Name string `path:"name.val"`
		Age  int    `path:"age.val"`
	}

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r = pathvar.WithVars(r, map[string]string{
		"name.val": "foo",
		"age.val":  "18",
	})
	err := Parse(r, &v)
	assert.Nil(t, err)
	assert.Equal(t, "foo", v.Name)
	assert.Equal(t, 18, v.Age)
}

func TestParseWithFloatPtr(t *testing.T) {
	t.Run("has float32 pointer", func(t *testing.T) {
		var v struct {
			WeightFloat32 *float32 `json:"weightFloat32,optional"`
		}
		body := `{"weightFloat32": 3.2}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(ContentType, header.ContentTypeJson)

		if assert.NoError(t, Parse(r, &v)) {
			assert.Equal(t, float32(3.2), *v.WeightFloat32)
		}
	})
}

func TestParseWithEscapedParams(t *testing.T) {
	t.Run("escaped", func(t *testing.T) {
		var v struct {
			Dev string `form:"dev"`
		}
		r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/api/v2/dev/test?dev=se205%5fy1205%5fj109%26verRelease=v01%26iid1=863494061186673%26iid2=863494061186681%26mcc=636%26mnc=1", http.NoBody)
		if assert.NoError(t, Parse(r, &v)) {
			assert.Equal(t, "se205_y1205_j109&verRelease=v01&iid1=863494061186673&iid2=863494061186681&mcc=636&mnc=1", v.Dev)
		}
	})
}

func TestCustomUnmarshalerStructRequest(t *testing.T) {
	reqBody := `{"name": "hello"}`
	r := httptest.NewRequest(http.MethodPost, "/a", bytes.NewReader([]byte(reqBody)))
	r.Header.Set(ContentType, JsonContentType)
	v := struct {
		Foo *mockUnmarshaler `json:"name"`
	}{}
	assert.Nil(t, Parse(r, &v))
	assert.Equal(t, "hello", v.Foo.Name)
}

func TestParseJsonStringRequest(t *testing.T) {
	type GoodsInfo struct {
		Sku int64 `json:"sku,optional"`
	}

	type GetReq struct {
		GoodsList []*GoodsInfo `json:"goods_list"`
	}

	input := `{"goods_list":"[{\"sku\":11},{\"sku\":22}]"}`
	r := httptest.NewRequest(http.MethodPost, "/a", strings.NewReader(input))
	r.Header.Set(ContentType, JsonContentType)
	var v GetReq
	assert.NotPanics(t, func() {
		assert.NoError(t, Parse(r, &v))
		assert.Equal(t, 2, len(v.GoodsList))
		assert.ElementsMatch(t, []int64{11, 22}, []int64{v.GoodsList[0].Sku, v.GoodsList[1].Sku})
	})
}

type valid1 struct{}

func (v valid1) Validate(*http.Request, any) error { return nil }

type valid2 struct{}

func (v valid2) Validate(*http.Request, any) error { return nil }

func TestSetValidatorTwice(t *testing.T) {
	// panic: sync/atomic: store of inconsistently typed value into Value
	assert.NotPanics(t, func() {
		SetValidator(valid1{})
		SetValidator(valid2{})
	})
}

func BenchmarkParseRaw(b *testing.B) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", http.NoBody)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		v := struct {
			Name    string  `form:"name"`
			Age     int     `form:"age"`
			Percent float64 `form:"percent,optional"`
		}{}

		v.Name = r.FormValue("name")
		v.Age, err = strconv.Atoi(r.FormValue("age"))
		if err != nil {
			b.Fatal(err)
		}
		v.Percent, err = strconv.ParseFloat(r.FormValue("percent"), 64)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseAuto(b *testing.B) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", http.NoBody)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		v := struct {
			Name    string  `form:"name"`
			Age     int     `form:"age"`
			Percent float64 `form:"percent,optional"`
		}{}

		if err = Parse(r, &v); err != nil {
			b.Fatal(err)
		}
	}
}

type mockValidator struct {
	nop bool
}

func (m mockValidator) Validate(r *http.Request, data any) error {
	if m.nop {
		return nil
	}

	if r.URL.Path == "/a" {
		val := reflect.ValueOf(data).Elem().FieldByName("Name").String()
		if val != "hello" {
			return errors.New("name is not hello")
		}
	}

	return nil
}

type mockRequest struct {
	Name string `json:"name,optional"`
}

func (m mockRequest) Validate() error {
	if m.Name != "hello" {
		return errors.New("name is not hello")
	}

	return nil
}

type mockUnmarshaler struct {
	Name string
}

func (m *mockUnmarshaler) UnmarshalJSON(b []byte) error {
	m.Name = string(b)
	return nil
}
