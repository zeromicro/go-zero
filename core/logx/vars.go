package logx

import (
	"errors"

	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	// DebugLevel logs everything
	DebugLevel uint32 = iota
	// InfoLevel does not include debugs
	InfoLevel
	// ErrorLevel includes errors, slows, stacks
	ErrorLevel
	// SevereLevel only log severe messages
	SevereLevel
)

const (
	jsonEncodingType = iota
	plainEncodingType
)

const (
	plainEncoding    = "plain"
	plainEncodingSep = '\t'
	sizeRotationRule = "size"

	accessFilename = "access.log"
	errorFilename  = "error.log"
	severeFilename = "severe.log"
	slowFilename   = "slow.log"
	statFilename   = "stat.log"

	fileMode   = "file"
	volumeMode = "volume"

	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"
	levelDebug  = "debug"

	backupFileDelimiter = "-"
	flags               = 0x0
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	durationKey  = "duration"
	levelKey     = "level"
	spanKey      = "span"
	timestampKey = "@timestamp"
	traceKey     = "trace"
	truncatedKey = "truncated"
)

var (
	// ErrLogPathNotSet is an error that indicates the log path is not set.
	ErrLogPathNotSet = errors.New("log path must be set")
	// ErrLogServiceNameNotSet is an error that indicates that the service name is not set.
	ErrLogServiceNameNotSet = errors.New("log service name must be set")
	// ExitOnFatal defines whether to exit on fatal errors, defined here to make it easier to test.
	ExitOnFatal = syncx.ForAtomicBool(true)

	truncatedField = Field(truncatedKey, true)
)
