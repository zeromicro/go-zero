package new

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/urfave/cli"
)

const apiTemplate = `
type Request struct {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + ` // 框架自动验证请求参数是否合法
}

type Response struct {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service {{.name}}-api {
  @server(
    handler: GreetHandler
  )
  get /greet/from/:name(Request) returns (Response);
}
`

func NewService(c *cli.Context) error {
	var args = os.Args
	if len(args) <= 2 {
		return errors.New("invalid args, eg: goctl api new [greet]")
	}
	name := "greet"
	if args[len(args)-1] != "new" {
		name = args[len(args)-1]
	}
	location := name
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return err
	}

	filename := name + ".api"
	goPath := filepath.Join(location, filename)
	fp, err := os.Create(goPath)
	if err != nil {
		return err
	}

	defer fp.Close()
	t := template.Must(template.New("template").Parse(apiTemplate))
	if err := t.Execute(fp, map[string]string{
		"name": name,
	}); err != nil {
		return err
	}

	_, err = execx.Run(fmt.Sprintf("goctl api go -api %s -dir .", filename), location)
	return err
}
