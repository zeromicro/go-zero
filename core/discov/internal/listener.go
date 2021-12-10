package internal

// Listener interface wraps the OnUpdate method.
type Listener interface {
	OnUpdate(keys, values []string, newKey string)
}
