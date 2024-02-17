package service

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
)

var (
	number = 1
	mutex  sync.Mutex
	done   = make(chan struct{})
)

func TestServiceGroup(t *testing.T) {
	multipliers := []int{2, 3, 5, 7}
	want := 1

	group := NewServiceGroup()
	for _, multiplier := range multipliers {
		want *= multiplier
		service := newMockedService(multiplier)
		group.Add(service)
	}

	go group.Start()

	for i := 0; i < len(multipliers); i++ {
		<-done
	}

	group.Stop()
	proc.Shutdown()

	mutex.Lock()
	defer mutex.Unlock()
	assert.Equal(t, want, number)
}

func TestServiceGroup_WithStart(t *testing.T) {
	multipliers := []int{2, 3, 5, 7}
	want := 1

	var wait sync.WaitGroup
	var lock sync.Mutex
	wait.Add(len(multipliers))
	group := NewServiceGroup()
	for _, multiplier := range multipliers {
		mul := multiplier
		group.Add(WithStart(func() {
			lock.Lock()
			want *= mul
			lock.Unlock()
			wait.Done()
		}))
	}

	go group.Start()
	wait.Wait()
	group.Stop()

	lock.Lock()
	defer lock.Unlock()
	assert.Equal(t, 210, want)
}

func TestServiceGroup_WithStarter(t *testing.T) {
	multipliers := []int{2, 3, 5, 7}
	want := 1

	var wait sync.WaitGroup
	var lock sync.Mutex
	wait.Add(len(multipliers))
	group := NewServiceGroup()
	for _, multiplier := range multipliers {
		mul := multiplier
		group.Add(WithStarter(mockedStarter{
			fn: func() {
				lock.Lock()
				want *= mul
				lock.Unlock()
				wait.Done()
			},
		}))
	}

	go group.Start()
	wait.Wait()
	group.Stop()

	lock.Lock()
	defer lock.Unlock()
	assert.Equal(t, 210, want)
}

type mockedStarter struct {
	fn func()
}

func (s mockedStarter) Start() {
	s.fn()
}

type mockedService struct {
	quit       chan struct{}
	multiplier int
}

func newMockedService(multiplier int) *mockedService {
	return &mockedService{
		quit:       make(chan struct{}),
		multiplier: multiplier,
	}
}

func (s *mockedService) Start() {
	mutex.Lock()
	number *= s.multiplier
	mutex.Unlock()
	done <- struct{}{}
	<-s.quit
}

func (s *mockedService) Stop() {
	close(s.quit)
}
