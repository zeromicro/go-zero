package internal

import (
	"strconv"
	"sync"
	"testing"

	"zero/core/logx"
	"zero/core/queue"
	"zero/core/stores/redis"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func init() {
	logx.Disable()
}

func TestRedisQueue(t *testing.T) {
	const (
		total = 1000
		key   = "queue"
	)
	r, err := miniredis.Run()
	assert.Nil(t, err)

	c := RedisKeyConf{
		RedisConf: redis.RedisConf{
			Host: r.Addr(),
			Type: redis.NodeType,
		},
		Key: key,
	}

	pusher := NewPusher(c.NewRedis(), key, WithTime())
	assert.True(t, len(pusher.Name()) > 0)
	for i := 0; i < total; i++ {
		err := pusher.Push(strconv.Itoa(i))
		assert.Nil(t, err)
	}

	consumer := new(mockedConsumer)
	consumer.wait.Add(total)
	q := queue.NewQueue(func() (queue.Producer, error) {
		return c.NewProducer(TimeSensitive(5))
	}, func() (queue.Consumer, error) {
		return consumer, nil
	})
	q.SetNumProducer(1)
	q.SetNumConsumer(1)
	go func() {
		q.Start()
	}()
	consumer.wait.Wait()
	q.Stop()

	var expect int
	for i := 0; i < total; i++ {
		expect ^= i
	}
	assert.Equal(t, expect, consumer.xor)
}

type mockedConsumer struct {
	wait sync.WaitGroup
	xor  int
}

func (c *mockedConsumer) Consume(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	c.xor ^= val
	c.wait.Done()
	return nil
}

func (c *mockedConsumer) OnEvent(event interface{}) {
}
