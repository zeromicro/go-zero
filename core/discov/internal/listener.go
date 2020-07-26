//go:generate mockgen -package internal -destination listener_mock.go -source listener.go Listener
package internal

type Listener interface {
	OnUpdate(keys []string, values []string, newKey string)
}
