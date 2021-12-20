//go:build linux
// +build linux

package zrpc

func (rs *RpcServer) startBallast(size int) {
	rs.ballast = make([]byte, size)
}
