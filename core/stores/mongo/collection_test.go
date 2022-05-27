package mongo

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/mongo/internal"
	"github.com/zeromicro/go-zero/core/stringx"
)

var errDummy = errors.New("dummy")

func init() {
	logx.Disable()
}

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

func TestNewCollection(t *testing.T) {
	col := newCollection(&mgo.Collection{
		Database: nil,
		Name:     "foo",
		FullName: "bar",
	}, breaker.GetBreaker("localhost"))
	assert.Equal(t, "bar", col.(*decoratedCollection).name)
}

func TestCollectionFind(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var query mgo.Query
	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(&query)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	actual := c.Find(nil)
	switch v := actual.(type) {
	case promisedQuery:
		assert.Equal(t, &query, v.Query)
		assert.Equal(t, errDummy, v.promise.keep(errDummy))
	default:
		t.Fail()
	}
	c.brk = new(dropBreaker)
	actual = c.Find(nil)
	assert.Equal(t, rejectedQuery{}, actual)
}

func TestCollectionFindId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var query mgo.Query
	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().FindId(gomock.Any()).Return(&query)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	actual := c.FindId(nil)
	switch v := actual.(type) {
	case promisedQuery:
		assert.Equal(t, &query, v.Query)
		assert.Equal(t, errDummy, v.promise.keep(errDummy))
	default:
		t.Fail()
	}
	c.brk = new(dropBreaker)
	actual = c.FindId(nil)
	assert.Equal(t, rejectedQuery{}, actual)
}

func TestCollectionInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Insert(nil, nil).Return(errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	err := c.Insert(nil, nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	err = c.Insert(nil, nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionPipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var pipe mgo.Pipe
	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Pipe(gomock.Any()).Return(&pipe)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	actual := c.Pipe(nil)
	switch v := actual.(type) {
	case promisedPipe:
		assert.Equal(t, &pipe, v.Pipe)
		assert.Equal(t, errDummy, v.promise.keep(errDummy))
	default:
		t.Fail()
	}
	c.brk = new(dropBreaker)
	actual = c.Pipe(nil)
	assert.Equal(t, rejectedPipe{}, actual)
}

func TestCollectionRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Remove(gomock.Any()).Return(errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	err := c.Remove(nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	err = c.Remove(nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionRemoveAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().RemoveAll(gomock.Any()).Return(nil, errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	_, err := c.RemoveAll(nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	_, err = c.RemoveAll(nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionRemoveId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().RemoveId(gomock.Any()).Return(errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	err := c.RemoveId(nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	err = c.RemoveId(nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	err := c.Update(nil, nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	err = c.Update(nil, nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionUpdateId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().UpdateId(gomock.Any(), gomock.Any()).Return(errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	err := c.UpdateId(nil, nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	err = c.UpdateId(nil, nil)
	assert.Equal(t, errDummy, err)
}

func TestCollectionUpsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	col.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil, errDummy)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}
	_, err := c.Upsert(nil, nil)
	assert.Equal(t, errDummy, err)
	c.brk = new(dropBreaker)
	_, err = c.Upsert(nil, nil)
	assert.Equal(t, errDummy, err)
}

func Test_logDuration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := internal.NewMockMgoCollection(ctrl)
	c := decoratedCollection{
		collection: col,
		brk:        breaker.NewBreaker(),
	}

	var buf strings.Builder
	w := logx.NewWriter(&buf)
	o := logx.Reset()
	logx.SetWriter(w)

	defer func() {
		logx.Reset()
		logx.SetWriter(o)
	}()

	buf.Reset()
	c.logDuration("foo", time.Millisecond, nil, "bar")
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	c.logDuration("foo", time.Millisecond, errors.New("bar"), make(chan int))
	assert.Contains(t, buf.String(), "bar")

	buf.Reset()
	c.logDuration("foo", slowThreshold.Load()+time.Millisecond, errors.New("bar"))
	assert.Contains(t, buf.String(), "bar")
	assert.Contains(t, buf.String(), "slowcall")

	buf.Reset()
	c.logDuration("foo", slowThreshold.Load()+time.Millisecond, nil)
	assert.Contains(t, buf.String(), "foo")
	assert.Contains(t, buf.String(), "slowcall")
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

type dropBreaker struct{}

func (d *dropBreaker) Name() string {
	return "dummy"
}

func (d *dropBreaker) Allow() (breaker.Promise, error) {
	return nil, errDummy
}

func (d *dropBreaker) Do(req func() error) error {
	return nil
}

func (d *dropBreaker) DoWithAcceptable(req func() error, acceptable breaker.Acceptable) error {
	return errDummy
}

func (d *dropBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return nil
}

func (d *dropBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error,
	acceptable breaker.Acceptable) error {
	return nil
}
