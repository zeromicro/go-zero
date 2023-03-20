package model

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"

	"github.com/zeromicro/go-zero/tools/goctl/model/mongo"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/command"
)

var (
	// Cmd describes a model command.
	Cmd = &cobra.Command{
		Use:   "model",
		Short: flags.Get("model.short"),
	}

	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: flags.Get("mysql.short"),
	}

	ddlCmd = &cobra.Command{
		Use:   "ddl",
		Short: flags.Get("model.mysal.ddl.short"),
		RunE:  command.MysqlDDL,
	}

	datasourceCmd = &cobra.Command{
		Use:   "datasource",
		Short: flags.Get("model.mysql.datasource.short"),
		RunE:  command.MySqlDataSource,
	}

	pgCmd = &cobra.Command{
		Use:   "pg",
		Short: flags.Get("model.pg.short"),
		RunE:  command.PostgreSqlDataSource,
	}

	pgDatasourceCmd = &cobra.Command{
		Use:   "datasource",
		Short: flags.Get("model.pg.datasource.short"),
		RunE:  command.PostgreSqlDataSource,
	}

	mongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: flags.Get("model.mongo.short"),
		RunE:  mongo.Action,
	}
)

func init() {
	var (
		ddlCmdFlags          = ddlCmd.Flags()
		datasourceCmdFlags   = datasourceCmd.Flags()
		pgDatasourceCmdFlags = pgDatasourceCmd.Flags()
		mongoCmdFlags        = mongoCmd.Flags()
	)

	ddlCmdFlags.StringVarP(&command.VarStringSrc, "src", "s", "", flags.Get("model.mysql.ddl.src"))
	ddlCmdFlags.StringVarP(&command.VarStringDir, "dir", "d", "", flags.Get("model.mysql.ddl.dir"))
	ddlCmdFlags.StringVar(&command.VarStringStyle, "style", "", flags.Get("model.mysql.ddl.style"))
	ddlCmdFlags.BoolVarP(&command.VarBoolCache, "cache", "c", false, flags.Get("model.mysql.ddl.cache"))
	ddlCmdFlags.BoolVar(&command.VarBoolIdea, "idea", false, flags.Get("model.mysql.ddl.idea"))
	ddlCmdFlags.StringVar(&command.VarStringDatabase, "database", "", flags.Get("model.mysql.ddl.database"))
	ddlCmdFlags.StringVar(&command.VarStringHome, "home", "", flags.Get("model.mysql.ddl.home"))
	ddlCmdFlags.StringVar(&command.VarStringRemote, "remote", "", flags.Get("model.mysql.ddl.remote"))
	ddlCmdFlags.StringVar(&command.VarStringBranch, "branch", "", flags.Get("model.mysql.ddl.branch"))

	datasourceCmdFlags.StringVar(&command.VarStringURL, "url", "", flags.Get("model.mysql.datasource.url"))
	datasourceCmdFlags.StringSliceVarP(&command.VarStringSliceTable, "table", "t", nil, flags.Get("model.mysql.datasource.table"))
	datasourceCmdFlags.BoolVarP(&command.VarBoolCache, "cache", "c", false, flags.Get("model.mysql.datasource.cache"))
	datasourceCmdFlags.StringVarP(&command.VarStringDir, "dir", "d", "", flags.Get("model.mysql.datasource.dir"))
	datasourceCmdFlags.StringVar(&command.VarStringStyle, "style", "", flags.Get("model.mysql.datasource.style"))
	datasourceCmdFlags.BoolVar(&command.VarBoolIdea, "idea", false, flags.Get("model.mysql.datasource.idea"))
	datasourceCmdFlags.StringVar(&command.VarStringHome, "home", "", flags.Get("model.mysql.datasource.home"))
	datasourceCmdFlags.StringVar(&command.VarStringRemote, "remote", "", flags.Get("model.mysql.datasource.remote"))
	datasourceCmdFlags.StringVar(&command.VarStringBranch, "branch", "", flags.Get("model.mysql.datasource.branch"))

	pgDatasourceCmdFlags.StringVar(&command.VarStringURL, "url", "", flags.Get("model.pg.datasource.url"))
	pgDatasourceCmdFlags.StringVarP(&command.VarStringTable, "table", "t", "", flags.Get("model.pg.datasource.table"))
	pgDatasourceCmdFlags.StringVarP(&command.VarStringSchema, "schema", "s", "public", flags.Get("model.pg.datasource.schema"))
	pgDatasourceCmdFlags.BoolVarP(&command.VarBoolCache, "cache", "c", false, flags.Get("model.pg.datasource.cache"))
	pgDatasourceCmdFlags.StringVarP(&command.VarStringDir, "dir", "d", "", flags.Get("model.pg.datasource.dir"))
	pgDatasourceCmdFlags.StringVar(&command.VarStringStyle, "style", "", flags.Get("model.pg.datasource.style"))
	pgDatasourceCmdFlags.BoolVar(&command.VarBoolIdea, "idea", false, flags.Get("model.pg.datasource.idea"))
	pgDatasourceCmdFlags.BoolVar(&command.VarBoolStrict, "strict", false, flags.Get("model.pg.datasource.strict"))
	pgDatasourceCmdFlags.StringVar(&command.VarStringHome, "home", "", flags.Get("model.pg.datasource.home"))
	pgDatasourceCmdFlags.StringVar(&command.VarStringRemote, "remote", "", flags.Get("model.pg.datasource.remote"))
	pgDatasourceCmdFlags.StringVar(&command.VarStringBranch, "branch", "", flags.Get("model.pg.datasource.branch"))

	mongoCmdFlags.StringSliceVarP(&mongo.VarStringSliceType, "type", "t", nil, flags.Get("model.mongo.type"))
	mongoCmdFlags.BoolVarP(&mongo.VarBoolCache, "cache", "c", false, flags.Get("model.mongo.cache"))
	mongoCmdFlags.BoolVarP(&mongo.VarBoolEasy, "easy", "e", false, flags.Get("model.mongo.easy"))
	mongoCmdFlags.StringVarP(&mongo.VarStringDir, "dir", "d", "", flags.Get("model.mongo.dir"))
	mongoCmdFlags.StringVar(&mongo.VarStringStyle, "style", "", flags.Get("model.mongo.style"))
	mongoCmdFlags.StringVar(&mongo.VarStringHome, "home", "", flags.Get("model.mongo.home"))
	mongoCmdFlags.StringVar(&mongo.VarStringRemote, "remote", "", flags.Get("model.mongo.remote"))
	mongoCmdFlags.StringVar(&mongo.VarStringBranch, "branch", "", flags.Get("model.mongo.branch"))

	mysqlCmd.PersistentFlags().BoolVar(&command.VarBoolStrict, "strict", false, flags.Get("model.mysql.strict"))
	mysqlCmd.PersistentFlags().StringSliceVarP(&command.VarStringSliceIgnoreColumns, "ignore-columns", "i", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, flags.Get("model.mysql.ignore-columns"))

	mysqlCmd.AddCommand(datasourceCmd, ddlCmd)
	pgCmd.AddCommand(pgDatasourceCmd)
	Cmd.AddCommand(mysqlCmd, mongoCmd, pgCmd)
}
