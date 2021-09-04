package new

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const apiTemplate = `
type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + ` 
}

type Response {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service {{.name}}-api {
  @handler {{.handler}}Handler
  get /from/:name(Request) returns (Response);
}
`

// CreateServiceCommand fast create service
func CreateServiceCommand(c *cli.Context) error {
	args := c.Args()
	dirName := args.First()
	if len(dirName) == 0 {
		dirName = "greet"
	}

	if strings.Contains(dirName, "-") {
		errorx.Must(errors.New("api new command service name not support strikethrough, because this will used by function name"))
	}

	abs, err := filepath.Abs(dirName)
	errorx.Must(err)

	err = util.MkdirIfNotExist(abs)
	errorx.Must(err)

	dirName = filepath.Base(filepath.Clean(abs))
	filename := dirName + ".api"
	apiFilePath := filepath.Join(abs, filename)
	fp, err := os.Create(apiFilePath)
	errorx.Must(err)

	defer fp.Close()
	t := template.Must(template.New("template").Parse(apiTemplate))
	errorx.Must(t.Execute(fp, map[string]string{
		"name":    dirName,
		"handler": strings.Title(dirName),
	}))

	errorx.Must(gogen.DoGenProject(apiFilePath, abs, conf.DefaultFormat))
	return nil
}
