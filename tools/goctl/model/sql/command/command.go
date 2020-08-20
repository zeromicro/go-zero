package command

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/urfave/cli"
)

const (
	flagSrc   = "src"
	flagDir   = "dir"
	flagCache = "cache"
	flagIdea  = "idea"
	flagUrl   = "url"
	flagTable = "table"
)

func MysqlDDL(ctx *cli.Context) error {
	src := ctx.String(flagSrc)
	dir := ctx.String(flagDir)
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	log := console.NewConsole(idea)
	fileSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(fileSrc)
	if err != nil {
		return err
	}
	source := string(data)
	generator := gen.NewDefaultGenerator(source, dir, gen.WithConsoleOption(log))
	err = generator.Start(cache)
	if err != nil {
		log.Error("%v", err)
	}
	return nil
}

func MyDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagUrl))
	dir := strings.TrimSpace(ctx.String(flagDir))
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	table := strings.TrimSpace(ctx.String(flagTable))
	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of mysql, but is empty")
		return nil
	}
	if len(table) == 0 {
		log.Error("%v", "expected table(s), but nothing found")
		return nil
	}
	logx.Disable()
	conn := sqlx.NewMysql(url)
	m := model.NewDDLModel(conn)
	tables := collection.NewSet()
	for _, item := range strings.Split(table, ",") {
		item = strings.TrimSpace(item)
		if len(item) == 0 {
			continue
		}
		tables.AddStr(item)
	}
	ddl, err := m.ShowDDL(tables.KeysStr()...)
	if err != nil {
		log.Error("%v", err)
		return nil
	}
	generator := gen.NewDefaultGenerator(strings.Join(ddl, "\n"), dir, gen.WithConsoleOption(log))
	err = generator.Start(cache)
	if err != nil {
		log.Error("%v", err)
	}
	return nil
}
