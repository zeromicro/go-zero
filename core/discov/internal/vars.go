package internal

import "time"

const (
	// Delimiter is a separator that separates the etcd path.
	Delimiter = '/'

	autoSyncInterval   = time.Minute
	coolDownInterval   = time.Second
	dialTimeout        = 5 * time.Second
	dialKeepAliveTime  = 5 * time.Second
	requestTimeout     = 3 * time.Second
	endpointsSeparator = ","
)

var (
	// DialTimeout is the dial timeout.
	DialTimeout = dialTimeout
	// RequestTimeout is the request timeout.
	RequestTimeout = requestTimeout
	// NewClient is used to create etcd clients.
	NewClient = DialClient
)
