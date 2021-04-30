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
	Message string ` + "`" + `json:"message"` + "`" + `
}
service A-api {
	@server(
		handler: GreetHandler
	)
	get /greet/from/:name(Request) returns (Response)
}`
)

func TestFormat(t *testing.T) {
	r, err := apiFormat(notFormattedStr)
	assert.Nil(t, err)
	assert.Equal(t, r, formattedStr)
}
