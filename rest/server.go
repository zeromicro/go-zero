package rest

import (
	"crypto/tls"
	"errors"
	"net/http"
	"path"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/chain"
	"github.com/zeromicro/go-zero/rest/handler"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/internal"
	"github.com/zeromicro/go-zero/rest/internal/cors"
	"github.com/zeromicro/go-zero/rest/internal/fileserver"
	"github.com/zeromicro/go-zero/rest/router"
)

type (
	// RunOption defines the method to customize a Server.
	RunOption func(*Server)

	// StartOption defines the method to customize http server.
	StartOption = internal.StartOption

	// A Server is a http server.
	Server struct {
		ngin   *engine
		router httpx.Router
	}
)

// MustNewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
// The process will exit if error occurs.
func MustNewServer(c RestConf, opts ...RunOption) *Server {
	server, err := NewServer(c, opts...)
	if err != nil {
		logx.Must(err)
	}

	return server
}

// NewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
func NewServer(c RestConf, opts ...RunOption) (*Server, error) {
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	server := &Server{
		ngin:   newEngine(c),
		router: router.NewRouter(),
	}

	opts = append([]RunOption{WithNotFoundHandler(nil)}, opts...)
	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}

// AddRoutes add given routes into the Server.
func (s *Server) AddRoutes(rs []Route, opts ...RouteOption) {
	r := featuredRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	s.ngin.addRoutes(r)
}

// AddRoute adds given route into the Server.
func (s *Server) AddRoute(r Route, opts ...RouteOption) {
	s.AddRoutes([]Route{r}, opts...)
}

// PrintRoutes prints the added routes to stdout.
func (s *Server) PrintRoutes() {
	s.ngin.print()
}

// Routes returns the HTTP routers that registered in the server.
func (s *Server) Routes() []Route {
	routes := make([]Route, 0, len(s.ngin.routes))

	for _, r := range s.ngin.routes {
		routes = append(routes, r.routes...)
	}

	return routes
}

// ServeHTTP is for test purpose, allow developer to do a unit test with
// all defined router without starting an HTTP Server.
//
// For example:
//
//	server := MustNewServer(...)
//	server.addRoute(...) // router a
//	server.addRoute(...) // router b
//	server.addRoute(...) // router c
//
//	r, _ := http.NewRequest(...)
//	w := httptest.NewRecorder(...)
//	server.ServeHTTP(w, r)
//	// verify the response
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ngin.bindRoutes(s.router)
	s.router.ServeHTTP(w, r)
}

// Start starts the Server.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (s *Server) Start() {
	handleError(s.ngin.start(s.router))
}

// StartWithOpts starts the Server.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (s *Server) StartWithOpts(opts ...StartOption) {
	handleError(s.ngin.start(s.router, opts...))
}

// Stop stops the Server.
func (s *Server) Stop() {
	logx.Close()
}

// Use adds the given middleware in the Server.
func (s *Server) Use(middleware Middleware) {
	s.ngin.use(middleware)
}

// ToMiddleware converts the given handler to a Middleware.
func ToMiddleware(handler func(next http.Handler) http.Handler) Middleware {
	return func(handle http.HandlerFunc) http.HandlerFunc {
		return handler(handle).ServeHTTP
	}
}

// WithChain returns a RunOption that uses the given chain to replace the default chain.
// JWT auth middleware and the middlewares that added by svr.Use() will be appended.
func WithChain(chn chain.Chain) RunOption {
	return func(svr *Server) {
		svr.ngin.chain = chn
	}
}

// WithCors returns a func to enable CORS for given origin, or default to all origins (*).
func WithCors(origin ...string) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(nil, origin...))
		server.router = newCorsRouter(server.router, nil, origin...)
	}
}

// WithCorsHeaders returns a RunOption to enable CORS with given headers.
func WithCorsHeaders(headers ...string) RunOption {
	const allDomains = "*"

	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(nil, allDomains))
		server.router = newCorsRouter(server.router, func(header http.Header) {
			cors.AddAllowHeaders(header, headers...)
		}, allDomains)
	}
}

// WithCustomCors returns a func to enable CORS for given origin, or default to all origins (*),
// fn lets caller customizing the response.
func WithCustomCors(middlewareFn func(header http.Header), notAllowedFn func(http.ResponseWriter),
	origin ...string) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(notAllowedFn, origin...))
		server.router = newCorsRouter(server.router, middlewareFn, origin...)
	}
}

