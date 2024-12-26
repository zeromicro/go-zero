package subscriber

// Subscriber is the interface for configcenter subscribers.
type Subscriber interface {
	// AddListener adds a listener to the subscriber.
	AddListener(listener func()) error
	// Value returns the value of the subscriber.
	Value() (string, error)
}
