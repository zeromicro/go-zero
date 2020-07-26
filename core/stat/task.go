package stat

import "time"

type Task struct {
	Drop        bool
	Duration    time.Duration
	Description string
}
