package syncx

import (
	"sync"

	"zero/core/lang"
)

type DoneChan struct {
	done chan lang.PlaceholderType
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan lang.PlaceholderType),
	}
}

func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan lang.PlaceholderType {
	return dc.done
}
