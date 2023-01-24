package queue

type (
	// A Consumer interface represents a consumer that can consume string messages.
	Consumer interface {
		Consume(string) error
		OnEvent(event any)
	}

	// ConsumerFactory defines the factory to generate consumers.
	ConsumerFactory func() (Consumer, error)
)
