package logx

import (
	"context"
	"fmt"
	"io"
	"testing"
)

func benchmarkLogger(b *testing.B, numFields int, cache bool) {
	w := NewWriter(io.Discard)
	old := writer.Swap(w)
	defer writer.Store(old)

	if cache {
		EnableCache()
	} else {
		DisableCache()
	}

	fields := make([]LogField, numFields)
	for i := 0; i < numFields; i++ {
		fields[i] = Field(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	logger := WithContext(context.Background()).WithFields(fields...)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkWithCache10Fields(b *testing.B) {
	benchmarkLogger(b, 10, true)
}

func BenchmarkWithoutCache10Fields(b *testing.B) {
	benchmarkLogger(b, 10, false)
}

func BenchmarkWithCache1000Fields(b *testing.B) {
	benchmarkLogger(b, 1000, true)
}

func BenchmarkWithoutCache1000Fields(b *testing.B) {
	benchmarkLogger(b, 1000, false)
}
