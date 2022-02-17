package mongo

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/zeromicro/go-zero/core/breaker"
)

// A Model is a mongo model.
type Model struct {
	session    *concurrentSession
	db         *mgo.Database
	collection string
	brk        breaker.Breaker
	opts       []Option
}

// MustNewModel returns a Model, exits on errors.
func MustNewModel(url, collection string, opts ...Option) *Model {
	model, err := NewModel(url, collection, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

// NewModel returns a Model.
func NewModel(url, collection string, opts ...Option) (*Model, error) {
	session, err := getConcurrentSession(url)
	if err != nil {
		return nil, err
	}

	return &Model{
		session: session,
		// If name is empty, the database name provided in the dialed URL is used instead
		db:         session.DB(""),
		collection: collection,
		brk:        breaker.GetBreaker(url),
		opts:       opts,
	}, nil
}

// Find finds a record with given query.
func (mm *Model) Find(query interface{}) (Query, error) {
	return mm.query(func(c Collection) Query {
		return c.Find(query)
	})
}

// FindId finds a record with given id.
func (mm *Model) FindId(id interface{}) (Query, error) {
	return mm.query(func(c Collection) Query {
		return c.FindId(id)
	})
}

// GetCollection returns a Collection with given session.
func (mm *Model) GetCollection(session *mgo.Session) Collection {
	return newCollection(mm.db.C(mm.collection).With(session), mm.brk)
}

// Insert inserts docs into mm.
func (mm *Model) Insert(docs ...interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Insert(docs...)
	})
}

// Pipe returns a Pipe with given pipeline.
func (mm *Model) Pipe(pipeline interface{}) (Pipe, error) {
	return mm.pipe(func(c Collection) Pipe {
		return c.Pipe(pipeline)
	})
}

// PutSession returns the given session.
func (mm *Model) PutSession(session *mgo.Session) {
	mm.session.putSession(session)
}

// Remove removes the records with given selector.
func (mm *Model) Remove(selector interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Remove(selector)
	})
}

// RemoveAll removes all with given selector and returns a mgo.ChangeInfo.
func (mm *Model) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return mm.change(func(c Collection) (*mgo.ChangeInfo, error) {
		return c.RemoveAll(selector)
	})
}

// RemoveId removes a record with given id.
func (mm *Model) RemoveId(id interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.RemoveId(id)
	})
}

// TakeSession gets a session.
func (mm *Model) TakeSession() (*mgo.Session, error) {
	return mm.session.takeSession(mm.opts...)
}

// Update updates a record with given selector.
func (mm *Model) Update(selector, update interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.Update(selector, update)
	})
}

// UpdateId updates a record with given id.
func (mm *Model) UpdateId(id, update interface{}) error {
	return mm.execute(func(c Collection) error {
		return c.UpdateId(id, update)
	})
}

// Upsert upserts a record with given selector, and returns a mgo.ChangeInfo.
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

// WithTimeout customizes an operation with given timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}
