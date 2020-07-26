package internal

import "testing"

func BenchmarkRefreshCpu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RefreshCpu()
	}
}
