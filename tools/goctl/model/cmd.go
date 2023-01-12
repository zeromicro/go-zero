package model

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/model/mongo"
)

var (
	// Cmd describes a model command.
	Cmd = &cobra.Command{
		Use:   "model",
		Short: "Generate model code",
	}

	mongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: "Generate mongo model",
		RunE:  mongo.Action,
	}
)

func init() {
	mongoCmd.Flags().StringSliceVarP(&mongo.VarStringSliceType, "type", "t", nil, "Specified model type name")
	mongoCmd.Flags().BoolVarP(&mongo.VarBoolCache, "cache", "c", false, "Generate code with cache [optional]")
	mongoCmd.Flags().BoolVarP(&mongo.VarBoolEasy, "easy", "e", false, "Generate code with auto generated CollectionName for easy declare [optional]")
	mongoCmd.Flags().StringVarP(&mongo.VarStringDir, "dir", "d", "", "The target dir")
	mongoCmd.Flags().StringVar(&mongo.VarStringStyle, "style", "", "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	mongoCmd.Flags().StringVar(&mongo.VarStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	mongoCmd.Flags().StringVar(&mongo.VarStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	mongoCmd.Flags().StringVar(&mongo.VarStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")

	Cmd.AddCommand(mongoCmd)
}
