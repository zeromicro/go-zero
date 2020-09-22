package parser

import (
	"errors"
)

var (
	unSupportDDL        = errors.New("unexpected type")
	tableBodyIsNotFound = errors.New("create table spec not found")
	errPrimaryKey       = errors.New("unexpected joint primary key")
)
