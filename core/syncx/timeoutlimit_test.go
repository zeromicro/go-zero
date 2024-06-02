package syncx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutLimit(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
	}{
		{
			name: "no wait",
		},
		{
			name:     "wait",
			interval: time.Millisecond * 100,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			limit := NewTimeoutLimit(2)
			assert.Nil(t, limit.Borrow(time.Millisecond*200))
			assert.Nil(t, limit.Borrow(time.Millisecond*200))
			var wait1, wait2, wait3 sync.WaitGroup
			wait1.Add(1)
			wait2.Add(1)
			wait3.Add(1)
			go func() {
				wait1.Wait()
				wait2.Done()
				time.Sleep(test.interval)
				assert.Nil(t, limit.Return())
				wait3.Done()
			}()
			wait1.Done()
			wait2.Wait()
			assert.Nil(t, limit.Borrow(time.Second))
			wait3.Wait()
			assert.Equal(t, ErrTimeout, limit.Borrow(time.Millisecond*100))
			assert.Nil(t, limit.Return())
			assert.Nil(t, limit.Return())
			assert.Equal(t, ErrLimitReturn, limit.Return())
		})
	}
}
