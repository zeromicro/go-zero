package ngin

import (
	"log"
	"net/http"

	"zero/core/logx"
	"zero/ngin/handler"
	"zero/ngin/internal/router"
)

type (
	runOptions struct {
		start func(*server) error
	}

	RunOption func(*Engine)

	Engine struct {
		srv  *server
		opts runOptions
	}
)

func MustNewEngine(c NgConf, opts ...RunOption) *Engine {
	engine, err := NewEngine(c, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return engine
}

func NewEngine(c NgConf, opts ...RunOption) (*Engine, error) {
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	engine := &Engine{
		srv: newServer(c),
		opts: runOptions{
			start: func(srv *server) error {
				return srv.Start()
			},
		},
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine, nil
}

func (e *Engine) AddRoutes(rs []Route, opts ...RouteOption) {
	r := featuredRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	e.srv.AddRoutes(r)
}

func (e *Engine) AddRoute(r Route, opts ...RouteOption) {
	e.AddRoutes([]Route{r}, opts...)
}

func (e *Engine) Start() {
	handleError(e.opts.start(e.srv))
}

func (e *Engine) Stop() {
	logx.Close()
}

func (e *Engine) Use(middleware Middleware) {
	e.srv.use(middleware)
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

func WithRouter(router router.Router) RunOption {
	return func(engine *Engine) {
		engine.opts.start = func(srv *server) error {
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
	return func(engine *Engine) {
		engine.srv.SetUnauthorizedCallback(callback)
	}
}

func WithUnsignedCallback(callback handler.UnsignedCallback) RunOption {
	return func(engine *Engine) {
		engine.srv.SetUnsignedCallback(callback)
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
