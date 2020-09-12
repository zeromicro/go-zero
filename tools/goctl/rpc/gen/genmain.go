package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const mainTemplate = `{{.head}}

package main

import (
	"flag"
	"fmt"
	"log"

	{{.imports}}

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rpcx"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	{{.srv}}

	s, err := rpcx.NewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{.registers}}
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
`

func (g *defaultRpcGenerator) genMain() error {
	mainPath := g.dirM[dirTarget]
	file := g.ast
	pkg := file.Package

	fileName := filepath.Join(mainPath, fmt.Sprintf("%v.go", g.Ctx.ServiceName.Lower()))
	imports := make([]string, 0)
	pbImport := fmt.Sprintf(`%v "%v"`, pkg, g.mustGetPackage(dirPb))
	svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
	remoteImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirServer))
	configImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirConfig))
	imports = append(imports, configImport, pbImport, remoteImport, svcImport)
	srv, registers := g.genServer(pkg, file.Service)
	head := util.GetHead(g.Ctx.ProtoSource)
	return util.With("main").GoFmt(true).Parse(mainTemplate).SaveTo(map[string]interface{}{
		"head":        head,
		"package":     pkg,
		"serviceName": g.Ctx.ServiceName.Lower(),
		"srv":         srv,
		"registers":   registers,
		"imports":     strings.Join(imports, "\n"),
	}, fileName, true)
}

func (g *defaultRpcGenerator) genServer(pkg string, list []*parser.RpcService) (string, string) {
	list1 := make([]string, 0)
	list2 := make([]string, 0)
	for _, item := range list {
		name := item.Name.UnTitle()
		list1 = append(list1, fmt.Sprintf("%sSrv := server.New%sServer(ctx)", name, item.Name.Title()))
		list2 = append(list2, fmt.Sprintf("%s.Register%sServer(grpcServer, %sSrv)", pkg, item.Name.Title(), name))
	}
	return strings.Join(list1, "\n"), strings.Join(list2, "\n")
}
