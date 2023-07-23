package devserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/internal/health"
)

var once sync.Once

// Server is inner http server, expose some useful observability information of app.
// For example health check, metrics and pprof.
type Server struct {
	config Config
	server *http.ServeMux
	routes []string
}

// NewServer returns a new inner http Server.
func NewServer(config Config) *Server {
	return &Server{
		config: config,
		server: http.NewServeMux(),
	}
}

func (s *Server) addRoutes() {
	// route path, routes list
	s.handleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(s.routes)
	})
	// health
	s.handleFunc(s.config.HealthPath, health.CreateHttpHandler())

	// metrics
	if s.config.EnableMetrics {
		// enable prometheus global switch
		prometheus.Enable()
		s.handleFunc(s.config.MetricsPath, promhttp.Handler().ServeHTTP)
	}
	// pprof
	if s.config.EnablePprof {
		s.handleFunc("/debug/pprof/", pprof.Index)
		s.handleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.handleFunc("/debug/pprof/profile", pprof.Profile)
		s.handleFunc("/debug/pprof/symbol", pprof.Symbol)
		s.handleFunc("/debug/pprof/trace", pprof.Trace)
	}
}

func (s *Server) handleFunc(pattern string, handler http.HandlerFunc) {
	s.server.HandleFunc(pattern, handler)
	s.routes = append(s.routes, pattern)
}

// StartAsync start inner http server background.
func (s *Server) StartAsync() {
	s.addRoutes()
	threading.GoSafe(func() {
		addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
		logx.Infof("Starting dev http server at %s", addr)
		if err := http.ListenAndServe(addr, s.server); err != nil {
			logx.Error(err)
		}
	})
}

// StartAgent start inner http server by config.
func StartAgent(c Config) {
	once.Do(func() {
		if c.Enabled {
			s := NewServer(c)
			s.StartAsync()
		}
	})
}
