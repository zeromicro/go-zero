package model

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/model/mongo"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/command"
)

var (
	Cmd = &cobra.Command{
		Use:   "model",
		Short: "generate model code",
	}

	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "generate mysql model",
	}

	ddlCmd = &cobra.Command{
		Use:   "ddl",
		Short: "generate mysql model from ddl",
		RunE:  command.MysqlDDL,
	}

	datasourceCmd = &cobra.Command{
		Use:   "datasource",
		Short: "generate model from datasource",
		RunE:  command.MySqlDataSource,
	}

	pgCmd = &cobra.Command{
		Use:   "pg",
		Short: "generate postgresql model",
		RunE:  command.PostgreSqlDataSource,
	}

	mongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: "generate mongo model",
		RunE:  mongo.Action,
	}
)

func init() {
	ddlCmd.Flags().StringVarP(&command.VarStringSrc, "src", "s", "", "the path or path globbing patterns of the ddl")
	ddlCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "the target dir")
	ddlCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	ddlCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "generate code with cache [optional]")
	ddlCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "for idea plugin [optional]")
	ddlCmd.Flags().StringVar(&command.VarStringDatabase, "database", "", "the name of database [optional]")
	ddlCmd.Flags().StringVar(&command.VarStringHome, "home", "", "the goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	ddlCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "the remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	ddlCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "the branch of the remote repo, it does work with --remote")

	datasourceCmd.Flags().StringVar(&command.VarStringURL, "url", "", `the data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`)
	datasourceCmd.Flags().StringSliceVarP(&command.VarStringSliceTable, "table", "t", nil, "the table or table globbing patterns in the database")
	datasourceCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "generate code with cache [optional]")
	datasourceCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "the target dir")
	datasourceCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	datasourceCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "for idea plugin [optional]")
	datasourceCmd.Flags().StringVar(&command.VarStringHome, "home", "", "the goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	datasourceCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "the remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	datasourceCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "the branch of the remote repo, it does work with --remote")

	pgCmd.Flags().StringVar(&command.VarStringURL, "url", "", `the data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`)
	pgCmd.Flags().StringVarP(&command.VarStringTable, "table", "t", "", "the table or table globbing patterns in the database")
	pgCmd.Flags().StringVarP(&command.VarStringSchema, "schema", "s", "public", "the table schema")
	pgCmd.Flags().BoolVarP(&command.VarBoolCache, "cache", "c", false, "generate code with cache [optional]")
	pgCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "the target dir")
	pgCmd.Flags().StringVar(&command.VarStringStyle, "style", "", "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	pgCmd.Flags().BoolVar(&command.VarBoolIdea, "idea", false, "for idea plugin [optional]")
	pgCmd.Flags().StringVar(&command.VarStringHome, "home", "", "the goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	pgCmd.Flags().StringVar(&command.VarStringRemote, "remote", "", "the remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	pgCmd.Flags().StringVar(&command.VarStringBranch, "branch", "", "the branch of the remote repo, it does work with --remote")

	mongoCmd.Flags().StringSliceVarP(&mongo.VarStringSliceType, "type", "t", nil, "specified model type name")
	mongoCmd.Flags().BoolVarP(&mongo.VarBoolCache, "cache", "c", false, "generate code with cache [optional]")
	mongoCmd.Flags().StringVarP(&mongo.VarStringDir, "dir", "d", "", "the target dir")
	mongoCmd.Flags().StringVar(&mongo.VarStringStyle, "style", "", "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	mongoCmd.Flags().StringVar(&mongo.VarStringHome, "home", "", "the goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	mongoCmd.Flags().StringVar(&mongo.VarStringRemote, "remote", "", "the remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	mongoCmd.Flags().StringVar(&mongo.VarStringBranch, "branch", "", "the branch of the remote repo, it does work with --remote")

	mysqlCmd.AddCommand(datasourceCmd)
	mysqlCmd.AddCommand(ddlCmd)
	mysqlCmd.AddCommand(pgCmd)
	Cmd.AddCommand(mysqlCmd)
	Cmd.AddCommand(mongoCmd)
}
