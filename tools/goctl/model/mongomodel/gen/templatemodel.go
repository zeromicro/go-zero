package gen

const (
	quotationMark = "`"
	//templates that do not use caching
	noCacheTemplate = `package model

import (
	{{.importArray}}
)

var ErrNotFound = mongoc.ErrNotFound

type (
	{{.modelName}}Model struct {
		*mongoc.Model
	}

	{{.modelName}} struct {
		{{.modelFields}}
	}
)

func New{{.modelName}}Model(url, database, collection string, c cache.CacheConf, opts ...cache.Option) *{{.modelName}}Model {
	return &{{.modelName}}Model{mongoc.MustNewModel(url, database, collection, c, opts...)}
}

func (m *{{.modelName}}Model) FindOne(id string) (*{{.modelName}}, error) {
	session, err := m.Model.TakeSession()
	if err != nil {
		return nil, err
	}
	defer m.Model.PutSession(session)

	var result {{.modelName}}
	err = m.GetCollection(session).FindOneIdNoCache(&result,bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *{{.modelName}}Model) Delete(id string) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)
	return m.GetCollection(session).RemoveIdNoCache(bson.ObjectIdHex(id))
}

func (m *{{.modelName}}Model) Insert(data *{{.modelName}}) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	return m.GetCollection(session).Insert(data)
}

func (m *{{.modelName}}Model) Update(data *{{.modelName}}) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	data.UpdateTime = time.Now()
	return m.GetCollection(session).UpdateIdNoCache(data.Id, data)
}
`

	//use cache template
	cacheTemplate = `package model

import (
	{{.importArray}}
)

var ErrNotFound = errors.New("not found")

const (
	Prefix{{.modelName}}CacheKey = "#{{.modelName}}#cache" //todo please modify this prefix
)

type (
	{{.modelName}}Model struct {
		*mongoc.Model
	}

	{{.modelName}} struct {
		{{.modelFields}}
	}
)

func New{{.modelName}}Model(url, database, collection string, c cache.CacheConf, opts ...cache.Option) *{{.modelName}}Model {
	return &{{.modelName}}Model{mongoc.MustNewModel(url, database, collection, c, opts...)}
}

func (m *{{.modelName}}Model) FindOne(id string) (*{{.modelName}}, error) {
	key := Prefix{{.modelName}}CacheKey + id
	session, err := m.Model.TakeSession()
	if err != nil {
		return nil, err
	}
	defer m.Model.PutSession(session)

	var result {{.modelName}}
	err = m.GetCollection(session).FindOneId(&result, key, bson.ObjectIdHex(id))
	switch err {
	case nil:
		return &result, nil
	case mongoc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *{{.modelName}}Model) Delete(id string) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	key := Prefix{{.modelName}}CacheKey + id
	return m.GetCollection(session).RemoveId(bson.ObjectIdHex(id), key)
}

func (m *{{.modelName}}Model) Insert(data *{{.modelName}}) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	return m.GetCollection(session).Insert(data)
}

func (m *{{.modelName}}Model) Update(data *{{.modelName}}) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	data.UpdateTime = time.Now()
	key := Prefix{{.modelName}}CacheKey + data.Id.Hex()
	return m.GetCollection(session).UpdateId(data.Id, data, key)
}
`
	cacheSetFieldtemplate = `func (m *{{.modelName}}Model) Set{{.Name}}(id string, {{.name}} {{.type}}) error {
	_, err := m.cache.Del(Prefix{{.modelName}}CacheKey + id)
	if err != nil {
		return err
	}

	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	update := bson.M{"$set": bson.M{"{{.name}}": {{.name}}, "updateTime": time.Now()}}
	return m.GetCollection(session).UpdateId(bson.ObjectIdHex(id), update)
}`

	noCacheSetFieldtemplate = `func (m *{{.modelName}}Model) Set{{.Name}}(id string, {{.name}} {{.type}}) error {
	session, err := m.TakeSession()
	if err != nil {
		return err
	}
	defer m.PutSession(session)

	update := bson.M{"$set": bson.M{"{{.name}}": {{.name}}, "updateTime": time.Now()}}
	return m.GetCollection(session).UpdateId(bson.ObjectIdHex(id), update)
}`

	noCacheGetTemplate = `func (m *{{.modelName}}Model) GetBy{{.Name}}({{.name}} {{.type}}) (*{{.modelName}},error) {
	session, err := m.TakeSession()
	if err != nil {
		return nil,err
	}
	defer m.PutSession(session)
	var result {{.modelName}}
	query := bson.M{"{{.name}}":{{.name}}}
	err = m.GetCollection(session).FindOneNoCache(&result, query)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil,ErrNotFound
		}
		return nil,err
	}
	return &result,nil
}`
	// GetByField return single model
	getTemplate = `func (m *{{.modelName}}Model) GetBy{{.Name}}({{.name}} {{.type}}) (*{{.modelName}},error) {
	session, err := m.TakeSession()
	if err != nil {
		return nil,err
	}
	defer m.PutSession(session)
	var result {{.modelName}}
	query := bson.M{"{{.name}}":{{.name}}}
	key := getCachePrimaryKeyBy{{.Name}}({{.name}})
	err = m.GetCollection(session).FindOne(&result,key,query)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil,ErrNotFound
		}
		return nil,err
	}
	return &result,nil
}

func getCachePrimaryKeyBy{{.Name}}({{.name}} {{.type}}) string {
	return "" //todo 请补全这里
}
`

	findTemplate = `func (m *{{.modelName}}Model) FindBy{{.Name}}({{.name}} string) ([]{{.modelName}},error) {
	session, err := m.TakeSession()
	if err != nil {
		return nil,err
	}
	defer m.PutSession(session)
	
	query := bson.M{"{{.name}}":{{.name}}}
	var result []{{.modelName}}
	err = m.GetCollection(session).FindAllNoCache(&result,query)
	if err != nil {
		return nil,err
	}
	return result,nil
}`
)
