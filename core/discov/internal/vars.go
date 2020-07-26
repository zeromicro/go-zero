package internal

import "time"

const (
	autoSyncInterval   = time.Minute
	coolDownInterval   = time.Second
	dialTimeout        = 5 * time.Second
	dialKeepAliveTime  = 5 * time.Second
	requestTimeout     = 3 * time.Second
	Delimiter          = '/'
	endpointsSeparator = ","
)

var (
	DialTimeout    = dialTimeout
	RequestTimeout = requestTimeout
	NewClient      = DialClient
)
