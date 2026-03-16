package rest

import "net/http"

// Serverless is a wrapper around Server that allows it to be used in serverless environments.
type Serverless struct {
	server *Server
}

// NewServerless creates a new Serverless instance from the provided Server.
func NewServerless(server *Server) (*Serverless, error) {
	// Ensure the server is built before using it in a serverless context.
	// Why not call server.build() when serving requests,
	// is because we need to ensure fail fast behavior.
	if err := server.build(); err != nil {
		return nil, err
	}

	return &Serverless{
		server: server,
	}, nil
}

// Serve handles HTTP requests by delegating them to the underlying Server instance.
func (s *Serverless) Serve(w http.ResponseWriter, r *http.Request) {
	s.server.serve(w, r)
}
