package syncx

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExclusiveCallDo(t *testing.T) {
	g := NewSingleFlight()
	v, err := g.Do("key", func() (any, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestExclusiveCallDoErr(t *testing.T) {
	g := NewSingleFlight()
	someErr := errors.New("some error")
	v, err := g.Do("key", func() (any, error) {
		return nil, someErr
	})
	if !errors.Is(err, someErr) {
		t.Errorf("Do error = %v; want someErr", err)
	}
	if v != nil {
		t.Errorf("unexpected non-nil value %#v", v)
	}
}

func TestExclusiveCallDoDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	c := make(chan string)
	var calls int32
	fn := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, err := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
}

func TestExclusiveCallDoDiffDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	broadcast := make(chan struct{})
	var calls int32
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d"}

	var wg sync.WaitGroup
	for _, key := range tests {
		wg.Add(1)
		go func(k string) {
			<-broadcast // get all goroutines ready
			_, err := g.Do(k, func() (any, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond)
				return nil, nil
			})
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			wg.Done()
		}(key)
	}

	time.Sleep(100 * time.Millisecond) // let goroutines above block
	close(broadcast)
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 5 {
		// five letters
		t.Errorf("number of calls = %d; want 5", got)
	}
}

func TestExclusiveCallDoExDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	c := make(chan string)
	var calls int32
	fn := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	var freshes int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, fresh, err := g.DoEx("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if fresh {
				atomic.AddInt32(&freshes, 1)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
	if got := atomic.LoadInt32(&freshes); got != 1 {
		t.Errorf("freshes = %d; want 1", got)
	}
}
