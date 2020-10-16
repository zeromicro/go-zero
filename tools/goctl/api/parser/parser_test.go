package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testApiTemplate = `
info(
    title: doc title
    desc: >
    doc description first part,
    doc description second part<
    version: 1.0
)

type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service A-api {
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)

  @server(
    handler: NoResponseHandler
  )
  get /greet/get(Request) returns
}
`

const testMultiServiceTemplate = `
info(
    title: doc title
    desc: >
    doc description first part,
    doc description second part<
    version: 1.0
)

type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + ` 
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service A-api {
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}

service A-api {
  @server(
    handler: NoResponseHandler
  )
  get /greet/get(Request) returns
}
`

const apiNoInfo = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service A-api {
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

const invalidApiFile = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service A-api
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

func TestParser(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(testApiTemplate), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	assert.Equal(t, len(api.Types), 2)
	assert.Equal(t, len(api.Service.Routes), 2)

	assert.Equal(t, api.Service.Routes[0].Path, "/greet/from/:name")
	assert.Equal(t, api.Service.Routes[1].Path, "/greet/get")

	assert.Equal(t, api.Service.Routes[1].RequestType.Name, "Request")
	assert.Equal(t, api.Service.Routes[1].ResponseType.Name, "")
}

func TestMultiService(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(testMultiServiceTemplate), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	assert.Equal(t, len(api.Service.Routes), 2)
	assert.Equal(t, len(api.Service.Groups), 2)
}

func TestApiNoInfo(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(apiNoInfo), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)
}

func TestInvalidApiFile(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(invalidApiFile), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	_, err = NewParser(filename)
	assert.NotNil(t, err)
}
