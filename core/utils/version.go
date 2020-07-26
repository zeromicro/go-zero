package utils

import (
	"strconv"
	"strings"
)

// returns -1 if the first version is lower than the second, 0 if they are equal, and 1 if the second is lower.
func CompareVersions(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	var loop int
	if len(as) > len(bs) {
		loop = len(as)
	} else {
		loop = len(bs)
	}

	for i := 0; i < loop; i++ {
		var x, y string
		if len(as) > i {
			x = as[i]
		}
		if len(bs) > i {
			y = bs[i]
		}
		xi, _ := strconv.Atoi(x)
		yi, _ := strconv.Atoi(y)
		if xi > yi {
			return 1
		} else if xi < yi {
			return -1
		}
	}

	return 0
}
