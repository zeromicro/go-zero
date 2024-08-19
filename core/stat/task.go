package stat

import "time"

// A Task is a task reported to Metrics.
type Task struct {
	Drop        bool
	Duration    time.Duration
	Description string
}
