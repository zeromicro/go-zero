package client

import (
	"github.com/tal-tech/go-zero/zrpc/internal/clientinterceptors"
)

var (
	// RetryWithDisable is an alias of clientinterceptors.RetryWithDisable.
	RetryWithDisable = clientinterceptors.RetryWithDisable
	// RetryWithCodes is an alias of clientinterceptors.RetryWithCodes.
	RetryWithCodes = clientinterceptors.RetryWithCodes
	// RetryWithBackoff is an alias of clientinterceptors.RetryWithBackoff.
	RetryWithBackoff = clientinterceptors.RetryWithBackoff
	// RetryWithPerRetryTimeout is an alias of clientinterceptors.RetryWithPerRetryTimeout.
	RetryWithPerRetryTimeout = clientinterceptors.RetryWithPerRetryTimeout
	// RetryWithMax is an alias of clientinterceptors.RetryWithMax.
	RetryWithMax = clientinterceptors.RetryWithMax
)
