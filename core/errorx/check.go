package errorx

import "errors"

// In checks if the given err is one of errs.
func In(err error, errs ...error) bool {
	for _, each := range errs {
		if errors.Is(err, each) {
			return true
		}
	}

	return false
}
