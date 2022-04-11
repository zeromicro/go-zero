package mongox

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
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
	assert.Equal(t, ErrNotFound, kp.accept(ErrNotFound))
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
			err:      ErrNotFound,
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
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	col := newCollection(mt.Coll, breaker.GetBreaker("localhost"))
	assert.Equal(t, "bar", col.(*decoratedCollection).name)
}

func TestCollectionFind(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	actual, err := c.Find(context.Background(), nil)
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	actual, err = c.Find(context.Background(), nil)
	assert.Equal(t, rejectedQuery{}, actual)
}

func TestCollectionFindOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	actual := c.FindOne(context.Background(), nil)
	assert.NotNil(t, actual)
	c.brk = new(dropBreaker)
	actual = c.FindOne(context.Background(), nil)
	assert.Equal(t, rejectedQuery{}, actual)
}

func TestCollectionInsert(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	res, err := c.InsertOne(context.Background(), bson.D{{"foo", "bar"}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	c.brk = new(dropBreaker)
	_, err = c.InsertOne(context.Background(), bson.D{{"foo", "bar"}})
	assert.NotNil(t, err)
}

func TestCollectionRemove(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
	res, err := c.DeleteOne(context.Background(), bson.D{{"foo", "bar"}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	c.brk = new(dropBreaker)
	_, err = c.DeleteOne(context.Background(), bson.D{{"foo", "bar"}})
	assert.NotNil(t, err)
}

func TestCollectionRemoveAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	res, err := c.DeleteMany(context.Background(), bson.D{{"foo", "bar"}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	c.brk = new(dropBreaker)
	_, err = c.DeleteMany(context.Background(), bson.D{{"foo", "bar"}})
	assert.NotNil(t, err)
}

func TestDecoratedCollection_FindOneAndDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	err := c.FindOneAndDelete(context.Background(), bson.D{{"foo", "bar"}})
	assert.Nil(t, err)
	c.brk = new(dropBreaker)
	err = c.FindOneAndDelete(context.Background(), bson.D{{"foo", "bar"}})
	assert.NotNil(t, err)
}

func TestCollectionUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
	resp, err := c.UpdateOne(context.Background(), bson.D{{"foo", "bar"}},
		bson.D{{"$set", bson.D{{"baz", "qux"}}}})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	c.brk = new(dropBreaker)
	_, err = c.UpdateOne(context.Background(), bson.D{{"foo", "bar"}},
		bson.D{{"$set", bson.D{{"baz", "qux"}}}})
	assert.NotNil(t, err)
}

func TestCollectionUpdateId(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	c := decoratedCollection{
		Collection: mt.Coll,
		brk:        breaker.NewBreaker(),
	}
	mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{"ok", 1}}...))
	resp, err := c.UpdateByID(context.Background(), primitive.NewObjectID(),
		bson.D{{"$set", bson.D{{"baz", "qux"}}}})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	c.brk = new(dropBreaker)
	_, err = c.UpdateByID(context.Background(), primitive.NewObjectID(),
		bson.D{{"$set", bson.D{{"baz", "qux"}}}})
	assert.NotNil(t, err)
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
