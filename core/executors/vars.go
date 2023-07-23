package executors

import "time"

const defaultFlushInterval = time.Second

// Execute defines the method to execute tasks.
type Execute func(tasks []any)
