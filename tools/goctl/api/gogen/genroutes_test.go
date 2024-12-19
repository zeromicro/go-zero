package gogen

import (
	"testing"
	"time"
)

func Test_formatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0 * time.Nanosecond"},
		{time.Nanosecond, "1 * time.Nanosecond"},
		{100 * time.Nanosecond, "100 * time.Nanosecond"},
		{500 * time.Microsecond, "500 * time.Microsecond"},
		{2 * time.Millisecond, "2 * time.Millisecond"},
		{time.Second, "1000 * time.Millisecond"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %v; want %v", test.duration, result, test.expected)
		}
	}
}
