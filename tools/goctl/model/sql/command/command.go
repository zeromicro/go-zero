package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/urfave/cli"
)

const (
	flagSrc   = "src"
	flagDir   = "dir"
	flagCache = "cache"
	flagIdea  = "idea"
	flagStyle = "style"
	flagUrl   = "url"
	flagTable = "table"
)

func MysqlDDL(ctx *cli.Context) error {
	src := ctx.String(flagSrc)
	dir := ctx.String(flagDir)
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	namingStyle := strings.TrimSpace(ctx.String(flagStyle))
	log := console.NewConsole(idea)
	src = strings.TrimSpace(src)
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}

	switch namingStyle {
	case gen.NamingLower, gen.NamingCamel, gen.NamingSnake:
	case "":
		namingStyle = gen.NamingLower
	default:
		return fmt.Errorf("unexpected naming style: %s", namingStyle)
	}

	files, err := util.MatchFiles(src)
	if err != nil {
		return err
	}

	var source []string
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		source = append(source, string(data))
	}
	generator := gen.NewDefaultGenerator(strings.Join(source, "\n"), dir, namingStyle, gen.WithConsoleOption(log))
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
	namingStyle := strings.TrimSpace(ctx.String(flagStyle))
	pattern := strings.TrimSpace(ctx.String(flagTable))
	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of mysql, but nothing found")
		return nil
	}

	if len(pattern) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	switch namingStyle {
	case gen.NamingLower, gen.NamingCamel, gen.NamingSnake:
	case "":
		namingStyle = gen.NamingLower
	default:
		return fmt.Errorf("unexpected naming style: %s", namingStyle)
	}

	cfg, err := mysql.ParseDSN(url)
	if err != nil {
		return err
	}

	logx.Disable()
	conn := sqlx.NewMysql(url)
	databaseSource := strings.TrimSuffix(url, "/"+cfg.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	m := model.NewDDLModel(conn)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(cfg.DBName)
	if err != nil {
		return err
	}

	var matchTables []string
	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		matchTables = append(matchTables, item)
	}
	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	ddl, err := m.ShowDDL(matchTables...)
	if err != nil {
		log.Error("%v", err)
		return nil
	}

	generator := gen.NewDefaultGenerator(strings.Join(ddl, "\n"), dir, namingStyle, gen.WithConsoleOption(log))
	err = generator.Start(cache)
	if err != nil {
		log.Error("%v", err)
	}

	return nil
}
