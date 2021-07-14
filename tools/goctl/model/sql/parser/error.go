package parser

import (
	"errors"
)

var (
	errUnsupportDDL      = errors.New("unexpected type")
	errTableBodyNotFound = errors.New("create table spec not found")
	errPrimaryKey        = errors.New("unexpected join primary key")
)
