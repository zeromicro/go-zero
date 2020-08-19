package console

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

type (
	Console interface {
		Success(format string, a ...interface{})
		Warning(format string, a ...interface{})
		Error(format string, a ...interface{})
	}
	colorConsole struct {
	}
	// for idea log
	ideaConsole struct {
	}
)

func NewColorConsole() *colorConsole {
	return &colorConsole{}
}

func (c *colorConsole) Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(aurora.Green(msg))
}

func (c *colorConsole) Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(aurora.Yellow(msg))
}

func (c *colorConsole) Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(aurora.Red(msg))
}

func NewIdeaConsole() *ideaConsole {
	return &ideaConsole{}
}

func (i *ideaConsole) Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[SUCCESS]: ", msg)
}

func (i *ideaConsole) Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[WARNING]: ", msg)
}

func (i *ideaConsole) Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[ERROR]: ", msg)
}
