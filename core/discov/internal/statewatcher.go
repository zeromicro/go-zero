//go:generate mockgen -package internal -destination statewatcher_mock.go -source statewatcher.go etcdConn
package internal

import (
	"context"
	"sync"

	"google.golang.org/grpc/connectivity"
)

type (
	etcdConn interface {
		GetState() connectivity.State
		WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
	}

	stateWatcher struct {
		disconnected bool
		currentState connectivity.State
		listeners    []func()
		lock         sync.Mutex
	}
)

func newStateWatcher() *stateWatcher {
	return new(stateWatcher)
}

func (sw *stateWatcher) addListener(l func()) {
	sw.lock.Lock()
	sw.listeners = append(sw.listeners, l)
	sw.lock.Unlock()
}

func (sw *stateWatcher) watch(conn etcdConn) {
	sw.currentState = conn.GetState()
	for {
		if conn.WaitForStateChange(context.Background(), sw.currentState) {
			newState := conn.GetState()
			sw.lock.Lock()
			sw.currentState = newState

			switch newState {
			case connectivity.TransientFailure, connectivity.Shutdown:
				sw.disconnected = true
			case connectivity.Ready:
				if sw.disconnected {
					sw.disconnected = false
					for _, l := range sw.listeners {
						l()
					}
				}
			}

			sw.lock.Unlock()
		}
	}
}
