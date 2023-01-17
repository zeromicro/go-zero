package cache

import (
	"testing"

	"github.com/zeromicro/go-zero/core/timex"
)

func TestCacheStat_statLoop(t *testing.T) {
	t.Run("stat loop total 0", func(t *testing.T) {
		var stat Stat
		ticker := timex.NewFakeTicker()
		go stat.statLoop(ticker)
		ticker.Tick()
		ticker.Tick()
		ticker.Stop()
	})

	t.Run("stat loop total not 0", func(t *testing.T) {
		var stat Stat
		stat.IncrementTotal()
		ticker := timex.NewFakeTicker()
		go stat.statLoop(ticker)
		ticker.Tick()
		ticker.Tick()
		ticker.Stop()
	})
}
