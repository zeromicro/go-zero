package executors

import "time"

const defaultFlushInterval = time.Second

type Execute func(tasks []interface{})
