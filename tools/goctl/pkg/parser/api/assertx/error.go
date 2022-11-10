package assertx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ErrorOrigin(t *testing.T, source string, err ...error) {
	if len(err) == 0 {
		t.Fatalf("expected errors, got 0 error")
		return
	}
	for _, e := range err {
		fmt.Printf("<%s>: %v\n", source, e)
		assert.Error(t, e)
	}
}

func Error(t *testing.T, err ...error) {
	if len(err) == 0 {
		t.Fatalf("expected errors, got 0 error")
		return
	}
	for _, e := range err {
		fmt.Println(e)
		assert.Error(t, e)
	}
}
