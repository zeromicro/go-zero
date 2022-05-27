package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestNewEngine(t *testing.T) {
	yamls := []string{
		`Name: foo
Port: 54321
`,
		`Name: foo
Port: 54321
CpuThreshold: 500
`,
		`Name: foo
Port: 54321
CpuThreshold: 500
Verbose: true
`,
	}

	routes := []featuredRoutes{
		{
			jwt:       jwtSetting{},
			signature: signatureSetting{},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority:  true,
			jwt:       jwtSetting{},
			signature: signatureSetting{},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled: true,
			},
			signature: signatureSetting{},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled:    true,
				prevSecret: "thesecret",
			},
			signature: signatureSetting{},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled: true,
			},
			signature: signatureSetting{},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled: true,
			},
			signature: signatureSetting{
				enabled: true,
			},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled: true,
			},
			signature: signatureSetting{
				enabled: true,
				SignatureConf: SignatureConf{
					Strict: true,
				},
			},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
		{
			priority: true,
			jwt: jwtSetting{
				enabled: true,
			},
			signature: signatureSetting{
				enabled: true,
				SignatureConf: SignatureConf{
					Strict: true,
					PrivateKeys: []PrivateKeyConf{
						{
							Fingerprint: "a",
							KeyFile:     "b",
						},
					},
				},
			},
			routes: []Route{{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
			}},
		},
	}

	for _, yaml := range yamls {
		for _, route := range routes {
			var cnf RestConf
			assert.Nil(t, conf.LoadFromYamlBytes([]byte(yaml), &cnf))
			ng := newEngine(cnf)
			ng.addRoutes(route)
			ng.use(func(next http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				}
			})
			assert.NotNil(t, ng.start(mockedRouter{}))
		}
	}
}

func TestEngine_checkedTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		expect  time.Duration
	}{
		{
			name:   "not set",
			expect: time.Second,
		},
		{
			name:    "less",
			timeout: time.Millisecond * 500,
			expect:  time.Millisecond * 500,
		},
		{
			name:    "equal",
			timeout: time.Second,
			expect:  time.Second,
		},
		{
			name:    "more",
			timeout: time.Millisecond * 1500,
			expect:  time.Millisecond * 1500,
		},
	}

	ng := newEngine(RestConf{
		Timeout: 1000,
	})
	for _, test := range tests {
		assert.Equal(t, test.expect, ng.checkedTimeout(test.timeout))
	}
}

func TestEngine_checkedMaxBytes(t *testing.T) {
	tests := []struct {
		name     string
		maxBytes int64
		expect   int64
	}{
		{
			name:   "not set",
			expect: 1000,
		},
		{
			name:     "less",
			maxBytes: 500,
			expect:   500,
		},
		{
			name:     "equal",
			maxBytes: 1000,
			expect:   1000,
		},
		{
			name:     "more",
			maxBytes: 1500,
			expect:   1500,
		},
	}

	ng := newEngine(RestConf{
		MaxBytes: 1000,
	})
	for _, test := range tests {
		assert.Equal(t, test.expect, ng.checkedMaxBytes(test.maxBytes))
	}
}

func TestEngine_notFoundHandler(t *testing.T) {
	logx.Disable()

	ng := newEngine(RestConf{})
	ts := httptest.NewServer(ng.notFoundHandler(nil))
	defer ts.Close()

	client := ts.Client()
	err := func(ctx context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", nil)
		assert.Nil(t, err)
		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		return res.Body.Close()
	}(context.Background())

	assert.Nil(t, err)
}

func TestEngine_notFoundHandlerNotNil(t *testing.T) {
	logx.Disable()

	ng := newEngine(RestConf{})
	var called int32
	ts := httptest.NewServer(ng.notFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
	})))
	defer ts.Close()

	client := ts.Client()
	err := func(ctx context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", nil)
		assert.Nil(t, err)
		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		return res.Body.Close()
	}(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&called))
}

func TestEngine_notFoundHandlerNotNilWriteHeader(t *testing.T) {
	logx.Disable()

	ng := newEngine(RestConf{})
	var called int32
	ts := httptest.NewServer(ng.notFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusExpectationFailed)
	})))
	defer ts.Close()

	client := ts.Client()
	err := func(ctx context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", nil)
		assert.Nil(t, err)
		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusExpectationFailed, res.StatusCode)
		return res.Body.Close()
	}(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&called))
}

func TestEngine_withTimeout(t *testing.T) {
	logx.Disable()

	tests := []struct {
		name    string
		timeout int64
	}{
		{
			name: "not set",
		},
		{
			name:    "set",
			timeout: 1000,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ng := newEngine(RestConf{Timeout: test.timeout})
			svr := &http.Server{}
			ng.withTimeout()(svr)

			assert.Equal(t, time.Duration(test.timeout)*time.Millisecond*4/5, svr.ReadTimeout)
			assert.Equal(t, time.Duration(0), svr.ReadHeaderTimeout)
			assert.Equal(t, time.Duration(test.timeout)*time.Millisecond*9/10, svr.WriteTimeout)
			assert.Equal(t, time.Duration(0), svr.IdleTimeout)
		})
	}
}

type mockedRouter struct{}

func (m mockedRouter) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}

func (m mockedRouter) Handle(_, _ string, _ http.Handler) error {
	return errors.New("foo")
}

func (m mockedRouter) SetNotFoundHandler(_ http.Handler) {
}

func (m mockedRouter) SetNotAllowedHandler(_ http.Handler) {
}
