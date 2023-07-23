package queue

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	consumers = 4
	rounds    = 100
)

func TestQueue(t *testing.T) {
	producer := newMockedProducer(rounds)
	consumer := newMockedConsumer()
	consumer.wait.Add(consumers)
	q := NewQueue(func() (Producer, error) {
		return producer, nil
	}, func() (Consumer, error) {
		return consumer, nil
	})
	q.AddListener(new(mockedListener))
	q.SetName("mockqueue")
	q.SetNumConsumer(consumers)
	q.SetNumProducer(1)
	q.pause()
	q.resume()
	go func() {
		producer.wait.Wait()
		q.Stop()
	}()
	q.Start()
	assert.Equal(t, int32(rounds), atomic.LoadInt32(&consumer.count))
}

func TestQueue_Broadcast(t *testing.T) {
	producer := newMockedProducer(math.MaxInt32)
	consumer := newMockedConsumer()
	consumer.wait.Add(consumers)
	q := NewQueue(func() (Producer, error) {
		return producer, nil
	}, func() (Consumer, error) {
		return consumer, nil
	})
	q.AddListener(new(mockedListener))
	q.SetName("mockqueue")
	q.SetNumConsumer(consumers)
	q.SetNumProducer(1)
	go func() {
		time.Sleep(time.Millisecond * 100)
		q.Stop()
	}()
	go q.Start()
	time.Sleep(time.Millisecond * 50)
	q.Broadcast("message")
	consumer.wait.Wait()
	assert.Equal(t, int32(consumers), atomic.LoadInt32(&consumer.events))
}

func TestQueue_PauseResume(t *testing.T) {
	producer := newMockedProducer(rounds)
	consumer := newMockedConsumer()
	consumer.wait.Add(consumers)
	q := NewQueue(func() (Producer, error) {
		return producer, nil
	}, func() (Consumer, error) {
		return consumer, nil
	})
	q.AddListener(new(mockedListener))
	q.SetName("mockqueue")
	q.SetNumConsumer(consumers)
	q.SetNumProducer(1)
	go func() {
		producer.wait.Wait()
		q.Stop()
	}()
	q.Start()
	producer.listener.OnProducerPause()
	assert.Equal(t, int32(0), atomic.LoadInt32(&q.active))
	producer.listener.OnProducerResume()
	assert.Equal(t, int32(1), atomic.LoadInt32(&q.active))
	assert.Equal(t, int32(rounds), atomic.LoadInt32(&consumer.count))
}

func TestQueue_ConsumeError(t *testing.T) {
	producer := newMockedProducer(rounds)
	consumer := newMockedConsumer()
	consumer.consumeErr = errors.New("consume error")
	consumer.wait.Add(consumers)
	q := NewQueue(func() (Producer, error) {
		return producer, nil
	}, func() (Consumer, error) {
		return consumer, nil
	})
	q.AddListener(new(mockedListener))
	q.SetName("mockqueue")
	q.SetNumConsumer(consumers)
	q.SetNumProducer(1)
	go func() {
		producer.wait.Wait()
		q.Stop()
	}()
	q.Start()
	assert.Equal(t, int32(rounds), atomic.LoadInt32(&consumer.count))
}

type mockedConsumer struct {
	count      int32
	events     int32
	consumeErr error
	wait       sync.WaitGroup
}

func newMockedConsumer() *mockedConsumer {
	return new(mockedConsumer)
}

func (c *mockedConsumer) Consume(string) error {
	atomic.AddInt32(&c.count, 1)
	return c.consumeErr
}

func (c *mockedConsumer) OnEvent(any) {
	if atomic.AddInt32(&c.events, 1) <= consumers {
		c.wait.Done()
	}
}

type mockedProducer struct {
	total    int32
	count    int32
	listener ProduceListener
	wait     sync.WaitGroup
}

func newMockedProducer(total int32) *mockedProducer {
	p := new(mockedProducer)
	p.total = total
	p.wait.Add(int(total))
	return p
}

func (p *mockedProducer) AddListener(listener ProduceListener) {
	p.listener = listener
}

func (p *mockedProducer) Produce() (string, bool) {
	if atomic.AddInt32(&p.count, 1) <= p.total {
		p.wait.Done()
		return "item", true
	}

	time.Sleep(time.Second)
	return "", false
}

type mockedListener struct{}

func (l *mockedListener) OnPause() {
}

func (l *mockedListener) OnResume() {
}
