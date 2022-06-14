package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func blue(s string) string {
	if runtime.GOOS == vars.OsWindows {
		return s
	}

	return color.New(color.FgHiBlue).Sprintf("%s", s)
}

func green(s string) string {
	if runtime.GOOS == vars.OsWindows {
		return s
	}

	return color.New(color.FgHiGreen).Sprintf("%s", s)
}

func rainbow(s string) string {
	if runtime.GOOS == vars.OsWindows {
		return s
	}

	s0 := s[0]
	return color.New(color.Attribute(91+s0%6)).Sprintf("%s", s)
}

// rpadx adds padding to the right of a string.
func rpadx(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return rainbow(fmt.Sprintf(template, s))
}
