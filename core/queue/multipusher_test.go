package queue

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiQueuePusher(t *testing.T) {
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

	pusher := NewMultiPusher(pushers)
	assert.True(t, len(pusher.Name()) > 0)

	for i := 0; i < 1000; i++ {
		_ = pusher.Push("item")
	}

	var counts []int
	for _, p := range mockedPushers {
		counts = append(counts, p.count)
	}
	mean := calcMean(counts)
	variance := calcVariance(mean, counts)
	assert.True(t, math.Abs(mean-1000*(1-failProba)) < 10)
	assert.True(t, variance < 100, fmt.Sprintf("too big variance - %.2f", variance))
}
