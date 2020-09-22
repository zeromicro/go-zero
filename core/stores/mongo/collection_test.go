package mongo

import (
	"errors"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
)

func TestKeepPromise_accept(t *testing.T) {
	p := new(mockPromise)
	kp := keepablePromise{
		promise: p,
		log:     func(error) {},
	}
	assert.Nil(t, kp.accept(nil))
	assert.Equal(t, mgo.ErrNotFound, kp.accept(mgo.ErrNotFound))
}

func TestKeepPromise_keep(t *testing.T) {
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
			p := new(mockPromise)
			kp := keepablePromise{
				promise: p,
				log:     func(error) {},
			}
			assert.Equal(t, test.err, kp.keep(test.err))
			assert.Equal(t, test.accepted, p.accepted)
			assert.Equal(t, test.reason, p.reason)
		})
	}
}

type mockPromise struct {
	accepted bool
	reason   string
}

func (p *mockPromise) Accept() {
	p.accepted = true
}

func (p *mockPromise) Reject(reason string) {
	p.reason = reason
}
