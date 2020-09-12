package console

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

type (
	Console interface {
		Success(format string, a ...interface{})
		Info(format string, a ...interface{})
		Warning(format string, a ...interface{})
		Error(format string, a ...interface{})
		Fatalln(format string, a ...interface{})
		Must(err error)
	}
	colorConsole struct {
	}
	// for idea log
	ideaConsole struct {
	}
)

func NewConsole(idea bool) Console {
	if idea {
		return NewIdeaConsole()
	}
	return NewColorConsole()
}

func NewColorConsole() *colorConsole {
	return &colorConsole{}
}

func (c *colorConsole) Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
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

func (c *colorConsole) Fatalln(format string, a ...interface{}) {
	c.Error(format, a...)
	os.Exit(1)
}

func (c *colorConsole) Must(err error) {
	if err != nil {
		c.Fatalln("%+v", err)
	}
}

func NewIdeaConsole() *ideaConsole {
	return &ideaConsole{}
}

func (i *ideaConsole) Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
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

func (i *ideaConsole) Fatalln(format string, a ...interface{}) {
	i.Error(format, a...)
	os.Exit(1)
}

func (i *ideaConsole) Must(err error) {
	if err != nil {
		i.Fatalln("%+v", err)
	}
}
