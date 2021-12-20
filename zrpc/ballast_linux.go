//go:build linux
// +build linux

package zrpc

func (rs *Server) startBallast(size int) {
	rs.ballast = make([]byte, size)
}
