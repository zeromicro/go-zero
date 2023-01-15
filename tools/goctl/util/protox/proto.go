package protox

import (
	"strings"
)

func FindBeginEndOfService(service, serviceName string) (begin, end int) {
	beginIndex := strings.Index(service, serviceName)
	begin, end = -1, -1
	if beginIndex > 0 {
		for i := beginIndex; i < len(service); i++ {
			if service[i] == '}' {
				end = i
				break
			}
		}

		for i := beginIndex; i >= 0; i-- {
			if service[i] == 's' {
				begin = i
				break
			}
		}
	}
	return begin, end
}
