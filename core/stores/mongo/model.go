package mongo

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
)

type (
	options struct {
		timeout time.Duration
	}

	Option func(opts *options)

	Model struct {
		session    *concurrentSession
		db         *mgo.Database
		collection string
		opts       []Option
	}
)

func MustNewModel(url, database, collection string, opts ...Option) *Model {
	model, err := NewModel(url, database, collection, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

func NewModel(url, database, collection string, opts ...Option) (*Model, error) {
	session, err := getConcurrentSession(url)
	if err != nil {
		return nil, err
	}

	return &Model{
		session:    session,
		db:         session.DB(database),
		collection: collection,
		opts:       opts,
	}, nil
}

func (mm *Model) Find(query interface{}) (Query, error) {
	return mm.query(func(c Collection) Query {
		return c.Find(query)
	})
}

func (mm *Model) FindId(id interface{}) (Query, error) {
	return mm.query(func(c Collection) Query {
		return c.FindId(id)
	})
}

func (mm *Model) GetCollection(session *mgo.Session) Collection {
	return newCollection(mm.db.C(mm.collection).With(session))
}

func (mm *Model) Insert(docs ...interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Insert(docs...)
	})
}

func (mm *Model) Pipe(pipeline interface{}) (Pipe, error) {
	return mm.pipe(func(c Collection) Pipe {
		return c.Pipe(pipeline)
	})
}

func (mm *Model) PutSession(session *mgo.Session) {
	mm.session.putSession(session)
}

func (mm *Model) Remove(selector interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Remove(selector)
	})
}

func (mm *Model) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c Collection) (*mgo.ChangeInfo, error) {
		return c.RemoveAll(selector)
	})
}

func (mm *Model) RemoveId(id interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.RemoveId(id)
	})
}

func (mm *Model) TakeSession() (*mgo.Session, error) {
	return mm.session.takeSession(mm.opts...)
}

func (mm *Model) Update(selector, update interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Update(selector, update)
	})
}

func (mm *Model) UpdateId(id, update interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.UpdateId(id, update)
	})
}

func (mm *Model) Upsert(selector, update interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c Collection) (*mgo.ChangeInfo, error) {
		return c.Upsert(selector, update)
	})
}

func (mm *Model) change(fn func(c Collection) (*mgo.ChangeInfo, error)) (*mgo.ChangeInfo, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) execute(fn func(c Collection) error) error {
	session, err := mm.TakeSession()
	if err != nil {
		return err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session))
}

func (mm *Model) pipe(fn func(c Collection) Pipe) (Pipe, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session)), nil
}

func (mm *Model) query(fn func(c Collection) Query) (Query, error) {
	session, err := mm.TakeSession()
	if err != nil {
		return nil, err
	}
	defer mm.PutSession(session)

	return fn(mm.GetCollection(session)), nil
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}
