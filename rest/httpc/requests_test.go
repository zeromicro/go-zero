package httpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/core/trace/tracetest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"github.com/zeromicro/go-zero/rest/router"
	tcodes "go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestDoRequest(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:         "go-zero-test",
		Endpoint:     "http://localhost:14268",
		OtlpHttpPath: "/v1/traces",
		Batcher:      "otlphttp",
		Sampler:      1.0,
	})
	defer ztrace.StopAgent()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	spanContext := trace.SpanContextFromContext(resp.Request.Context())
	assert.True(t, spanContext.IsValid())
}

func TestDoRequest_NotFound(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()
	req, err := http.NewRequest(http.MethodPost, svr.URL, nil)
	assert.Nil(t, err)
	req.Header.Set(header.ContentType, header.ContentTypeJson)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDoRequest_Moved(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	_, err = DoRequest(req)
	// too many redirects
	assert.NotNil(t, err)
}

func TestDo(t *testing.T) {
	me := tracetest.NewInMemoryExporter(t)
	type Data struct {
		Key    string `path:"key"`
		Value  int    `form:"value"`
		Header string `header:"X-Header"`
		Body   string `json:"body"`
	}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodPost, "/nodes/:key",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req Data
			assert.Nil(t, httpx.Parse(r, &req))
		}))
	assert.Nil(t, err)

	svr := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer svr.Close()

	data := Data{
		Key:    "foo",
		Value:  10,
		Header: "my-header",
		Body:   "my body",
	}
	resp, err := Do(context.Background(), http.MethodPost, svr.URL+"/nodes/:key", data)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, len(me.GetSpans()))
	span := me.GetSpans()[0].Snapshot()
	assert.Equal(t, sdktrace.Status{
		Code: tcodes.Unset,
	}, span.Status())
	assert.Equal(t, 0, len(span.Events()))
	assert.Equal(t, 7, len(span.Attributes()))
}

func TestDo_Ptr(t *testing.T) {
	type Data struct {
		Key    string `path:"key"`
		Value  int    `form:"value"`
		Header string `header:"X-Header"`
		Body   string `json:"body"`
	}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodPost, "/nodes/:key",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req Data
			assert.Nil(t, httpx.Parse(r, &req))
			assert.Equal(t, "foo", req.Key)
			assert.Equal(t, 10, req.Value)
			assert.Equal(t, "my-header", req.Header)
			assert.Equal(t, "my body", req.Body)
		}))
	assert.Nil(t, err)

	svr := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer svr.Close()

	data := &Data{
		Key:    "foo",
		Value:  10,
		Header: "my-header",
		Body:   "my body",
	}
	resp, err := Do(context.Background(), http.MethodPost, svr.URL+"/nodes/:key", data)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDo_BadRequest(t *testing.T) {
	_, err := Do(context.Background(), http.MethodPost, ":/nodes/:key", nil)
	assert.NotNil(t, err)

	val1 := struct {
		Value string `json:"value,options=[a,b]"`
	}{
		Value: "c",
	}
	_, err = Do(context.Background(), http.MethodPost, "/nodes/:key", val1)
	assert.NotNil(t, err)

	val2 := struct {
		Value string `path:"val"`
	}{
		Value: "",
	}
	_, err = Do(context.Background(), http.MethodPost, "/nodes/:key", val2)
	assert.NotNil(t, err)

	val3 := struct {
		Value string `path:"key"`
		Body  string `json:"body"`
	}{
		Value: "foo",
	}
	_, err = Do(context.Background(), http.MethodGet, "/nodes/:key", val3)
	assert.NotNil(t, err)

	_, err = Do(context.Background(), "\n", "rtmp://nodes", nil)
	assert.NotNil(t, err)

	val4 := struct {
		Value string `path:"val"`
	}{
		Value: "",
	}
	_, err = Do(context.Background(), http.MethodPost, "/nodes/:val", val4)
	assert.NotNil(t, err)

	val5 := struct {
		Value   string `path:"val"`
		Another int    `path:"foo"`
	}{
		Value:   "1",
		Another: 2,
	}
	_, err = Do(context.Background(), http.MethodPost, "/nodes/:val", val5)
	assert.NotNil(t, err)
}

func TestDo_Json(t *testing.T) {
	type Data struct {
		Key    string   `path:"key"`
		Value  int      `form:"value"`
		Header string   `header:"X-Header"`
		Body   chan int `json:"body"`
	}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodPost, "/nodes/:key",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req Data
			assert.Nil(t, httpx.Parse(r, &req))
		}))
	assert.Nil(t, err)

	svr := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer svr.Close()

	data := Data{
		Key:    "foo",
		Value:  10,
		Header: "my-header",
		Body:   make(chan int),
	}
	_, err = Do(context.Background(), http.MethodPost, svr.URL+"/nodes/:key", data)
	assert.NotNil(t, err)
}

func TestDo_WithClientHttpTrace(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer svr.Close()

	enter := false
	_, err := Do(httptrace.WithClientTrace(context.Background(),
		&httptrace.ClientTrace{
			GetConn: func(hostPort string) {
				assert.Equal(t, "127.0.0.1", strings.Split(hostPort, ":")[0])
				enter = true
			},
		}), http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	assert.True(t, enter)
}

func TestBuildRequestWithBody(t *testing.T) {
	testBody := struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}{
		Key:   "foo",
		Value: 10,
	}

	testcases := []struct {
		testName  string
		method    string
		url       string
		body      any
		wantedErr error
	}{
		{
			testName:  "GET Request with Body",
			method:    http.MethodGet,
			url:       "/ping",
			body:      testBody,
			wantedErr: ErrGetWithBody,
		},
		{
			testName:  "GET Request without Body",
			method:    http.MethodGet,
			url:       "/ping",
			body:      nil,
			wantedErr: nil,
		},
		{
			testName:  "HEAD Request with Body",
			method:    http.MethodHead,
			url:       "/ping",
			body:      testBody,
			wantedErr: ErrHeadWithBody,
		},
		{
			testName:  "HEAD Request without Body",
			method:    http.MethodHead,
			url:       "/ping",
			body:      nil,
			wantedErr: nil,
		},
		{
			testName:  "POST Request with Body",
			method:    http.MethodPost,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "PUT Request with Body",
			method:    http.MethodPut,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "PATCH Request with Body",
			method:    http.MethodPatch,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "DELETE Request with Body",
			method:    http.MethodDelete,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "CONNECT Request with Body",
			method:    http.MethodConnect,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "OPTIONS Request with Body",
			method:    http.MethodOptions,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
		{
			testName:  "TRACE Request with Body",
			method:    http.MethodTrace,
			url:       "/ping",
			body:      testBody,
			wantedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := buildRequest(context.Background(), tc.method, tc.url, tc.body)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}
