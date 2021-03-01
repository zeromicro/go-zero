//go:generate mockgen -package internal -destination collection_mock.go -source collection.go

package internal

import "github.com/globalsign/mgo"

// MgoCollection interface represents a mgo collection.
type MgoCollection interface {
	Find(query interface{}) *mgo.Query
	FindId(id interface{}) *mgo.Query
	Insert(docs ...interface{}) error
	Pipe(pipeline interface{}) *mgo.Pipe
	Remove(selector interface{}) error
	RemoveAll(selector interface{}) (*mgo.ChangeInfo, error)
	RemoveId(id interface{}) error
	Update(selector, update interface{}) error
	UpdateId(id, update interface{}) error
	Upsert(selector, update interface{}) (*mgo.ChangeInfo, error)
}
