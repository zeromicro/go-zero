//go:build linux
// +build linux

package rest

func (s *Server) startBallast(size int) {
	s.ballast = make([]byte, size)
}
