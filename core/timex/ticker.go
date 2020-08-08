package timex

import (
	"errors"
	"time"

	"github.com/tal-tech/go-zero/core/lang"
)

type (
	Ticker interface {
		Chan() <-chan time.Time
		Stop()
	}

	FakeTicker interface {
		Ticker
		Done()
		Tick()
		Wait(d time.Duration) error
	}

	fakeTicker struct {
		c    chan time.Time
		done chan lang.PlaceholderType
	}

	realTicker struct {
		*time.Ticker
	}
)

func NewTicker(d time.Duration) Ticker {
	return &realTicker{
		Ticker: time.NewTicker(d),
	}
}

func (rt *realTicker) Chan() <-chan time.Time {
	return rt.C
}

func NewFakeTicker() FakeTicker {
	return &fakeTicker{
		c:    make(chan time.Time, 1),
		done: make(chan lang.PlaceholderType, 1),
	}
}

func (ft *fakeTicker) Chan() <-chan time.Time {
	return ft.c
}

func (ft *fakeTicker) Done() {
	ft.done <- lang.Placeholder
}

func (ft *fakeTicker) Stop() {
	close(ft.c)
}

func (ft *fakeTicker) Tick() {
	ft.c <- Time()
}

func (ft *fakeTicker) Wait(d time.Duration) error {
	select {
	case <-time.After(d):
		return errors.New("timeout")
	case <-ft.done:
		return nil
	}
}
