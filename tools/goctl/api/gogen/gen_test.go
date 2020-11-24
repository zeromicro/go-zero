package gogen

import (
	goformat "go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

const testApiTemplate = `
info(
    title: doc title
    desc: >
    doc description first part,
    doc description second part<
    version: 1.0
)

// TODO: test
// {
type Request struct {  // TODO: test
  // TOOD
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `   // }
} // TODO: test

// TODO: test
type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

@server(
    // C0
	group: greet/s1
)
// C1
service A-api {
  // C2
  @server( // C3
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)   // hello
	
  // C4
  @handler NoResponseHandler  // C5
  get /greet/get(Request)
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

const anonymousAnnotation = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request) returns (Response)
}
`

const apiHasMiddleware = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

@server(
	middleware: TokenValidate
)
service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request) returns (Response)
}
`

const apiJwt = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

@server(
	jwt: Auth
)
service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request) returns (Response)
}
`

const apiJwtWithMiddleware = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

@server(
	jwt: Auth
	middleware: TokenValidate
)
service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request) returns (Response)
}
`

const apiHasNoRequest = `
service A-api {
  @handler GreetHandler
  post /greet/ping ()
}
`

const apiRouteTest = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}
type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}
service A-api {
  @handler NormalHandler
  get /greet/from/:name(Request) returns (Response)
  @handler NoResponseHandler
  get /greet/from/:sex(Request)
  @handler NoRequestHandler
  get /greet/from/request returns (Response)
  @handler NoRequestNoResponseHandler
  get /greet/from
}
`

const hasCommentApiTest = `
type Inline struct {

}

type Request struct {
  Inline 
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + ` // name in path
}

type Response struct {
  Message string ` + "`" + `json:"msg"` + "`" + ` // message
}

service A-api {
  @doc(helloworld)
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

const hasInlineNoExistTest = `

type Request struct {
  Inline 
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + ` // message
}

service A-api {
  @doc(helloworld)
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

const importApi = `
type ImportData struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

`

const hasImportApi = `
import "importApi.api"

type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + ` // message
}

service A-api {
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

const noStructTagApi = `
type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}

type XXX {}

type (
	Response {
  		Message string ` + "`" + `json:"message"` + "`" + `
	}

	A {}

	B struct {}
)

service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request) returns (Response)
}
`

const nestTypeApi = `
type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
  XXX struct {
  }
}

service A-api {
  @handler GreetHandler
  get /greet/from/:name(Request)
}
`

func TestParser(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(testApiTemplate), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	assert.Equal(t, len(api.Types), 2)
	assert.Equal(t, len(api.Service.Routes()), 2)

	assert.Equal(t, api.Service.Routes()[0].Path, "/greet/from/:name")
	assert.Equal(t, api.Service.Routes()[1].Path, "/greet/get")

	assert.Equal(t, api.Service.Routes()[1].RequestType.Name, "Request")
	assert.Equal(t, api.Service.Routes()[1].ResponseType.Name, "")

	validate(t, filename)
}

func TestMultiService(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(testMultiServiceTemplate), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	assert.Equal(t, len(api.Service.Routes()), 2)
	assert.Equal(t, len(api.Service.Groups), 2)

	validate(t, filename)
}

func TestApiNoInfo(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(apiNoInfo), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestInvalidApiFile(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(invalidApiFile), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	_, err = parser.NewParser(filename)
	assert.NotNil(t, err)
}

func TestAnonymousAnnotation(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(anonymousAnnotation), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	assert.Equal(t, len(api.Service.Routes()), 1)
	assert.Equal(t, api.Service.Routes()[0].Annotations[0].Value, "GreetHandler")

	validate(t, filename)
}

func TestApiHasMiddleware(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(apiHasMiddleware), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestApiHasJwt(t *testing.T) {
	filename := "jwt.api"
	err := ioutil.WriteFile(filename, []byte(apiJwt), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestApiHasJwtAndMiddleware(t *testing.T) {
	filename := "jwt.api"
	err := ioutil.WriteFile(filename, []byte(apiJwtWithMiddleware), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestApiHasNoRequestBody(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(apiHasNoRequest), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestApiRoutes(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(apiRouteTest), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestHasCommentRoutes(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(hasCommentApiTest), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestInlineTypeNotExist(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(hasInlineNoExistTest), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	_, err = parser.Parse()
	assert.Nil(t, err)

	validate(t, filename)
}

func TestHasImportApi(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(hasImportApi), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	importApiName := "importApi.api"
	err = ioutil.WriteFile(importApiName, []byte(importApi), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(importApiName)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	api, err := parser.Parse()
	assert.Nil(t, err)

	var hasInline bool
	for _, ty := range api.Types {
		if ty.Name == "ImportData" {
			hasInline = true
			break
		}
	}
	assert.True(t, hasInline)

	validate(t, filename)
}

func TestNoStructApi(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(noStructTagApi), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)

	parser, err := parser.NewParser(filename)
	assert.Nil(t, err)

	spec, err := parser.Parse()
	assert.Nil(t, err)
	assert.Equal(t, len(spec.Types), 5)

	validate(t, filename)
}

func TestNestTypeApi(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(nestTypeApi), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)
	_, err = parser.NewParser(filename)

	assert.NotNil(t, err)
}

func TestCamelStyle(t *testing.T) {
	filename := "greet.api"
	err := ioutil.WriteFile(filename, []byte(testApiTemplate), os.ModePerm)
	assert.Nil(t, err)
	defer os.Remove(filename)
	_, err = parser.NewParser(filename)
	assert.Nil(t, err)

	validateWithCamel(t, filename, "GoZero")
}

func validate(t *testing.T, api string) {
	validateWithCamel(t, api, "gozero")
}

func validateWithCamel(t *testing.T, api, camel string) {
	dir := "_go"
	os.RemoveAll(dir)
	err := DoGenProject(api, dir, camel)
	defer os.RemoveAll(dir)
	assert.Nil(t, err)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			code, err := ioutil.ReadFile(path)
			assert.Nil(t, err)
			assert.Nil(t, validateCode(string(code)))
		}
		return nil
	})

	_, err = execx.Run("go test ./...", dir)
	assert.Nil(t, err)
}

func validateCode(code string) error {
	_, err := goformat.Source([]byte(code))
	return err
}
