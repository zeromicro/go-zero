package model

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/model/mongo"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/command"
)

var (
	// Cmd describes a model command.
	Cmd = &cobra.Command{
		Use:   "model",
		Short: "Generate model code",
	}

	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "Generate mysql model",
	}

	ddlCmd = &cobra.Command{
		Use:   "ddl",
		Short: "Generate mysql model from ddl",
		RunE:  command.MysqlDDL,
	}

	datasourceCmd = &cobra.Command{
		Use:   "datasource",
		Short: "Generate model from datasource",
		RunE:  command.MySqlDataSource,
	}

	pgCmd = &cobra.Command{
		Use:   "pg",
		Short: "Generate postgresql model",
		RunE:  command.PostgreSqlDataSource,
	}

	pgDatasourceCmd = &cobra.Command{
		Use:   "datasource",
		Short: "Generate model from datasource",
		RunE:  command.PostgreSqlDataSource,
	}

	mongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: "Generate mongo model",
		RunE:  mongo.Action,
	}
)

func init() {
	ddlCmd.Flags().StringVarP(&command.VarStringSrc, "src", "s", "", "The path or path globbing patterns of the ddl")
	ddlCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "The target dir")
	ddlCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	ddlCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "Generate code with cache [optional]")
	ddlCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "For idea plugin [optional]")
	ddlCmd.Flags().StringVar(&command.VarStringDatabase, "database", "", "The name of database [optional]")
	ddlCmd.Flags().StringVar(&command.VarStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	ddlCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	ddlCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")

	datasourceCmd.Flags().StringVar(&command.VarStringURL, "url", "", `The data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`)
	datasourceCmd.Flags().StringSliceVarP(&command.VarStringSliceTable, "table", "t", nil, "The table or table globbing patterns in the database")
	datasourceCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "Generate code with cache [optional]")
	datasourceCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "The target dir")
	datasourceCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	datasourceCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "For idea plugin [optional]")
	datasourceCmd.Flags().StringVar(&command.VarStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	datasourceCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	datasourceCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")

	pgDatasourceCmd.Flags().StringVar(&command.VarStringURL, "url", "", `The data source of database,like "postgres://root:password@127.0.0.1:5432/database?sslmode=disable"`)
	pgDatasourceCmd.Flags().StringVarP(&command.VarStringTable, "table", "t", "", "The table or table globbing patterns in the database")
	pgDatasourceCmd.Flags().StringVarP(&command.VarStringSchema, "schema", "s", "public", "The table schema")
	pgDatasourceCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "Generate code with cache [optional]")
	pgDatasourceCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "The target dir")
	pgDatasourceCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	pgDatasourceCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "For idea plugin [optional]")
	pgDatasourceCmd.Flags().BoolVar(&command.VarBoolStrict, "strict", false, "Generate model in strict mode")
	pgDatasourceCmd.Flags().StringVar(&command.VarStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	pgDatasourceCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	pgDatasourceCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")

	mongoCmd.Flags().StringSliceVarP(&mongo.VarStringSliceType, "type", "t", nil, "Specified model type name")
	mongoCmd.Flags().BoolVarP(&mongo.VarBoolCache, "cache", "c", false, "Generate code with cache [optional]")
	mongoCmd.Flags().BoolVarP(&mongo.VarBoolEasy, "easy", "e", false, "Generate code with auto generated CollectionName for easy declare [optional]")
	mongoCmd.Flags().StringVarP(&mongo.VarStringDir, "dir", "d", "", "The target dir")
	mongoCmd.Flags().StringVar(&mongo.VarStringStyle, "style", "", "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	mongoCmd.Flags().StringVar(&mongo.VarStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	mongoCmd.Flags().StringVar(&mongo.VarStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	mongoCmd.Flags().StringVar(&mongo.VarStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")

	mysqlCmd.PersistentFlags().BoolVar(&command.VarBoolStrict, "strict", false, "Generate model in strict mode")
	mysqlCmd.PersistentFlags().StringSliceVarP(&command.VarStringSliceIgnoreColumns, "ignore-columns", "i", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "Ignore columns while creating or updating rows")

	mysqlCmd.AddCommand(datasourceCmd)
	mysqlCmd.AddCommand(ddlCmd)
	pgCmd.AddCommand(pgDatasourceCmd)
	Cmd.AddCommand(mysqlCmd)
	Cmd.AddCommand(mongoCmd)
	Cmd.AddCommand(pgCmd)
}
