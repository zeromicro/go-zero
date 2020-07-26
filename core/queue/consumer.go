package queue

type (
	Consumer interface {
		Consume(string) error
		OnEvent(event interface{})
	}

	ConsumerFactory func() (Consumer, error)
)
