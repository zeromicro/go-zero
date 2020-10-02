package cmdline

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/lang"
)

func TestEnterToContinue(t *testing.T) {
	r, w, err := os.Pipe()
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ow := os.Stdout
		os.Stdout = w
		fmt.Println()
		os.Stdout = ow
	}()
	go func() {
		defer wg.Done()
		or := os.Stdin
		os.Stdin = r
		EnterToContinue()
		os.Stdin = or
	}()

	wait := make(chan lang.PlaceholderType)
	go func() {
		wg.Wait()
		close(wait)
	}()

	select {
	case <-time.After(time.Second):
		t.Error("timeout")
	case <-wait:
	}
}

func TestReadLine(t *testing.T) {
	r, w, err := os.Pipe()
	assert.Nil(t, err)

	const message = "hello"
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ow := os.Stdout
		os.Stdout = w
		fmt.Println(message)
		os.Stdout = ow
	}()
	go func() {
		defer wg.Done()
		or := os.Stdin
		os.Stdin = r
		input := ReadLine("")
		assert.Equal(t, message, input)
		os.Stdin = or
	}()

	wait := make(chan lang.PlaceholderType)
	go func() {
		wg.Wait()
		close(wait)
	}()

	select {
	case <-time.After(time.Second):
		t.Error("timeout")
	case <-wait:
	}
}
