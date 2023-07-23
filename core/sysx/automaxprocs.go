package sysx

import "go.uber.org/automaxprocs/maxprocs"

// Automatically set GOMAXPROCS to match Linux container CPU quota.
func init() {
	maxprocs.Set(maxprocs.Logger(nil))
}
