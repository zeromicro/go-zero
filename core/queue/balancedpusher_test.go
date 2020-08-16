package queue

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalancedQueuePusher(t *testing.T) {
	const numPushers = 100
	var pushers []Pusher
	var mockedPushers []*mockedPusher
	for i := 0; i < numPushers; i++ {
		p := &mockedPusher{
			name: "pusher:" + strconv.Itoa(i),
		}
		pushers = append(pushers, p)
		mockedPushers = append(mockedPushers, p)
	}

	pusher := NewBalancedPusher(pushers)
	assert.True(t, len(pusher.Name()) > 0)

	for i := 0; i < numPushers*1000; i++ {
		assert.Nil(t, pusher.Push("item"))
	}

	var counts []int
	for _, p := range mockedPushers {
		counts = append(counts, p.count)
	}
	mean := calcMean(counts)
	variance := calcVariance(mean, counts)
	assert.True(t, variance < 100, fmt.Sprintf("too big variance - %.2f", variance))
}

func TestBalancedQueuePusher_NoAvailable(t *testing.T) {
	pusher := NewBalancedPusher(nil)
	assert.True(t, len(pusher.Name()) == 0)
	assert.Equal(t, ErrNoAvailablePusher, pusher.Push("item"))
}
