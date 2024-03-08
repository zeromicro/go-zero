package queue

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/rescue"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/timex"
)

const queueName = "queue"

type (
	// A Queue is a message queue.
	Queue struct {
		name                 string
		metrics              *stat.Metrics
		producerFactory      ProducerFactory
		producerRoutineGroup *threading.RoutineGroup
		consumerFactory      ConsumerFactory
		consumerRoutineGroup *threading.RoutineGroup
		producerCount        int
		consumerCount        int
		active               int32
		channel              chan string
		quit                 chan struct{}
		listeners            []Listener
		eventLock            sync.Mutex
		eventChannels        []chan any
	}

	// A Listener interface represents a listener that can be notified with queue events.
	Listener interface {
		OnPause()
		OnResume()
	}

	// A Poller interface wraps the method Poll.
	Poller interface {
		Name() string
		Poll() string
	}

	// A Pusher interface wraps the method Push.
	Pusher interface {
		Name() string
		Push(string) error
	}
)

// NewQueue returns a Queue.
func NewQueue(producerFactory ProducerFactory, consumerFactory ConsumerFactory) *Queue {
	q := &Queue{
		metrics:              stat.NewMetrics(queueName),
		producerFactory:      producerFactory,
		producerRoutineGroup: threading.NewRoutineGroup(),
		consumerFactory:      consumerFactory,
		consumerRoutineGroup: threading.NewRoutineGroup(),
		producerCount:        runtime.NumCPU(),
		consumerCount:        runtime.NumCPU() << 1,
		channel:              make(chan string),
		quit:                 make(chan struct{}),
	}
	q.SetName(queueName)

	return q
}

// AddListener adds a listener to q.
func (q *Queue) AddListener(listener Listener) {
	q.listeners = append(q.listeners, listener)
}

// Broadcast broadcasts the message to all event channels.
func (q *Queue) Broadcast(message any) {
	go func() {
		q.eventLock.Lock()
		defer q.eventLock.Unlock()

		for _, channel := range q.eventChannels {
			channel <- message
		}
	}()
}

// SetName sets the name of q.
func (q *Queue) SetName(name string) {
	q.name = name
	q.metrics.SetName(name)
}

// SetNumConsumer sets the number of consumers.
func (q *Queue) SetNumConsumer(count int) {
	q.consumerCount = count
}

// SetNumProducer sets the number of producers.
func (q *Queue) SetNumProducer(count int) {
	q.producerCount = count
}

// Start starts q.
func (q *Queue) Start() {
	q.startProducers(q.producerCount)
	q.startConsumers(q.consumerCount)

	q.producerRoutineGroup.Wait()
	close(q.channel)
	q.consumerRoutineGroup.Wait()
}

// Stop stops q.
func (q *Queue) Stop() {
	close(q.quit)
}

func (q *Queue) consume(eventChan chan any) {
	var consumer Consumer

	for {
		var err error
		if consumer, err = q.consumerFactory(); err != nil {
			logx.Errorf("Error on creating consumer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	for {
		select {
		case message, ok := <-q.channel:
			if ok {
				q.consumeOne(consumer, message)
			} else {
				logx.Info("Task channel was closed, quitting consumer...")
				return
			}
		case event := <-eventChan:
			consumer.OnEvent(event)
		}
	}
}

func (q *Queue) consumeOne(consumer Consumer, message string) {
	threading.RunSafe(func() {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			q.metrics.Add(stat.Task{
				Duration: duration,
			})
			logx.WithDuration(duration).Infof("%s", message)
		}()

		if err := consumer.Consume(message); err != nil {
			logx.Errorf("Error occurred while consuming %v: %v", message, err)
		}
	})
}

func (q *Queue) pause() {
	for _, listener := range q.listeners {
		listener.OnPause()
	}
}

func (q *Queue) produce() {
	var producer Producer

	for {
		var err error
		if producer, err = q.producerFactory(); err != nil {
			logx.Errorf("Error on creating producer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	atomic.AddInt32(&q.active, 1)
	producer.AddListener(routineListener{
		queue: q,
	})

	for {
		select {
		case <-q.quit:
			logx.Info("Quitting producer")
			return
		default:
			if v, ok := q.produceOne(producer); ok {
				q.channel <- v
			}
		}
	}
}

func (q *Queue) produceOne(producer Producer) (string, bool) {
	// avoid panic quit the producer, log it and continue
	defer rescue.Recover()

	return producer.Produce()
}

func (q *Queue) resume() {
	for _, listener := range q.listeners {
		listener.OnResume()
	}
}

func (q *Queue) startConsumers(number int) {
	for i := 0; i < number; i++ {
		eventChan := make(chan any)
		q.eventLock.Lock()
		q.eventChannels = append(q.eventChannels, eventChan)
		q.eventLock.Unlock()
		q.consumerRoutineGroup.Run(func() {
			q.consume(eventChan)
		})
	}
}

func (q *Queue) startProducers(number int) {
	for i := 0; i < number; i++ {
		q.producerRoutineGroup.Run(func() {
			q.produce()
		})
	}
}

type routineListener struct {
	queue *Queue
}

func (rl routineListener) OnProducerPause() {
	if atomic.AddInt32(&rl.queue.active, -1) <= 0 {
		rl.queue.pause()
	}
}

func (rl routineListener) OnProducerResume() {
	if atomic.AddInt32(&rl.queue.active, 1) == 1 {
		rl.queue.resume()
	}
}
