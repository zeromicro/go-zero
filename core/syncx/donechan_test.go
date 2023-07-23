package syncx

import (
	"sync"
	"testing"
)

func TestDoneChanClose(t *testing.T) {
	doneChan := NewDoneChan()

	for i := 0; i < 5; i++ {
		doneChan.Close()
	}
}

func TestDoneChanDone(t *testing.T) {
	var waitGroup sync.WaitGroup
	doneChan := NewDoneChan()

	waitGroup.Add(1)
	go func() {
		<-doneChan.Done()
		waitGroup.Done()
	}()

	for i := 0; i < 5; i++ {
		doneChan.Close()
	}

	waitGroup.Wait()
}
