package internal

// Listener interface wraps the OnUpdate method.
type Listener interface {
	OnUpdate(keys []string, values []string, newKey string)
}
