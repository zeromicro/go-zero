package new

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/util"
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
	args := c.Args()
	name := "greet"
	if len(args) > 0 {
		name = args.First()
	}

	abs, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	err = util.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	name = filepath.Base(filepath.Clean(abs))
	filename := name + ".api"
	apiFilePath := filepath.Join(abs, filename)
	fp, err := os.Create(apiFilePath)
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

	err = gogen.DoGenProject(apiFilePath, abs)
	return err
}
