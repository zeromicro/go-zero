package errorx

import "bytes"

type (
	// A BatchError is an error that can hold multiple errors.
	BatchError struct {
		errs errorArray
	}

	errorArray []error
)

// Errors returns a slice containing zero or more errors that the supplied
// error is composed of. If the error is nil, a nil slice is returned.
func Errors(err error) []error {
	if err == nil {
		return nil
	}

	if be, ok := err.(*BatchError); ok {
		return be.Errors()
	}

	return []error{err}
}

// Add adds errs to be, nil errors are ignored.
func (be *BatchError) Add(errs ...error) {
	for _, err := range errs {
		if err != nil {
			be.errs = append(be.errs, err)
		}
	}
}

// Err returns an error that represents all errors.
func (be *BatchError) Err() error {
	switch len(be.errs) {
	case 0:
		return nil
	case 1:
		return be.errs[0]
	default:
		return be.errs
	}
}

// NotNil checks if any error inside.
func (be *BatchError) NotNil() bool {
	return len(be.errs) > 0
}

// Errors returns the list of underlying errors.
// Callers of this function are free to modify the returned slice.
func (be *BatchError) Errors() []error {
	if be == nil {
		return nil
	}
	return append([]error(nil), be.errs...)
}

// Error implements error interface. This allows users to convert error to BatchError.
func (be *BatchError) Error() string {
	if be == nil || be.Err() == nil {
		return ""
	}
	return be.Err().Error()
}

// Error returns a string that represents inside errors.
func (ea errorArray) Error() string {
	var buf bytes.Buffer

	for i := range ea {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(ea[i].Error())
	}

	return buf.String()
}
