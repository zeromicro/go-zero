package mon

import (
	"log"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	Collection
	collName string
	cli      *mongo.Client
	brk      breaker.Breaker
	opts     []Option
}

// MustNewModel returns a Model, exits on errors.
func MustNewModel(uri, db, collection string, opts ...Option) *Model {
	model, err := NewModel(uri, db, collection, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

// NewModel returns a Model.
func NewModel(uri, db, collection string, opts ...Option) (*Model, error) {
	cli, err := getClient(uri)
	if err != nil {
		return nil, err
	}

	brk := breaker.GetBreaker(uri)
	coll := newCollection(cli.Database(db).Collection(collection), brk)
	return &Model{
		Collection: coll,
		collName:   collection,
		cli:        cli,
		brk:        brk,
		opts:       opts,
	}, nil
}

func (m *Model) StartSession(opts ...*mopt.SessionOptions) (sess mongo.Session, err error) {
	err = m.brk.DoWithAcceptable(func() error {
		starTime := timex.Now()
		defer func() {
			logDuration(m.collName, "StartSession", starTime, err)
		}()

		sess, err = m.cli.StartSession(opts...)
		return err
	}, acceptable)
	return
}
