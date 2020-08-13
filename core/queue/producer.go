package queue

type (
	Producer interface {
		AddListener(listener ProduceListener)
		Produce() (string, bool)
	}

	ProduceListener interface {
		OnProducerPause()
		OnProducerResume()
	}

	ProducerFactory func() (Producer, error)
)
