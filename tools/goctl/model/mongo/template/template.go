package template

// Text provides the default template for model to generate
var Text = `package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type {{.Type}}Model interface {
	Insert(ctx context.Context, data *{{.Type}}) error
	FindOne(ctx context.Context, id string) (*{{.Type}}, error)
	Update(ctx context.Context, data *{{.Type}}) error
	Delete(ctx context.Context, id string) error
}

type default{{.Type}}Model struct {
	collection *mongo.Collection
}

func New{{.Type}}Model(url string) {{.Type}}Model {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	collection := client.Database("{{.Db}}").Collection("{{.Collection}}")
	return &default{{.Type}}Model{collection: collection}
}

func (m *default{{.Type}}Model) Insert(ctx context.Context, data *{{.Type}}) error {
	_, err := m.collection.InsertOne(ctx, data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (m *default{{.Type}}Model) FindOne(ctx context.Context, id string) (*{{.Type}}, error) {
	var result {{.Type}}
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *default{{.Type}}Model) Update(ctx context.Context, data *{{.Type}}) error {
	_, err := m.collection.UpdateByID(ctx, data.Id, data)
	return err
}

func (m *default{{.Type}}Model) Delete(ctx context.Context, id string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
`

// Error provides the default template for error definition in mongo code generation.
var Error = `
package model

import "errors"

var ErrNotFound = errors.New("not found")
`
