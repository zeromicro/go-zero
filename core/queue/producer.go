package queue

type (
	// A Producer interface represents a producer that produces messages.
	Producer interface {
		AddListener(listener ProduceListener)
		Produce() (string, bool)
	}

	// A ProduceListener interface represents a produce listener.
	ProduceListener interface {
		OnProducerPause()
		OnProducerResume()
	}

	// ProducerFactory defines the method to generate a Producer.
	ProducerFactory func() (Producer, error)
)