// WithFileServer returns a RunOption to serve files from given dir with given path.
func WithFileServer(path string, fs http.FileSystem) RunOption {
	return func(server *Server) {
		server.router = newFileServingRouter(server.router, path, fs)
	}
}

// WithJwt returns a func to enable jwt authentication in given route.
func WithJwt(secret string) RouteOption {
	return func(r *featuredRoutes) {
		validateSecret(secret)
		r.jwt.enabled = true
		r.jwt.secret = secret
	}
}

// WithJwtTransition returns a func to enable jwt authentication as well as jwt secret transition.
// Which means old and new jwt secrets work together for a period.
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

// WithMaxBytes returns a RouteOption to set maxBytes with the given value.
func WithMaxBytes(maxBytes int64) RouteOption {
	return func(r *featuredRoutes) {
		r.maxBytes = maxBytes
	}
}

// WithMiddlewares adds given middlewares to given routes.
func WithMiddlewares(ms []Middleware, rs ...Route) []Route {
	for i := len(ms) - 1; i >= 0; i-- {
		rs = WithMiddleware(ms[i], rs...)
	}
	return rs
}

// WithMiddleware adds given middleware to given route.
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

// WithNotFoundHandler returns a RunOption with not found handler set to given handler.
func WithNotFoundHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		notFoundHandler := server.ngin.notFoundHandler(handler)
		server.router.SetNotFoundHandler(notFoundHandler)
	}
}

// WithNotAllowedHandler returns a RunOption with not allowed handler set to given handler.
func WithNotAllowedHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(handler)
	}
}

// WithPrefix adds group as a prefix to the route paths.
func WithPrefix(group string) RouteOption {
	return func(r *featuredRoutes) {
		routes := make([]Route, 0, len(r.routes))
		for _, rt := range r.routes {
			p := path.Join(group, rt.Path)
			routes = append(routes, Route{
				Method:  rt.Method,
				Path:    p,
				Handler: rt.Handler,
			})
		}
		r.routes = routes
	}
}

// WithPriority returns a RunOption with priority.
func WithPriority() RouteOption {
	return func(r *featuredRoutes) {
		r.priority = true
	}
}

// WithRouter returns a RunOption that make server run with given router.
func WithRouter(router httpx.Router) RunOption {
	return func(server *Server) {
		server.router = router
	}
}

// WithSignature returns a RouteOption to enable signature verification.
func WithSignature(signature SignatureConf) RouteOption {
	return func(r *featuredRoutes) {
		r.signature.enabled = true
		r.signature.Strict = signature.Strict
		r.signature.Expiry = signature.Expiry
		r.signature.PrivateKeys = signature.PrivateKeys
	}
}

// WithTimeout returns a RouteOption to set timeout with given value.
func WithTimeout(timeout time.Duration) RouteOption {
	return func(r *featuredRoutes) {
		r.timeout = timeout
	}
}

// WithTLSConfig returns a RunOption that with given tls config.
func WithTLSConfig(cfg *tls.Config) RunOption {
	return func(svr *Server) {
		svr.ngin.setTlsConfig(cfg)
	}
}

// WithUnauthorizedCallback returns a RunOption that with given unauthorized callback set.
func WithUnauthorizedCallback(callback handler.UnauthorizedCallback) RunOption {
	return func(svr *Server) {
		svr.ngin.setUnauthorizedCallback(callback)
	}
}

// WithUnsignedCallback returns a RunOption that with given unsigned callback set.
func WithUnsignedCallback(callback handler.UnsignedCallback) RunOption {
	return func(svr *Server) {
		svr.ngin.setUnsignedCallback(callback)
	}
}

func handleError(err error) {
	// ErrServerClosed means the server is closed manually
	if err == nil || errors.Is(err, http.ErrServerClosed) {
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

type corsRouter struct {
	httpx.Router
	middleware Middleware
}

func newCorsRouter(router httpx.Router, headerFn func(http.Header), origins ...string) httpx.Router {
	return &corsRouter{
		Router:     router,
		middleware: cors.Middleware(headerFn, origins...),
	}
}

func (c *corsRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.middleware(c.Router.ServeHTTP)(w, r)
}

type fileServingRouter struct {
	httpx.Router
	middleware Middleware
}

func newFileServingRouter(router httpx.Router, path string, fs http.FileSystem) httpx.Router {
	return &fileServingRouter{
		Router:     router,
		middleware: fileserver.Middleware(path, fs),
	}
}

func (f *fileServingRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.middleware(f.Router.ServeHTTP)(w, r)
}
