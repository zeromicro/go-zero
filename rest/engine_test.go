package rest

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/router"
)

const (
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC4TJk3onpqb2RYE3wwt23J9SHLFstHGSkUYFLe+nl1dEKHbD+/
Zt95L757J3xGTrwoTc7KCTxbrgn+stn0w52BNjj/kIE2ko4lbh/v8Fl14AyVR9ms
fKtKOnhe5FCT72mdtApr+qvzcC3q9hfXwkyQU32pv7q5UimZ205iKSBmgQIDAQAB
AoGAM5mWqGIAXj5z3MkP01/4CDxuyrrGDVD5FHBno3CDgyQa4Gmpa4B0/ywj671B
aTnwKmSmiiCN2qleuQYASixes2zY5fgTzt+7KNkl9JHsy7i606eH2eCKzsUa/s6u
WD8V3w/hGCQ9zYI18ihwyXlGHIgcRz/eeRh+nWcWVJzGOPUCQQD5nr6It/1yHb1p
C6l4fC4xXF19l4KxJjGu1xv/sOpSx0pOqBDEX3Mh//FU954392rUWDXV1/I65BPt
TLphdsu3AkEAvQJ2Qay/lffFj9FaUrvXuftJZ/Ypn0FpaSiUh3Ak3obBT6UvSZS0
bcYdCJCNHDtBOsWHnIN1x+BcWAPrdU7PhwJBAIQ0dUlH2S3VXnoCOTGc44I1Hzbj
Rc65IdsuBqA3fQN2lX5vOOIog3vgaFrOArg1jBkG1wx5IMvb/EnUN2pjVqUCQCza
KLXtCInOAlPemlCHwumfeAvznmzsWNdbieOZ+SXVVIpR6KbNYwOpv7oIk3Pfm9sW
hNffWlPUKhW42Gc+DIECQQDmk20YgBXwXWRM5DRPbhisIV088N5Z58K9DtFWkZsd
OBDT3dFcgZONtlmR1MqZO0pTh30lA4qovYj3Bx7A8i36
-----END RSA PRIVATE KEY-----`
)

func TestNewEngine(t *testing.T) {
	priKeyfile, err := fs.TempFilenameWithText(priKey)
	assert.Nil(t, err)
	defer os.Remove(priKeyfile)

	yamls := []string{
		`Name: foo
Host: localhost
Port: 0
Middlewares:
  Log: false
`,
		`Name: foo
Host: localhost
Port: 0
CpuThreshold: 500
Middlewares:
  Log: false
`,
		`Name: foo
Host: localhost
Port: 0
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
			timeout: time.Minute,
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
			timeout: time.Second,
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
							KeyFile:     priKeyfile,
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

	var index int32
	for _, yaml := range yamls {
		yaml := yaml
		for _, route := range routes {
			route := route
			t.Run(fmt.Sprintf("%s-%v", yaml, route.routes), func(t *testing.T) {
				var cnf RestConf
				assert.Nil(t, conf.LoadFromYamlBytes([]byte(yaml), &cnf))
				ng := newEngine(cnf)
				if atomic.AddInt32(&index, 1)%2 == 0 {
					ng.setUnsignedCallback(func(w http.ResponseWriter, r *http.Request,
						next http.Handler, strict bool, code int) {
					})
				}
				ng.addRoutes(route)
				ng.use(func(next http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					}
				})

				assert.NotNil(t, ng.start(mockedRouter{}, func(svr *http.Server) {
				}))

				timeout := time.Second * 3
				if route.timeout > timeout {
					timeout = route.timeout
				}
				assert.Equal(t, timeout, ng.timeout)
			})
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
	err := func(_ context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", http.NoBody)
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
	err := func(_ context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", http.NoBody)
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
	err := func(_ context.Context) error {
		req, err := http.NewRequest("GET", ts.URL+"/bad", http.NoBody)
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
			assert.Equal(t, time.Duration(test.timeout)*time.Millisecond*11/10, svr.WriteTimeout)
			assert.Equal(t, time.Duration(0), svr.IdleTimeout)
		})
	}
}

func TestEngine_start(t *testing.T) {
	logx.Disable()

	t.Run("http", func(t *testing.T) {
		ng := newEngine(RestConf{
			Host: "localhost",
			Port: -1,
		})
		assert.Error(t, ng.start(router.NewRouter()))
	})

	t.Run("https", func(t *testing.T) {
		ng := newEngine(RestConf{
			Host:     "localhost",
			Port:     -1,
			CertFile: "foo",
			KeyFile:  "bar",
		})
		ng.tlsConfig = &tls.Config{}
		assert.Error(t, ng.start(router.NewRouter()))
	})
}

type mockedRouter struct {
}

func (m mockedRouter) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}

func (m mockedRouter) Handle(_, _ string, _ http.Handler) error {
	return errors.New("foo")
}

func (m mockedRouter) SetNotFoundHandler(_ http.Handler) {
}

func (m mockedRouter) SetNotAllowedHandler(_ http.Handler) {
}
