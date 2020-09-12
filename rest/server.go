package rest

import (
	"log"
	"net/http"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/handler"
	"github.com/tal-tech/go-zero/rest/httpx"
)

type (
	runOptions struct {
		start func(*engine) error
	}

	RunOption func(*Server)

	Server struct {
		ngin *engine
		opts runOptions
	}
)

func MustNewServer(c RestConf, opts ...RunOption) *Server {
	engine, err := NewServer(c, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return engine
}

func NewServer(c RestConf, opts ...RunOption) (*Server, error) {
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	server := &Server{
		ngin: newEngine(c),
		opts: runOptions{
			start: func(srv *engine) error {
				return srv.Start()
			},
		},
	}

	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

func (e *Server) AddRoutes(rs []Route, opts ...RouteOption) {
	r := featuredRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	e.ngin.AddRoutes(r)
}

func (e *Server) AddRoute(r Route, opts ...RouteOption) {
	e.AddRoutes([]Route{r}, opts...)
}

func (e *Server) Start() {
	handleError(e.opts.start(e.ngin))
}

func (e *Server) Stop() {
	logx.Close()
}

func (e *Server) Use(middleware Middleware) {
	e.ngin.use(middleware)
}

func ToMiddleware(handler func(next http.Handler) http.Handler) Middleware {
	return func(handle http.HandlerFunc) http.HandlerFunc {
		return handler(handle).ServeHTTP
	}
}

func WithJwt(secret string) RouteOption {
	return func(r *featuredRoutes) {
		validateSecret(secret)
		r.jwt.enabled = true
		r.jwt.secret = secret
	}
}

func WithJwtTransition(secret, prevSecret string) RouteOption {
	return func(r *featuredRoutes) {
		// why not validate prevSecret, because prevSecret is an already used one,
		// even it not meet our requirement, we still need to allow the transition.
		validateSecret(secret)
		r.jwt.enabled = true
		r.jwt.secret = secret
		r.jwt.prevSecret = prevSecret
	}
}

func WithMiddleware(middleware Middleware, rs ...Route) []Route {
	routes := make([]Route, len(rs))

	for i := range rs {
		route := rs[i]
		routes[i] = Route{
			Method:  route.Method,
			Path:    route.Path,
			Handler: middleware(route.Handler),
		}
	}

	return routes
}

func WithPriority() RouteOption {
	return func(r *featuredRoutes) {
		r.priority = true
	}
}

func WithRouter(router httpx.Router) RunOption {
	return func(server *Server) {
		server.opts.start = func(srv *engine) error {
			return srv.StartWithRouter(router)
		}
	}
}

func WithSignature(signature SignatureConf) RouteOption {
	return func(r *featuredRoutes) {
		r.signature.enabled = true
		r.signature.Strict = signature.Strict
		r.signature.Expiry = signature.Expiry
		r.signature.PrivateKeys = signature.PrivateKeys
	}
}

func WithUnauthorizedCallback(callback handler.UnauthorizedCallback) RunOption {
	return func(engine *Server) {
		engine.ngin.SetUnauthorizedCallback(callback)
	}
}

func WithUnsignedCallback(callback handler.UnsignedCallback) RunOption {
	return func(engine *Server) {
		engine.ngin.SetUnsignedCallback(callback)
	}
}

func handleError(err error) {
	// ErrServerClosed means the server is closed manually
	if err == nil || err == http.ErrServerClosed {
		return
	}

	logx.Error(err)
	panic(err)
}

func validateSecret(secret string) {
	if len(secret) < 8 {
		panic("secret's length can't be less than 8")
	}
}
