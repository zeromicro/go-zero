package rq

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"zero/core/discov"
	"zero/core/errorx"
	"zero/core/jsonx"
	"zero/core/lang"
	"zero/core/logx"
	"zero/core/queue"
	"zero/core/stores/redis"
	"zero/core/threading"
	"zero/rq/internal"
	"zero/rq/internal/update"
)

const (
	retryTimes      = 3
	etcdRedisFields = 4
)

var ErrPusherTypeError = errors.New("not a QueuePusher instance")

type (
	KeyFn        func(string) (key, payload string, err error)
	KeysFn       func(string) (ctx context.Context, keys []string, err error)
	AssembleFn   func(context.Context, []string) (payload string, err error)
	PusherOption func(*Pusher) error

	// just push once or do it retryTimes, it's a choice.
	// because only when at least a server is alive, and
	// pushing to the server failed, we'll return with an error
	// if waken up, but the server is going down very quickly,
	// we're going to wait again. so it's safe to push once.
	pushStrategy interface {
		addListener(listener discov.Listener)
		push(string) error
	}

	batchConsistentStrategy struct {
		keysFn     KeysFn
		assembleFn AssembleFn
		subClient  *discov.BatchConsistentSubClient
	}

	consistentStrategy struct {
		keyFn     KeyFn
		subClient *discov.ConsistentSubClient
	}

	roundRobinStrategy struct {
		subClient *discov.RoundRobinSubClient
	}

	serverListener struct {
		updater *update.IncrementalUpdater
	}

	Pusher struct {
		name            string
		endpoints       []string
		key             string
		failovers       sync.Map
		strategy        pushStrategy
		serverSensitive bool
	}
)

func NewPusher(endpoints []string, key string, opts ...PusherOption) (*Pusher, error) {
	pusher := &Pusher{
		name:      getName(key),
		endpoints: endpoints,
		key:       key,
	}

	if len(opts) == 0 {
		opts = []PusherOption{WithRoundRobinStrategy()}
	}

	for _, opt := range opts {
		if err := opt(pusher); err != nil {
			return nil, err
		}
	}

	if pusher.serverSensitive {
		listener := new(serverListener)
		listener.updater = update.NewIncrementalUpdater(listener.update)
		pusher.strategy.addListener(listener)
	}

	return pusher, nil
}

func (pusher *Pusher) Name() string {
	return pusher.name
}

func (pusher *Pusher) Push(message string) error {
	return pusher.strategy.push(message)
}

func (pusher *Pusher) close(server string, conn interface{}) error {
	logx.Errorf("dropped redis node: %s", server)

	return pusher.failover(server)
}

func (pusher *Pusher) dial(server string) (interface{}, error) {
	pusher.failovers.Delete(server)

	p, err := newPusher(server)
	if err != nil {
		return nil, err
	}

	logx.Infof("new redis node: %s", server)

	return p, nil
}

func (pusher *Pusher) failover(server string) error {
	pusher.failovers.Store(server, lang.Placeholder)

	rds, key, option, err := newRedisWithKey(server)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		defer pusher.failovers.Delete(server)

		for {
			_, ok := pusher.failovers.Load(server)
			if !ok {
				logx.Infof("redis queue (%s) revived", server)
				return
			}

			message, err := rds.Lpop(key)
			if err != nil {
				logx.Error(err)
				return
			}

			if len(message) == 0 {
				logx.Infof("repush redis queue (%s) done", server)
				return
			}

			if option == internal.TimedQueueType {
				message, err = unwrapTimedMessage(message)
				if err != nil {
					logx.Errorf("invalid timedMessage: %s, error: %s", message, err.Error())
					return
				}
			}

			if err = pusher.strategy.push(message); err != nil {
				logx.Error(err)
				return
			}
		}
	})

	return nil
}

func UnmarshalPusher(server string) (queue.QueuePusher, error) {
	store, key, option, err := newRedisWithKey(server)
	if err != nil {
		return nil, err
	}

	if option == internal.TimedQueueType {
		return internal.NewPusher(store, key, internal.WithTime()), nil
	}

	return internal.NewPusher(store, key), nil
}

func WithBatchConsistentStrategy(keysFn KeysFn, assembleFn AssembleFn, opts ...discov.BalanceOption) PusherOption {
	return func(pusher *Pusher) error {
		subClient, err := discov.NewBatchConsistentSubClient(pusher.endpoints, pusher.key, pusher.dial,
			pusher.close, opts...)
		if err != nil {
			return err
		}

		pusher.strategy = batchConsistentStrategy{
			keysFn:     keysFn,
			assembleFn: assembleFn,
			subClient:  subClient,
		}

		return nil
	}
}

func WithConsistentStrategy(keyFn KeyFn, opts ...discov.BalanceOption) PusherOption {
	return func(pusher *Pusher) error {
		subClient, err := discov.NewConsistentSubClient(pusher.endpoints, pusher.key, pusher.dial, pusher.close, opts...)
		if err != nil {
			return err
		}

		pusher.strategy = consistentStrategy{
			keyFn:     keyFn,
			subClient: subClient,
		}

		return nil
	}
}

func WithRoundRobinStrategy() PusherOption {
	return func(pusher *Pusher) error {
		subClient, err := discov.NewRoundRobinSubClient(pusher.endpoints, pusher.key, pusher.dial, pusher.close)
		if err != nil {
			return err
		}

		pusher.strategy = roundRobinStrategy{
			subClient: subClient,
		}

		return nil
	}
}

