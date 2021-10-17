package mongo

import (
	"errors"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/syncx"
)

func TestClosableIter_Close(t *testing.T) {
	errs := []error{
		nil,
		mgo.ErrNotFound,
	}

	for _, err := range errs {
		t.Run(stringx.RandId(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cleaned := syncx.NewAtomicBool()
			iter := NewMockIter(ctrl)
			iter.EXPECT().Close().Return(err)
			ci := ClosableIter{
				Iter: iter,
				Cleanup: func() {
					cleaned.Set(true)
				},
			}
			assert.Equal(t, err, ci.Close())
			assert.True(t, cleaned.True())
		})
	}
}

func TestPromisedIter_AllAndClose(t *testing.T) {
	tests := []struct {
		err      error
		accepted bool
		reason   string
	}{
		{
			err:      nil,
			accepted: true,
			reason:   "",
		},
		{
			err:      mgo.ErrNotFound,
			accepted: true,
			reason:   "",
		},
		{
			err:      errors.New("any"),
			accepted: false,
			reason:   "any",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iter := NewMockIter(ctrl)
			iter.EXPECT().All(gomock.Any()).Return(test.err)
			promise := new(mockPromise)
			pi := promisedIter{
				Iter: iter,
				promise: keepablePromise{
					promise: promise,
					log:     func(error) {},
				},
			}
			assert.Equal(t, test.err, pi.All(nil))
			assert.Equal(t, test.accepted, promise.accepted)
			assert.Equal(t, test.reason, promise.reason)
		})
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iter := NewMockIter(ctrl)
			iter.EXPECT().Close().Return(test.err)
			promise := new(mockPromise)
			pi := promisedIter{
				Iter: iter,
				promise: keepablePromise{
					promise: promise,
					log:     func(error) {},
				},
			}
			assert.Equal(t, test.err, pi.Close())
			assert.Equal(t, test.accepted, promise.accepted)
			assert.Equal(t, test.reason, promise.reason)
		})
	}
}

func TestPromisedIter_Err(t *testing.T) {
	errs := []error{
		nil,
		mgo.ErrNotFound,
	}

	for _, err := range errs {
		t.Run(stringx.RandId(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iter := NewMockIter(ctrl)
			iter.EXPECT().Err().Return(err)
			promise := new(mockPromise)
			pi := promisedIter{
				Iter: iter,
				promise: keepablePromise{
					promise: promise,
					log:     func(error) {},
				},
			}
			assert.Equal(t, err, pi.Err())
		})
	}
}

func TestPromisedIter_For(t *testing.T) {
	tests := []struct {
		err      error
		accepted bool
		reason   string
	}{
		{
			err:      nil,
			accepted: true,
			reason:   "",
		},
		{
			err:      mgo.ErrNotFound,
			accepted: true,
			reason:   "",
		},
		{
			err:      errors.New("any"),
			accepted: false,
			reason:   "any",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iter := NewMockIter(ctrl)
			iter.EXPECT().For(gomock.Any(), gomock.Any()).Return(test.err)
			promise := new(mockPromise)
			pi := promisedIter{
				Iter: iter,
				promise: keepablePromise{
					promise: promise,
					log:     func(error) {},
				},
			}
			assert.Equal(t, test.err, pi.For(nil, nil))
			assert.Equal(t, test.accepted, promise.accepted)
			assert.Equal(t, test.reason, promise.reason)
		})
	}
}

func TestRejectedIter_All(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedIter).All(nil))
}

func TestRejectedIter_Close(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedIter).Close())
}

func TestRejectedIter_Done(t *testing.T) {
	assert.False(t, new(rejectedIter).Done())
}

func TestRejectedIter_Err(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedIter).Err())
}

func TestRejectedIter_For(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedIter).For(nil, nil))
}

func TestRejectedIter_Next(t *testing.T) {
	assert.False(t, new(rejectedIter).Next(nil))
}

func TestRejectedIter_State(t *testing.T) {
	n, raw := new(rejectedIter).State()
	assert.Equal(t, int64(0), n)
	assert.Nil(t, raw)
}

func TestRejectedIter_Timeout(t *testing.T) {
	assert.False(t, new(rejectedIter).Timeout())
}

func TestIter_Done(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iter := NewMockIter(ctrl)
	iter.EXPECT().Done().Return(true)
	ci := ClosableIter{
		Iter:    iter,
		Cleanup: nil,
	}
	assert.True(t, ci.Done())
}

func TestIter_Next(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iter := NewMockIter(ctrl)
	iter.EXPECT().Next(gomock.Any()).Return(true)
	ci := ClosableIter{
		Iter:    iter,
		Cleanup: nil,
	}
	assert.True(t, ci.Next(nil))
}

func TestIter_State(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iter := NewMockIter(ctrl)
	iter.EXPECT().State().Return(int64(1), nil)
	ci := ClosableIter{
		Iter:    iter,
		Cleanup: nil,
	}
	n, raw := ci.State()
	assert.Equal(t, int64(1), n)
	assert.Nil(t, raw)
}

func TestIter_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iter := NewMockIter(ctrl)
	iter.EXPECT().Timeout().Return(true)
	ci := ClosableIter{
		Iter:    iter,
		Cleanup: nil,
	}
	assert.True(t, ci.Timeout())
}
