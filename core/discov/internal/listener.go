package internal

type Listener interface {
	OnUpdate(keys []string, values []string, newKey string)
}
