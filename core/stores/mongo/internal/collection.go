//go:generate mockgen -package internal -destination collection_mock.go -source collection.go

package internal

import "github.com/globalsign/mgo"

// MgoCollection interface represents a mgo collection.
type MgoCollection interface {
	Find(query any) *mgo.Query
	FindId(id any) *mgo.Query
	Insert(docs ...any) error
	Pipe(pipeline any) *mgo.Pipe
	Remove(selector any) error
	RemoveAll(selector any) (*mgo.ChangeInfo, error)
	RemoveId(id any) error
	Update(selector, update any) error
	UpdateId(id, update any) error
	Upsert(selector, update any) (*mgo.ChangeInfo, error)
}
