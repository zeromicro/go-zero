package model

import (
    "context"

    "github.com/globalsign/mgo/bson"
     {{if .Cache}}cachec "github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/mongoc"{{else}}"github.com/zeromicro/go-zero/core/stores/mongo"{{end}}
)

{{if .Cache}}var prefix{{.Type}}CacheKey = "cache:{{.Type}}:"{{end}}

type {{.Type}}Model interface{
	Insert(ctx context.Context,data *{{.Type}}) error
	FindOne(ctx context.Context,id string) (*{{.Type}}, error)
	Update(ctx context.Context,data *{{.Type}}) error
	Delete(ctx context.Context,id string) error
}

type default{{.Type}}Model struct {
    {{if .Cache}}*mongoc.Model{{else}}*mongo.Model{{end}}
}

func New{{.Type}}Model(url, collection string{{if .Cache}}, c cachec.CacheConf{{end}}) {{.Type}}Model {
	return &default{{.Type}}Model{
		Model: {{if .Cache}}mongoc.MustNewModel(url, collection, c){{else}}mongo.MustNewModel(url, collection){{end}},
	}
}


func (m *default{{.Type}}Model) Insert(ctx context.Context, data *{{.Type}}) error {
    if !data.ID.Valid() {
        data.ID = bson.NewObjectId()
    }

    session, err := m.TakeSession()
    if err != nil {
        return err
    }

    defer m.PutSession(session)
    return m.GetCollection(session).Insert(data)
}

func (m *default{{.Type}}Model) FindOne(ctx context.Context, id string) (*{{.Type}}, error) {
    if !bson.IsObjectIdHex(id) {
        return nil, ErrInvalidObjectId
    }

    session, err := m.TakeSession()
    if err != nil {
        return nil, err
    }

    defer m.PutSession(session)
    var data {{.Type}}
    {{if .Cache}}key := prefix{{.Type}}CacheKey + id
    err = m.GetCollection(session).FindOneId(&data, key, bson.ObjectIdHex(id))
	{{- else}}
	err = m.GetCollection(session).FindId(bson.ObjectIdHex(id)).One(&data)
	{{- end}}
    switch err {
    case nil:
        return &data,nil
    case {{if .Cache}}mongoc.ErrNotFound{{else}}mongo.ErrNotFound{{end}}:
        return nil,ErrNotFound
    default:
        return nil,err
    }
}

func (m *default{{.Type}}Model) Update(ctx context.Context, data *{{.Type}}) error {
    session, err := m.TakeSession()
    if err != nil {
        return err
    }

    defer m.PutSession(session)
	{{if .Cache}}key := prefix{{.Type}}CacheKey + data.ID.Hex()
    return m.GetCollection(session).UpdateId(data.ID, data, key)
	{{- else}}
	return m.GetCollection(session).UpdateId(data.ID, data)
	{{- end}}
}

func (m *default{{.Type}}Model) Delete(ctx context.Context, id string) error {
    session, err := m.TakeSession()
    if err != nil {
        return err
    }

    defer m.PutSession(session)
    {{if .Cache}}key := prefix{{.Type}}CacheKey + id
    return m.GetCollection(session).RemoveId(bson.ObjectIdHex(id), key)
	{{- else}}
	return m.GetCollection(session).RemoveId(bson.ObjectIdHex(id))
	{{- end}}
}
