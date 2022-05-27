package logx

import "errors"

const (
	// InfoLevel logs everything
	InfoLevel uint32 = iota
	// ErrorLevel includes errors, slows, stacks
	ErrorLevel
	// SevereLevel only log severe messages
	SevereLevel
)

const (
	jsonEncodingType = iota
	plainEncodingType

	jsonEncoding     = "json"
	plainEncoding    = "plain"
	plainEncodingSep = '\t'
)

const (
	accessFilename = "access.log"
	errorFilename  = "error.log"
	severeFilename = "severe.log"
	slowFilename   = "slow.log"
	statFilename   = "stat.log"

	consoleMode = "console"
	fileMode    = "file"
	volumeMode  = "volume"

	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"

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
)

var (
	// ErrLogPathNotSet is an error that indicates the log path is not set.
	ErrLogPathNotSet = errors.New("log path must be set")
	// ErrLogServiceNameNotSet is an error that indicates that the service name is not set.
	ErrLogServiceNameNotSet = errors.New("log service name must be set")
)
