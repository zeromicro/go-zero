package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	err1 = "first error"
	err2 = "second error"
)

func TestBatchErrorNil(t *testing.T) {
	var batch BatchError
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
	batch.Add(nil)
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
}

func TestBatchErrorNilFromFunc(t *testing.T) {
	err := func() error {
		var be BatchError
		return be.Err()
	}()
	assert.True(t, err == nil)
}

func TestBatchErrorOneError(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	assert.NotNil(t, batch)
	assert.Equal(t, err1, batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchErrorWithErrors(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	batch.Add(errors.New(err2))
	assert.NotNil(t, batch)
	assert.Equal(t, fmt.Sprintf("%s\n%s", err1, err2), batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchError_Error(t *testing.T) {
	type fields struct {
		be *BatchError
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"nil",
			fields{nil},
			"",
		},
		{
			"nil errors",
			fields{&BatchError{}},
			"",
		},
		{
			"one error",
			fields{&BatchError{errs: errorArray{errors.New(err1)}}},
			err1,
		},
		{
			"two errors",
			fields{&BatchError{errs: errorArray{errors.New(err1), errors.New(err2)}}},
			fmt.Sprintf("%s\n%s", err1, err2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := tt.fields.be
			assert.Equalf(t, tt.want, be.Error(), "Error()")
		})
	}
}

func TestErrors(t *testing.T) {
	e1 := errors.New(err1)
	e2 := errors.New(err2)

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			"nil",
			args{nil},
			nil,
		},
		{
			"nil errors",
			args{&BatchError{}},
			nil,
		},
		{
			"one error",
			args{e1},
			[]error{e1},
		},
		{
			"BatchError - one error",
			args{&BatchError{errs: errorArray{e1}}},
			[]error{e1},
		},
		{
			"BatchError - two errors",
			args{&BatchError{errs: errorArray{e1, e2}}},
			[]error{e1, e2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Errors(tt.args.err), "Errors(%v)", tt.args.err)
		})
	}
}

func TestBatchError_Errors(t *testing.T) {
	var be *BatchError
	assert.Nil(t, be.Errors())

	e1 := errors.New(err1)
	e2 := errors.New(err2)

	be = &BatchError{}
	be.Add(e1, e2)
	assert.Equal(t, []error{e1, e2}, be.Errors())
}
