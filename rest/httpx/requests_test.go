package httpx

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestParseForm(t *testing.T) {
	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "/a?name=hello&age=18&percent=3.4", nil)
	assert.Nil(t, err)
	assert.Nil(t, Parse(r, &v))
	assert.Equal(t, "hello", v.Name)
	assert.Equal(t, 18, v.Age)
	assert.Equal(t, 3.4, v.Percent)
}

func TestParseForm_Error(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	r := httptest.NewRequest(http.MethodGet, "/a?name=hello;", nil)
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

	r := httptest.NewRequest(http.MethodGet, "/", nil)
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

	r := httptest.NewRequest(http.MethodGet, "/", nil)
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
		r, err := http.NewRequest(http.MethodGet, test.url, nil)
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
		r.Header.Set(ContentType, header.JsonContentType)

		assert.Nil(t, Parse(r, &v))
		assert.Equal(t, "kevin", v.Name)
		assert.Equal(t, 18, v.Age)
	})

	t.Run("hasn't body", func(t *testing.T) {
		var v struct {
			Name string `json:"name,optional"`
			Age  int    `json:"age,optional"`
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		assert.Nil(t, Parse(r, &v))
		assert.Equal(t, "", v.Name)
		assert.Equal(t, 0, v.Age)
	})
}

func TestParseRequired(t *testing.T) {
	v := struct {
		Name    string  `form:"name"`
		Percent float64 `form:"percent"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "/a?name=hello", nil)
	assert.Nil(t, err)
	assert.NotNil(t, Parse(r, &v))
}

func TestParseOptions(t *testing.T) {
	v := struct {
		Position int8 `form:"pos,options=1|2"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "/a?pos=4", nil)
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
	request, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("name", "chenquan")
	request.Header.Set("percent", "1")
	request.Header.Add("addrs", "addr1")
	request.Header.Add("addrs", "addr2")
	request.Header.Add("X-Forwarded-For", "10.0.10.11")
	request.Header.Add("x-real-ip", "10.0.11.10")
	request.Header.Add("Accept", header.JsonContentType)
	err = ParseHeaders(request, &v)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "chenquan", v.Name)
	assert.Equal(t, "1", v.Percent)
	assert.Equal(t, []string{"addr1", "addr2"}, v.Addrs)
	assert.Equal(t, "10.0.10.11", v.XForwardedFor)
	assert.Equal(t, "10.0.11.10", v.XRealIP)
	assert.Equal(t, header.JsonContentType, v.Accept)
}

func TestParseHeaders_Error(t *testing.T) {
	v := struct {
		Name string `header:"name"`
		Age  int    `header:"age"`
	}{}

	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("name", "foo")
	assert.NotNil(t, Parse(r, &v))
}

func BenchmarkParseRaw(b *testing.B) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
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
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
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
