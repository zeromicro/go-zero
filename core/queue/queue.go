package queue

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/rescue"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/threading"
	"github.com/tal-tech/go-zero/core/timex"
)

const queueName = "queue"

type (
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
		eventChannels        []chan interface{}
	}

	Listener interface {
		OnPause()
		OnResume()
	}

	Poller interface {
		Name() string
		Poll() string
	}

	Pusher interface {
		Name() string
		Push(string) error
	}
)

func NewQueue(producerFactory ProducerFactory, consumerFactory ConsumerFactory) *Queue {
	queue := &Queue{
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
	queue.SetName(queueName)

	return queue
}

func (queue *Queue) AddListener(listener Listener) {
	queue.listeners = append(queue.listeners, listener)
}

func (queue *Queue) Broadcast(message interface{}) {
	go func() {
		queue.eventLock.Lock()
		defer queue.eventLock.Unlock()

		for _, channel := range queue.eventChannels {
			channel <- message
		}
	}()
}

func (queue *Queue) SetName(name string) {
	queue.name = name
	queue.metrics.SetName(name)
}

func (queue *Queue) SetNumConsumer(count int) {
	queue.consumerCount = count
}

func (queue *Queue) SetNumProducer(count int) {
	queue.producerCount = count
}

func (queue *Queue) Start() {
	queue.startProducers(queue.producerCount)
	queue.startConsumers(queue.consumerCount)

	queue.producerRoutineGroup.Wait()
	close(queue.channel)
	queue.consumerRoutineGroup.Wait()
}

func (queue *Queue) Stop() {
	close(queue.quit)
}

func (queue *Queue) consume(eventChan chan interface{}) {
	var consumer Consumer

	for {
		var err error
		if consumer, err = queue.consumerFactory(); err != nil {
			logx.Errorf("Error on creating consumer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	for {
		select {
		case message, ok := <-queue.channel:
			if ok {
				queue.consumeOne(consumer, message)
			} else {
				logx.Info("Task channel was closed, quitting consumer...")
				return
			}
		case event := <-eventChan:
			consumer.OnEvent(event)
		}
	}
}

func (queue *Queue) consumeOne(consumer Consumer, message string) {
	threading.RunSafe(func() {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			queue.metrics.Add(stat.Task{
				Duration: duration,
			})
			logx.WithDuration(duration).Infof("%s", message)
		}()

		if err := consumer.Consume(message); err != nil {
			logx.Errorf("Error occurred while consuming %v: %v", message, err)
		}
	})
}

func (queue *Queue) pause() {
	for _, listener := range queue.listeners {
		listener.OnPause()
	}
}

func (queue *Queue) produce() {
	var producer Producer

	for {
		var err error
		if producer, err = queue.producerFactory(); err != nil {
			logx.Errorf("Error on creating producer: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	atomic.AddInt32(&queue.active, 1)
	producer.AddListener(routineListener{
		queue: queue,
	})

	for {
		select {
		case <-queue.quit:
			logx.Info("Quitting producer")
			return
		default:
			if v, ok := queue.produceOne(producer); ok {
				queue.channel <- v
			}
		}
	}
}

func (queue *Queue) produceOne(producer Producer) (string, bool) {
	// avoid panic quit the producer, just log it and continue
	defer rescue.Recover()

	return producer.Produce()
}

func (queue *Queue) resume() {
	for _, listener := range queue.listeners {
		listener.OnResume()
	}
}

func (queue *Queue) startConsumers(number int) {
	for i := 0; i < number; i++ {
		eventChan := make(chan interface{})
		queue.eventLock.Lock()
		queue.eventChannels = append(queue.eventChannels, eventChan)
		queue.eventLock.Unlock()
		queue.consumerRoutineGroup.Run(func() {
			queue.consume(eventChan)
		})
	}
}

func (queue *Queue) startProducers(number int) {
	for i := 0; i < number; i++ {
		queue.producerRoutineGroup.Run(func() {
			queue.produce()
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
