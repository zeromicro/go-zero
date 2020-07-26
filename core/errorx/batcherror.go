package errorx

import "bytes"

type (
	BatchError struct {
		errs errorArray
	}

	errorArray []error
)

func (be *BatchError) Add(err error) {
	if err != nil {
		be.errs = append(be.errs, err)
	}
}

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

func (be *BatchError) NotNil() bool {
	return len(be.errs) > 0
}

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
