package dq

import "time"

const (
	PriHigh   = 1
	PriNormal = 2
	PriLow    = 3

	defaultTimeToRun = time.Second * 5
	reserveTimeout   = time.Second * 5

	idSep   = ","
	timeSep = '/'
)
