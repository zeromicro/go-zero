package rq

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueueWithTimeout(t *testing.T) {
	consumer, err := wrapWithTimeout(WithHandle(func(string) error {
		time.Sleep(time.Minute)
		return nil
	}), 100)()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ErrTimeout, consumer.Consume("any"))
}

func TestQueueWithoutTimeout(t *testing.T) {
	consumer, err := wrapWithTimeout(WithHandle(func(string) error {
		return nil
	}), 3600000)()
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, consumer.Consume("any"))
}

func BenchmarkQueue(b *testing.B) {
	b.ReportAllocs()

	consumer, err := WithHandle(func(string) error {
		return nil
	})()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		consumer.Consume(strconv.Itoa(i))
	}
}

func BenchmarkQueueWithTimeout(b *testing.B) {
	b.ReportAllocs()

	consumer, err := wrapWithTimeout(WithHandle(func(string) error {
		return nil
	}), 1000)()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		consumer.Consume(strconv.Itoa(i))
	}
}
