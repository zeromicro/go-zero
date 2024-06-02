package parser

import (
	"fmt"
	"strings"
)

type errorManager struct {
	errors []string
}

func newErrorManager() *errorManager {
	return &errorManager{}
}

func (e *errorManager) add(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err.Error())
}

func (e *errorManager) error() error {
	if len(e.errors) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(e.errors, "\n"))
}
