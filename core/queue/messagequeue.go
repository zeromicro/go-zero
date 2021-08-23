package queue

// A MessageQueue interface represents a message queue.
type MessageQueue interface {
	Start()
	Stop()
}