func WithServerSensitive() PusherOption {
	return func(pusher *Pusher) error {
		pusher.serverSensitive = true
		return nil
	}
}

func (bcs batchConsistentStrategy) addListener(listener discov.Listener) {
	bcs.subClient.AddListener(listener)
}

func (bcs batchConsistentStrategy) balance(keys []string) map[interface{}][]string {
	// we need to make sure the servers are available, otherwise wait forever
	for {
		if mapping, ok := bcs.subClient.Next(keys); ok {
			return mapping
		} else {
			bcs.subClient.WaitForServers()
			// make sure we don't flood logs too much in extreme conditions
			time.Sleep(time.Second)
		}
	}
}

func (bcs batchConsistentStrategy) push(message string) error {
	ctx, keys, err := bcs.keysFn(message)
	if err != nil {
		return err
	}

	var batchError errorx.BatchError
	mapping := bcs.balance(keys)
	for conn, connKeys := range mapping {
		payload, err := bcs.assembleFn(ctx, connKeys)
		if err != nil {
			batchError.Add(err)
			continue
		}

		for i := 0; i < retryTimes; i++ {
			if err = bcs.pushOnce(conn, payload); err != nil {
				batchError.Add(err)
			} else {
				break
			}
		}
	}

	return batchError.Err()
}

func (bcs batchConsistentStrategy) pushOnce(server interface{}, payload string) error {
	pusher, ok := server.(queue.QueuePusher)
	if ok {
		return pusher.Push(payload)
	} else {
		return ErrPusherTypeError
	}
}

func (cs consistentStrategy) addListener(listener discov.Listener) {
	cs.subClient.AddListener(listener)
}

func (cs consistentStrategy) push(message string) error {
	var batchError errorx.BatchError

	key, payload, err := cs.keyFn(message)
	if err != nil {
		return err
	}

	for i := 0; i < retryTimes; i++ {
		if err = cs.pushOnce(key, payload); err != nil {
			batchError.Add(err)
		} else {
			return nil
		}
	}

	return batchError.Err()
}

func (cs consistentStrategy) pushOnce(key, payload string) error {
	// we need to make sure the servers are available, otherwise wait forever
	for {
		if server, ok := cs.subClient.Next(key); ok {
			pusher, ok := server.(queue.QueuePusher)
			if ok {
				return pusher.Push(payload)
			} else {
				return ErrPusherTypeError
			}
		} else {
			cs.subClient.WaitForServers()
			// make sure we don't flood logs too much in extreme conditions
			time.Sleep(time.Second)
		}
	}
}

func (rrs roundRobinStrategy) addListener(listener discov.Listener) {
	rrs.subClient.AddListener(listener)
}

func (rrs roundRobinStrategy) push(message string) error {
	var batchError errorx.BatchError

	for i := 0; i < retryTimes; i++ {
		if err := rrs.pushOnce(message); err != nil {
			batchError.Add(err)
		} else {
			return nil
		}
	}

	return batchError.Err()
}

func (rrs roundRobinStrategy) pushOnce(message string) error {
	if server, ok := rrs.subClient.Next(); ok {
		pusher, ok := server.(queue.QueuePusher)
		if ok {
			return pusher.Push(message)
		} else {
			return ErrPusherTypeError
		}
	} else {
		rrs.subClient.WaitForServers()
		return rrs.pushOnce(message)
	}
}

func getName(key string) string {
	return fmt.Sprintf("etcd:%s", key)
}

func newPusher(server string) (queue.QueuePusher, error) {
	if rds, key, option, err := newRedisWithKey(server); err != nil {
		return nil, err
	} else if option == internal.TimedQueueType {
		return internal.NewPusher(rds, key, internal.WithTime()), nil
	} else {
		return internal.NewPusher(rds, key), nil
	}
}

func newRedisWithKey(server string) (rds *redis.Redis, key, option string, err error) {
	fields := strings.Split(server, internal.Delimeter)
	if len(fields) < etcdRedisFields {
		err = fmt.Errorf("wrong redis queue: %s, should be ip:port/type/password/key/[option]", server)
		return
	}

	addr := fields[0]
	tp := fields[1]
	pass := fields[2]
	key = fields[3]

	if len(fields) > etcdRedisFields {
		option = fields[4]
	}

	rds = redis.NewRedis(addr, tp, pass)
	return
}

func (sl *serverListener) OnUpdate(keys []string, servers []string, newKey string) {
	sl.updater.Update(keys, servers, newKey)
}

func (sl *serverListener) OnReload() {
	sl.updater.Update(nil, nil, "")
}

func (sl *serverListener) update(change update.ServerChange) {
	content, err := change.Marshal()
	if err != nil {
		logx.Error(err)
	}

	if err = broadcast(change.Servers, content); err != nil {
		logx.Error(err)
	}
}

func broadcast(servers []string, message string) error {
	var be errorx.BatchError

	for _, server := range servers {
		q, err := UnmarshalPusher(server)
		if err != nil {
			be.Add(err)
		} else {
			q.Push(message)
		}
	}

	return be.Err()
}

func unwrapTimedMessage(message string) (string, error) {
	var tm internal.TimedMessage
	if err := jsonx.UnmarshalFromString(message, &tm); err != nil {
		return "", err
	}

	return tm.Payload, nil
}
