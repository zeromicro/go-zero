package assertx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ErrorOrigin is used to assert error and print source and error.
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

// Error is used to assert error.
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
