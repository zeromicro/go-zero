package contextx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextCancel(t *testing.T) {
	c := context.WithValue(context.Background(), "key", "value")
	c1, cancel := context.WithCancel(c)
	o := ValueOnlyFrom(c1)
	c2, _ := context.WithCancel(o)
	contexts := []context.Context{c1, c2}

	for _, c := range contexts {
		assert.NotNil(t, c.Done())
		assert.Nil(t, c.Err())

		select {
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	cancel()
	<-c1.Done()

	assert.Nil(t, o.Err())
	assert.Equal(t, context.Canceled, c1.Err())
	assert.NotEqual(t, context.Canceled, c2.Err())
}

func TestContextDeadline(t *testing.T) {
	c, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	o := ValueOnlyFrom(c)
	select {
	case <-time.After(100 * time.Millisecond):
	case <-o.Done():
		t.Fatal("ValueOnlyContext: context should not have timed out")
	}

	c, _ = context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	o = ValueOnlyFrom(c)
	c, _ = context.WithDeadline(o, time.Now().Add(20*time.Millisecond))
	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ValueOnlyContext+Deadline: context should have timed out")
	case <-c.Done():
	}
}
