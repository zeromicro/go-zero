package command

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/postgres"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/util"
	file "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/urfave/cli"
)

const (
	flagSrc      = "src"
	flagDir      = "dir"
	flagCache    = "cache"
	flagIdea     = "idea"
	flagURL      = "url"
	flagTable    = "table"
	flagStyle    = "style"
	flagDatabase = "database"
	flagSchema   = "schema"
	flagHome     = "home"
)

var errNotMatched = errors.New("sql not matched")

// MysqlDDL generates model code from ddl
func MysqlDDL(ctx *cli.Context) error {
	src := ctx.String(flagSrc)
	dir := ctx.String(flagDir)
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	style := ctx.String(flagStyle)
	database := ctx.String(flagDatabase)
	home := ctx.String(flagHome)
	remote := ctx.String("remote")
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		file.RegisterGoctlHome(home)
	}
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	return fromDDL(src, dir, cfg, cache, idea, database)
}

// MySqlDataSource generates model code from datasource
func MySqlDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagURL))
	dir := strings.TrimSpace(ctx.String(flagDir))
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	style := ctx.String(flagStyle)
	home := ctx.String("home")
	remote := ctx.String("remote")
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		file.RegisterGoctlHome(home)
	}

	pattern := strings.TrimSpace(ctx.String(flagTable))
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	return fromMysqlDataSource(url, pattern, dir, cfg, cache, idea)
}

// PostgreSqlDataSource generates model code from datasource
func PostgreSqlDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagURL))
	dir := strings.TrimSpace(ctx.String(flagDir))
	cache := ctx.Bool(flagCache)
	idea := ctx.Bool(flagIdea)
	style := ctx.String(flagStyle)
	schema := ctx.String(flagSchema)
	home := ctx.String("home")
	remote := ctx.String("remote")
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		file.RegisterGoctlHome(home)
	}

	if len(schema) == 0 {
		schema = "public"
	}

	pattern := strings.TrimSpace(ctx.String(flagTable))
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	return fromPostgreSqlDataSource(url, pattern, dir, schema, cfg, cache, idea)
}

func fromDDL(src, dir string, cfg *config.Config, cache, idea bool, database string) error {
	log := console.NewConsole(idea)
	src = strings.TrimSpace(src)
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}

	files, err := util.MatchFiles(src)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNotMatched
	}

	generator, err := gen.NewDefaultGenerator(dir, cfg, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = generator.StartFromDDL(file, cache, database)
		if err != nil {
			return err
		}
	}

	return nil
}

func fromMysqlDataSource(url, pattern, dir string, cfg *config.Config, cache, idea bool) error {
	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of mysql, but nothing found")
		return nil
	}

	if len(pattern) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	dsn, err := mysql.ParseDSN(url)
	if err != nil {
		return err
	}

	logx.Disable()
	databaseSource := strings.TrimSuffix(url, "/"+dsn.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := im.FindColumns(dsn.DBName, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, cfg, gen.WithConsoleOption(log))
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, cache)
}

func fromPostgreSqlDataSource(url, pattern, dir, schema string, cfg *config.Config, cache, idea bool) error {
	log := console.NewConsole(idea)
	if len(url) == 0 {
		log.Error("%v", "expected data source of postgresql, but nothing found")
		return nil
	}

	if len(pattern) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}
	db := postgres.New(url)
	im := model.NewPostgreSqlModel(db)

	tables, err := im.GetAllTables(schema)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := im.FindColumns(schema, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, cfg, gen.WithConsoleOption(log), gen.WithPostgreSql())
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, cache)
}
