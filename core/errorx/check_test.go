package errorx

import (
	"errors"
	"testing"
)

func TestIn(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	tests := []struct {
		name string
		err  error
		errs []error
		want bool
	}{
		{
			name: "Error matches one of the errors in the list",
			err:  err1,
			errs: []error{err1, err2},
			want: true,
		},
		{
			name: "Error does not match any errors in the list",
			err:  err3,
			errs: []error{err1, err2},
			want: false,
		},
		{
			name: "Empty error list",
			err:  err1,
			errs: []error{},
			want: false,
		},
		{
			name: "Nil error with non-nil list",
			err:  nil,
			errs: []error{err1, err2},
			want: false,
		},
		{
			name: "Non-nil error with nil in list",
			err:  err1,
			errs: []error{nil, err2},
			want: false,
		},
		{
			name: "Error matches nil error in the list",
			err:  nil,
			errs: []error{nil, err2},
			want: true,
		},
		{
			name: "Nil error with empty list",
			err:  nil,
			errs: []error{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := In(tt.err, tt.errs...); got != tt.want {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}
