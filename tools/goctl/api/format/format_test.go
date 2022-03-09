package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	notFormattedStr = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + `  
}
type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
  Students []Student ` + "`" + `json:"students"` + "`" + `
}
service A-api {
@server(
handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response)
}
`

	formattedStr = `type Request {
	Name string ` + "`" + `path:"name,options=you|me"` + "`" + `
}
type Response {
	Message  string    ` + "`" + `json:"message"` + "`" + `
	Students []Student ` + "`" + `json:"students"` + "`" + `
}
service A-api {
	@server(
		handler: GreetHandler
	)
	get /greet/from/:name(Request) returns (Response)
}`
)

func TestFormat(t *testing.T) {
	r, err := apiFormat(notFormattedStr, true)
	assert.Nil(t, err)
	assert.Equal(t, formattedStr, r)
	_, err = apiFormat(notFormattedStr, false)
	assert.Errorf(t, err, " line 7:13 can not found declaration 'Student' in context")
}
