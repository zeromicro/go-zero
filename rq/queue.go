package rq

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"zero/core/discov"
	"zero/core/logx"
	"zero/core/queue"
	"zero/core/redisqueue"
	"zero/core/service"
	"zero/core/stores/redis"
	"zero/core/stringx"
	"zero/core/threading"
	"zero/rq/constant"
	"zero/rq/update"
)

const keyLen = 6

var (
	ErrTimeout = errors.New("timeout error")

	eventHandlerPlaceholder = dummyEventHandler(0)
)

type (
	ConsumeHandle func(string) error

	ConsumeHandler interface {
		Consume(string) error
	}

	EventHandler interface {
		OnEvent(event interface{})
	}

	QueueOption func(queue *MessageQueue)

	queueOptions struct {
		renewId int64
	}

	MessageQueue struct {
		c               RmqConf
		redisQueue      *queue.Queue
		consumerFactory queue.ConsumerFactory
		options         queueOptions
		eventLock       sync.Mutex
		lastEvent       string
	}
)

func MustNewMessageQueue(c RmqConf, factory queue.ConsumerFactory, opts ...QueueOption) queue.MessageQueue {
	q, err := NewMessageQueue(c, factory, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return q
}

func NewMessageQueue(c RmqConf, factory queue.ConsumerFactory, opts ...QueueOption) (queue.MessageQueue, error) {
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	q := &MessageQueue{
		c: c,
	}

	if len(q.c.Redis.Key) == 0 {
		if len(q.c.Name) == 0 {
			q.c.Redis.Key = stringx.Randn(keyLen)
		} else {
			q.c.Redis.Key = fmt.Sprintf("%s-%s", q.c.Name, stringx.Randn(keyLen))
		}
	}
	if q.c.Timeout > 0 {
		factory = wrapWithTimeout(factory, time.Duration(q.c.Timeout)*time.Millisecond)
	}
	factory = wrapWithServerSensitive(q, factory)
	q.consumerFactory = factory
	q.redisQueue = q.buildQueue()

	for _, opt := range opts {
		opt(q)
	}

	return q, nil
}

func (q *MessageQueue) Start() {
	serviceGroup := service.NewServiceGroup()
	serviceGroup.Add(q.redisQueue)
	q.maybeAppendRenewer(serviceGroup, q.redisQueue)
	serviceGroup.Start()
}

func (q *MessageQueue) Stop() {
	logx.Close()
}

func (q *MessageQueue) buildQueue() *queue.Queue {
	inboundStore := redis.NewRedis(q.c.Redis.Host, q.c.Redis.Type, q.c.Redis.Pass)
	producerFactory := redisqueue.NewProducerFactory(inboundStore, q.c.Redis.Key,
		redisqueue.TimeSensitive(q.c.DropBefore))
	mq := queue.NewQueue(producerFactory, q.consumerFactory)

	if len(q.c.Name) > 0 {
		mq.SetName(q.c.Name)
	}
	if q.c.NumConsumers > 0 {
		mq.SetNumConsumer(q.c.NumConsumers)
	}
	if q.c.NumProducers > 0 {
		mq.SetNumProducer(q.c.NumProducers)
	}

	return mq
}

func (q *MessageQueue) compareAndSetEvent(event string) bool {
	q.eventLock.Lock()
	defer q.eventLock.Unlock()

	if q.lastEvent == event {
		return false
	}

	q.lastEvent = event
	return true
}

func (q *MessageQueue) maybeAppendRenewer(group *service.ServiceGroup, mq *queue.Queue) {
	if len(q.c.Etcd.Hosts) > 0 || len(q.c.Etcd.Key) > 0 {
		etcdValue := MarshalQueue(q.c.Redis)
		if q.c.DropBefore > 0 {
			etcdValue = strings.Join([]string{etcdValue, constant.TimedQueueType}, constant.Delimeter)
		}
		keepAliver := discov.NewRenewer(q.c.Etcd.Hosts, q.c.Etcd.Key, etcdValue, q.options.renewId)
		mq.AddListener(pauseResumeHandler{
			Renewer: keepAliver,
		})
		group.Add(keepAliver)
	}
}

func MarshalQueue(rds redis.RedisKeyConf) string {
	return strings.Join([]string{
		rds.Host,
		rds.Type,
		rds.Pass,
		rds.Key,
	}, constant.Delimeter)
}

func WithHandle(handle ConsumeHandle) queue.ConsumerFactory {
	return WithHandler(innerConsumerHandler{handle})
}

func WithHandler(handler ConsumeHandler, eventHandlers ...EventHandler) queue.ConsumerFactory {
	return func() (queue.Consumer, error) {
		if len(eventHandlers) < 1 {
			return eventConsumer{
				consumeHandler: handler,
				eventHandler:   eventHandlerPlaceholder,
			}, nil
		} else {
			return eventConsumer{
				consumeHandler: handler,
				eventHandler:   eventHandlers[0],
			}, nil
		}
	}
}

func WithHandlerFactory(factory func() (ConsumeHandler, error)) queue.ConsumerFactory {
	return func() (queue.Consumer, error) {
		if handler, err := factory(); err != nil {
			return nil, err
		} else {
			return eventlessHandler{handler}, nil
		}
	}
}

func WithRenewId(id int64) QueueOption {
	return func(mq *MessageQueue) {
		mq.options.renewId = id
	}
}

func wrapWithServerSensitive(mq *MessageQueue, factory queue.ConsumerFactory) queue.ConsumerFactory {
	return func() (queue.Consumer, error) {
		consumer, err := factory()
		if err != nil {
			return nil, err
		}

		return &serverSensitiveConsumer{
			mq:       mq,
			consumer: consumer,
		}, nil
	}
}

func wrapWithTimeout(factory queue.ConsumerFactory, dt time.Duration) queue.ConsumerFactory {
	return func() (queue.Consumer, error) {
		consumer, err := factory()
		if err != nil {
			return nil, err
		}

		return &timeoutConsumer{
			consumer: consumer,
			dt:       dt,
			timer:    time.NewTimer(dt),
		}, nil
	}
}

type innerConsumerHandler struct {
	handle ConsumeHandle
}

func (h innerConsumerHandler) Consume(v string) error {
	return h.handle(v)
}

type serverSensitiveConsumer struct {
	mq       *MessageQueue
	consumer queue.Consumer
}

func (c *serverSensitiveConsumer) Consume(msg string) error {
	if update.IsServerChange(msg) {
		change, err := update.UnmarshalServerChange(msg)
		if err != nil {
			return err
		}

		code := change.GetCode()
		if !c.mq.compareAndSetEvent(code) {
			return nil
		}

		oldHash := change.CreatePrevHash()
		newHash := change.CreateCurrentHash()
		hashChange := NewHashChange(oldHash, newHash)
		c.mq.redisQueue.Broadcast(hashChange)

		return nil
	}

	return c.consumer.Consume(msg)
}

func (c *serverSensitiveConsumer) OnEvent(event interface{}) {
	c.consumer.OnEvent(event)
}

type timeoutConsumer struct {
	consumer queue.Consumer
	dt       time.Duration
	timer    *time.Timer
}

func (c *timeoutConsumer) Consume(msg string) error {
	done := make(chan error)
	threading.GoSafe(func() {
		if err := c.consumer.Consume(msg); err != nil {
			done <- err
		}
		close(done)
	})

	c.timer.Reset(c.dt)
	select {
	case err, ok := <-done:
		c.timer.Stop()
		if ok {
			return err
		} else {
			return nil
		}
	case <-c.timer.C:
		return ErrTimeout
	}
}

func (c *timeoutConsumer) OnEvent(event interface{}) {
	c.consumer.OnEvent(event)
}

type pauseResumeHandler struct {
	discov.Renewer
}

func (pr pauseResumeHandler) OnPause() {
	pr.Pause()
}

func (pr pauseResumeHandler) OnResume() {
	pr.Resume()
}

type eventConsumer struct {
	consumeHandler ConsumeHandler
	eventHandler   EventHandler
}

func (ec eventConsumer) Consume(msg string) error {
	return ec.consumeHandler.Consume(msg)
}

func (ec eventConsumer) OnEvent(event interface{}) {
	ec.eventHandler.OnEvent(event)
}

type eventlessHandler struct {
	handler ConsumeHandler
}

func (h eventlessHandler) Consume(msg string) error {
	return h.handler.Consume(msg)
}

func (h eventlessHandler) OnEvent(event interface{}) {
}

type dummyEventHandler int

func (eh dummyEventHandler) OnEvent(event interface{}) {
}
