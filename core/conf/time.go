package conf

import "time"

const minDuration = 100 * time.Microsecond

// CheckedDuration returns the duration that guaranteed to be greater than 100us.
// Why we need this is because users sometimes intend to use 500 to represent 500ms.
// In config, duration less than 100us should always be missing ms etc.
func CheckedDuration(duration time.Duration) time.Duration {
	if duration > minDuration {
		return duration
	}

	return duration * time.Millisecond
}
