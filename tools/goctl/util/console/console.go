package console

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

type (
	// Console wraps from the fmt.Sprintf,
	// by default, it implemented the colorConsole to provide the colorful output to the console
	// and the ideaConsole to output with prefix for the plugin of intellij
	Console interface {
		Success(format string, a ...any)
		Info(format string, a ...any)
		Debug(format string, a ...any)
		Warning(format string, a ...any)
		Error(format string, a ...any)
		Fatalln(format string, a ...any)
		MarkDone()
		Must(err error)
	}

	colorConsole struct {
		enable bool
	}

	// for idea log
	ideaConsole struct{}
)

// NewConsole returns an instance of Console
func NewConsole(idea bool) Console {
	if idea {
		return NewIdeaConsole()
	}
	return NewColorConsole()
}

// NewColorConsole returns an instance of colorConsole
func NewColorConsole(enable ...bool) Console {
	logEnable := true
	for _, e := range enable {
		logEnable = e
	}
	return &colorConsole{
		enable: logEnable,
	}
}

func (c *colorConsole) Info(format string, a ...any) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
}

func (c *colorConsole) Debug(format string, a ...any) {
	if !c.enable {
		return
	}
	println(color.LightCyan.Sprintf(format, a...))
}

func (c *colorConsole) Success(format string, a ...any) {
	if !c.enable {
		return
	}
	println(color.LightGreen.Sprintf(format, a...))
}

func (c *colorConsole) Warning(format string, a ...any) {
	if !c.enable {
		return
	}
	println(color.LightYellow.Sprintf(format, a...))
}

func (c *colorConsole) Error(format string, a ...any) {
	if !c.enable {
		return
	}
	println(color.LightRed.Sprintf(format, a...))
}

func (c *colorConsole) Fatalln(format string, a ...any) {
	if !c.enable {
		return
	}
	c.Error(format, a...)
	os.Exit(1)
}

func (c *colorConsole) MarkDone() {
	if !c.enable {
		return
	}
	c.Success("Done.")
}

func (c *colorConsole) Must(err error) {
	if !c.enable {
		return
	}
	if err != nil {
		c.Fatalln("%+v", err)
	}
}

// NewIdeaConsole returns an instance of ideaConsole
func NewIdeaConsole() Console {
	return &ideaConsole{}
}

func (i *ideaConsole) Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
}

func (i *ideaConsole) Debug(format string, a ...any) {
	fmt.Println(color.LightCyan.Sprintf(format, a...))
}

func (i *ideaConsole) Success(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[SUCCESS]: ", msg)
}

func (i *ideaConsole) Warning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[WARNING]: ", msg)
}

func (i *ideaConsole) Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println("[ERROR]: ", msg)
}

func (i *ideaConsole) Fatalln(format string, a ...any) {
	i.Error(format, a...)
	os.Exit(1)
}

func (i *ideaConsole) MarkDone() {
	i.Success("Done.")
}

func (i *ideaConsole) Must(err error) {
	if err != nil {
		i.Fatalln("%+v", err)
	}
}

func println(msg any) {
	fmt.Println(msg)
}

var defaultConsole = &colorConsole{enable: true}

func Success(format string, a ...any) {
	defaultConsole.Success(format, a...)
}

func Info(format string, a ...any) {
	defaultConsole.Info(format, a...)
}

func Debug(format string, a ...any) {
	defaultConsole.Debug(format, a...)
}

func Warning(format string, a ...any) {
	defaultConsole.Warning(format, a...)
}

func Error(format string, a ...any) {
	defaultConsole.Error(format, a...)
}

func Fatalln(format string, a ...any) {
	defaultConsole.Fatalln(format, a...)
}

func MarkDone() {
	defaultConsole.MarkDone()
}

func Must(err error) {
	defaultConsole.Must(err)
}
