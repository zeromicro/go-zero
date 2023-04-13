package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator/ent"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringOutput describes the output.
	VarStringOutput string
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringSliceGoOut describes the go output.
	VarStringSliceGoOut []string
	// VarStringSliceGoGRPCOut describes the grpc output.
	VarStringSliceGoGRPCOut []string
	// VarStringSlicePlugin describes the protoc plugin.
	VarStringSlicePlugin []string
	// VarStringSliceProtoPath describes the proto path.
	VarStringSliceProtoPath []string
	// VarStringSliceGoOpt describes the go options.
	VarStringSliceGoOpt []string
	// VarStringSliceGoGRPCOpt describes the grpc options.
	VarStringSliceGoGRPCOpt []string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarStringZRPCOut describes the zRPC output.
	VarStringZRPCOut string
	// VarBoolIdea describes whether idea or not
	VarBoolIdea bool
	// VarBoolVerbose describes whether verbose.
	VarBoolVerbose bool
	// VarBoolMultiple describes whether support generating multiple rpc services or not.
	VarBoolMultiple bool
	// VarStringSchema describes the ent schema path
	VarStringSchema string
	// VarStringServiceName describes the service name
	VarStringServiceName string
	// VarStringProjectName describes the service name
	VarStringProjectName string
	// VarStringModelName describes which model for generating
	VarStringModelName string
	// VarIntSearchKeyNum describes the number of search keys
	VarIntSearchKeyNum int
	// VarBoolEnt describes whether the project use Ent
	VarBoolEnt bool
	// VarStringModuleName describes the module name
	VarStringModuleName string
	// VarStringGoZeroVersion describes the version of Go Zero
	VarStringGoZeroVersion string
	// VarStringToolVersion describes the version of Simple Admin Tools
	VarStringToolVersion string
	// VarIntServicePort describes the service port exposed
	VarIntServicePort int
	// VarBoolGitlab describes whether to use gitlab-ci
	VarBoolGitlab bool
	// VarStringGroupName describes whether to use group
	VarStringGroupName string
	// VarStringProtoPath describes the output proto file path for ent code generation
	VarStringProtoPath string
	// VarStringProtoFieldStyle describes the proto fields naming style
	VarStringProtoFieldStyle string
	// VarBoolDesc describes whether to create desc folder for splitting proto files
	VarBoolDesc bool
	// VarBoolOverwrite describes whether to overwrite the files, it will overwrite all generated files.
	VarBoolOverwrite bool
)

// RPCNew is to generate rpc greet service, this greet service can speed
// up your understanding of the zrpc service structure
func RPCNew(_ *cobra.Command, args []string) error {
	rpcname := args[0]
	ext := filepath.Ext(rpcname)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	verbose := VarBoolVerbose
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	protoName, err := format.FileNamingFormat(style, rpcname)
	protoName += ".proto"
	filename := filepath.Join(".", rpcname, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Src = src
	ctx.GoOutput = filepath.Join(filepath.Dir(src), "types")
	ctx.GrpcOutput = filepath.Join(filepath.Dir(src), "types")
	ctx.IsGooglePlugin = true
	ctx.Output = filepath.Dir(src)
	ctx.ProtocCmd = fmt.Sprintf("protoc -I=%s %s --go_out=%s --go-grpc_out=%s", filepath.Dir(src), filepath.Base(src), ctx.GoOutput, ctx.GrpcOutput)
	ctx.Ent = VarBoolEnt

	if VarStringModuleName != "" {
		ctx.ModuleName = VarStringModuleName
	} else {
		ctx.ModuleName = rpcname
	}

	ctx.GoZeroVersion = VarStringGoZeroVersion
	ctx.ToolVersion = VarStringToolVersion
	ctx.Port = VarIntServicePort
	ctx.MakeFile = true
	ctx.DockerFile = true
	ctx.Gitlab = VarBoolGitlab
	ctx.UseDescDir = VarBoolDesc
	ctx.RpcName = rpcname

	if err := pathx.MkdirIfNotExist(ctx.GoOutput); err != nil {
		return err
	}

	grpcOptList := VarStringSliceGoGRPCOpt
	if len(grpcOptList) > 0 {
		ctx.ProtocCmd += " --go-grpc_opt=" + strings.Join(grpcOptList, ",")
	}

	goOptList := VarStringSliceGoOpt
	if len(goOptList) > 0 {
		ctx.ProtocCmd += " --go_opt=" + strings.Join(goOptList, ",")
	}

	g := generator.NewGenerator(style, verbose)
	return g.Generate(&ctx)
}

// RPCTemplate is the entry for generate rpc template
func RPCTemplate(latest bool) error {
	if !latest {
		console.Warning("deprecated: goctl rpc template -o is deprecated and will be removed in the future, use goctl rpc -o instead")
	}
	protoFile := VarStringOutput
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(protoFile) == 0 {
		return errors.New("missing -o")
	}

	return generator.ProtoTmpl(protoFile)
}

// EntCRUDLogic is used to generate CRUD code with Ent
func EntCRUDLogic(_ *cobra.Command, _ []string) error {
	params := &ent.GenEntLogicContext{
		Schema:          VarStringSchema,
		Output:          VarStringOutput,
		ServiceName:     VarStringServiceName,
		ProjectName:     VarStringProjectName,
		Style:           VarStringStyle,
		ModelName:       VarStringModelName,
		Multiple:        VarBoolMultiple,
		SearchKeyNum:    VarIntSearchKeyNum,
		ModuleName:      VarStringModuleName,
		GroupName:       VarStringGroupName,
		ProtoOut:        VarStringProtoPath,
		ProtoFieldStyle: VarStringProtoFieldStyle,
		Overwrite:       VarBoolOverwrite,
	}

	if params.ProjectName == "" {
		params.ProjectName = params.ServiceName
	}

	err := params.Validate()
	if err != nil {
		return err
	}

	err = ent.GenEntLogic(params)
	if err != nil {
		return err
	}

	return err
}
