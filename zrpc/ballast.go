//go:build !linux
// +build !linux

package zrpc

func (rs *RpcServer) startBallast(_ int) {
}
