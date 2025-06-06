// Code generated by goctl. DO NOT EDIT.
// goctl {{.version}}

package model

import (
    "context"
    "time"

    {{if .Cache}}"github.com/zeromicro/go-zero/core/stores/monc"{{else}}"github.com/zeromicro/go-zero/core/stores/mon"{{end}}
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

{{if .Cache}}var prefix{{.Type}}CacheKey = "{{if .Prefix}}{{.Prefix}}:{{end}}cache:{{.lowerType}}:"{{end}}

type {{.lowerType}}Model interface{
    Insert(ctx context.Context,data *{{.Type}}) error
    FindOne(ctx context.Context,id string) (*{{.Type}}, error)
    Update(ctx context.Context,data *{{.Type}}) (*mongo.UpdateResult, error)
    Delete(ctx context.Context,id string) (int64, error)
}

type default{{.Type}}Model struct {
    conn {{if .Cache}}*monc.Model{{else}}*mon.Model{{end}}
}

func newDefault{{.Type}}Model(conn {{if .Cache}}*monc.Model{{else}}*mon.Model{{end}}) *default{{.Type}}Model {
    return &default{{.Type}}Model{conn: conn}
}


func (m *default{{.Type}}Model) Insert(ctx context.Context, data *{{.Type}}) error {
    if data.ID.IsZero() {
        data.ID = primitive.NewObjectID()
        data.CreateAt = time.Now()
        data.UpdateAt = time.Now()
    }

    {{if .Cache}}key := prefix{{.Type}}CacheKey + data.ID.Hex(){{end}}
    _, err := m.conn.InsertOne(ctx, {{if .Cache}}key, {{end}} data)
    return err
}

func (m *default{{.Type}}Model) FindOne(ctx context.Context, id string) (*{{.Type}}, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, ErrInvalidObjectId
    }

    var data {{.Type}}
    {{if .Cache}}key := prefix{{.Type}}CacheKey + id{{end}}
    err = m.conn.FindOne(ctx, {{if .Cache}}key, {{end}}&data, bson.M{"_id": oid})
    switch err {
    case nil:
        return &data, nil
    case {{if .Cache}}monc{{else}}mon{{end}}.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *default{{.Type}}Model) Update(ctx context.Context, data *{{.Type}}) (*mongo.UpdateResult, error) {
    data.UpdateAt = time.Now()
    {{if .Cache}}key := prefix{{.Type}}CacheKey + data.ID.Hex(){{end}}
    res, err := m.conn.UpdateOne(ctx, {{if .Cache}}key, {{end}}bson.M{"_id": data.ID}, bson.M{"$set": data})
    return res, err
}

func (m *default{{.Type}}Model) Delete(ctx context.Context, id string) (int64, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return 0, ErrInvalidObjectId
    }
	{{if .Cache}}key := prefix{{.Type}}CacheKey +id{{end}}
    res, err := m.conn.DeleteOne(ctx, {{if .Cache}}key, {{end}}bson.M{"_id": oid})
	return res, err
}
