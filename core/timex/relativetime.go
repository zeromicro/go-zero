package timex

import "time"

// Use the long enough past time as start time, in case timex.Now() - lastTime equals 0.
var initTime = time.Now().AddDate(-1, -1, -1)

func Now() time.Duration {
	return time.Since(initTime)
}

func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}

func Time() time.Time {
	return initTime.Add(Now())
}
