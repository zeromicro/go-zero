package command

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/command/migrationnotes"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/gen"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/util"
	file "github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringSrc describes the source file of sql.
	VarStringSrc string
	// VarStringDir describes the output directory of sql.
	VarStringDir string
	// VarBoolCache describes whether the cache is enabled.
	VarBoolCache bool
	// VarBoolIdea describes whether is idea or not.
	VarBoolIdea bool
	// VarStringURL describes the dsn of the sql.
	VarStringURL string
	// VarStringSliceTable describes tables.
	VarStringSliceTable []string
	// VarStringStyle describes the style.
	VarStringStyle string
	// VarStringDatabase describes the database.
	VarStringDatabase string
	// VarStringSchema describes the schema of postgresql.
	VarStringSchema string
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch of the repository.
	VarStringBranch string
	// VarBoolStrict describes whether the strict mode is enabled.
	VarBoolStrict bool
	// VarStringSliceIgnoreColumns represents the columns which are ignored.
	VarStringSliceIgnoreColumns []string
	// VarStringCachePrefix describes the prefix of cache.
	VarStringCachePrefix string
)

var errNotMatched = errors.New("sql not matched")

// MysqlDDL generates model code from ddl
func MysqlDDL(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	if VarBoolCache && len(VarStringCachePrefix) == 0 {
		return errors.New("cache prefix is empty")
	}
	src := VarStringSrc
	dir := VarStringDir
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	database := VarStringDatabase
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	arg := ddlArg{
		src:           src,
		dir:           dir,
		cfg:           cfg,
		cache:         cache,
		idea:          idea,
		database:      database,
		strict:        VarBoolStrict,
		ignoreColumns: mergeColumns(VarStringSliceIgnoreColumns),
		prefix:        VarStringCachePrefix,
	}
	return fromDDL(arg)
}

// MySqlDataSource generates model code from datasource
func MySqlDataSource(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	if VarBoolCache && len(VarStringCachePrefix) == 0 {
		return errors.New("cache prefix is empty")
	}
	url := strings.TrimSpace(VarStringURL)
	dir := strings.TrimSpace(VarStringDir)
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	tableValue := VarStringSliceTable
	patterns := parseTableList(tableValue)
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	arg := dataSourceArg{
		url:           url,
		dir:           dir,
		tablePat:      patterns,
		cfg:           cfg,
		cache:         cache,
		idea:          idea,
		strict:        VarBoolStrict,
		ignoreColumns: mergeColumns(VarStringSliceIgnoreColumns),
		prefix:        VarStringCachePrefix,
	}
	return fromMysqlDataSource(arg)
}

func mergeColumns(columns []string) []string {
	set := collection.NewSet()
	for _, v := range columns {
		fields := strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
		set.AddStr(fields...)
	}
	return set.KeysStr()
}

type pattern map[string]struct{}

func (p pattern) Match(s string) bool {
	for v := range p {
		match, err := filepath.Match(v, s)
		if err != nil {
			console.Error("%+v", err)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func (p pattern) list() []string {
	var ret []string
	for v := range p {
		ret = append(ret, v)
	}
	return ret
}

func parseTableList(tableValue []string) pattern {
	tablePattern := make(pattern)
	for _, v := range tableValue {
		fields := strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
		for _, f := range fields {
			tablePattern[f] = struct{}{}
		}
	}
	return tablePattern
}

// PostgreSqlDataSource generates model code from datasource
func PostgreSqlDataSource(_ *cobra.Command, _ []string) error {
	migrationnotes.BeforeCommands(VarStringDir, VarStringStyle)
	url := strings.TrimSpace(VarStringURL)
	dir := strings.TrimSpace(VarStringDir)
	cache := VarBoolCache
	idea := VarBoolIdea
	style := VarStringStyle
	schema := VarStringSchema
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(schema) == 0 {
		schema = "public"
	}

	patterns := parseTableList(VarStringSliceTable)
	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}
	ignoreColumns := mergeColumns(VarStringSliceIgnoreColumns)

	return fromPostgreSqlDataSource(url, patterns, dir, schema, cfg, cache, idea, VarBoolStrict, ignoreColumns)
}

type ddlArg struct {
	src, dir      string
	cfg           *config.Config
	cache, idea   bool
	database      string
	strict        bool
	ignoreColumns []string
	prefix        string
}

func fromDDL(arg ddlArg) error {
	log := console.NewConsole(arg.idea)
	src := strings.TrimSpace(arg.src)
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

	generator, err := gen.NewDefaultGenerator(arg.prefix, arg.dir, arg.cfg,
		gen.WithConsoleOption(log), gen.WithIgnoreColumns(arg.ignoreColumns))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = generator.StartFromDDL(file, arg.cache, arg.strict, arg.database)
		if err != nil {
			return err
		}
	}

	return nil
}

type dataSourceArg struct {
	url, dir      string
	tablePat      pattern
	cfg           *config.Config
	cache, idea   bool
	strict        bool
	ignoreColumns []string
	prefix        string
}

func fromMysqlDataSource(arg dataSourceArg) error {
	log := console.NewConsole(arg.idea)
	if len(arg.url) == 0 {
		log.Error("%v", "expected data source of mysql, but nothing found")
		return nil
	}

	if len(arg.tablePat) == 0 {
		log.Error("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	dsn, err := mysql.ParseDSN(arg.url)
	if err != nil {
		return err
	}

	logx.Disable()
	databaseSource := strings.TrimSuffix(arg.url, "/"+dsn.DBName) + "/information_schema"
	db := sqlx.NewMysql(databaseSource)
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		if !arg.tablePat.Match(item) {
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

	generator, err := gen.NewDefaultGenerator(arg.prefix, arg.dir, arg.cfg,
		gen.WithConsoleOption(log), gen.WithIgnoreColumns(arg.ignoreColumns))
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, arg.cache, arg.strict)
}

func fromPostgreSqlDataSource(url string, pattern pattern, dir, schema string, cfg *config.Config, cache, idea, strict bool, ignoreColumns []string) error {
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
		if !pattern.Match(item) {
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

	generator, err := gen.NewDefaultGenerator("", dir, cfg, gen.WithConsoleOption(log), gen.WithPostgreSql(), gen.WithIgnoreColumns(ignoreColumns))
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, cache, strict)
}
