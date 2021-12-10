package rest

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/conf"
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
			assert.Nil(t, conf.LoadConfigFromYamlBytes([]byte(yaml), &cnf))
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

type mockedRouter struct{}

func (m mockedRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}

func (m mockedRouter) Handle(method, path string, handler http.Handler) error {
	return errors.New("foo")
}

func (m mockedRouter) SetNotFoundHandler(handler http.Handler) {
}

func (m mockedRouter) SetNotAllowedHandler(handler http.Handler) {
}
