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

var errNotMatched = errors.New("sql not matched")

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
	return fromDDl(src, dir, namingStyle, cache, idea)
}

func MyDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagUrl))
	dir := strings.TrimSpace(ctx.String(flagDir))
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	namingStyle := strings.TrimSpace(ctx.String(flagStyle))
	pattern := strings.TrimSpace(ctx.String(flagTable))
	return fromDataSource(url, pattern, dir, namingStyle, cache, idea)
}

func fromDDl(src, dir, namingStyle string, cache, idea bool) error {
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

	if len(files) == 0 {
		return errNotMatched
	}

	var source []string
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		source = append(source, string(data))
	}
	generator, err := gen.NewDefaultGenerator(dir, namingStyle, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}

	err = generator.StartFromDDL(strings.Join(source, "\n"), cache)
	return err
}

func fromDataSource(url, pattern, dir, namingStyle string, cache, idea bool) error {
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
	databaseSource := strings.TrimSuffix(url, "/"+cfg.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(cfg.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string][]*model.Column)
	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}
		columns, err := im.FindByTableName(cfg.DBName, item)
		if err != nil {
			return err
		}
		matchTables[item] = columns
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, namingStyle, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}

	err = generator.StartFromInformationSchema(cfg.DBName, matchTables, cache)
	return err
}
