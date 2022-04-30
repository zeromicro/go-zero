package logx

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/timex"
)

func getCaller(callDepth int) string {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		buf.WriteString(short)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	return buf.String()
}

func getTimestamp() string {
	return timex.Time().Format(timeFormat)
}
