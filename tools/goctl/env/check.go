package env

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protoc"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protocgengo"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/protocgengogrpc"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

type bin struct {
	name   string
	exists bool
	get    func(cacheDir string) (string, error)
}

var bins = []bin{
	{
		name:   "protoc",
		exists: protoc.Exists(),
		get:    protoc.Install,
	},
	{
		name:   "protoc-gen-go",
		exists: protocgengo.Exists(),
		get:    protocgengo.Install,
	},
	{
		name:   "protoc-gen-go-grpc",
		exists: protocgengogrpc.Exists(),
		get:    protocgengogrpc.Install,
	},
}

func Check(ctx *cli.Context) error {
	install := ctx.Bool("install")
	force := ctx.Bool("force")
	return Prepare(install, force)
}

func Prepare(install, force bool) error {
	pending := true
	console.Info("[goctl-env]: preparing to check env")
	defer func() {
		if p := recover(); p != nil {
			console.Error("%+v", p)
			return
		}
		if pending {
			console.Success("\n[goctl-env]: congratulations! your goctl environment is ready!")
		} else {
			console.Error(`
[goctl-env]: check env finish, some dependencies is not found in PATH, you can execute
command 'goctl env check --install' to install it, for details, please execute command 
'goctl env check --help'`)
		}
	}()
	for _, e := range bins {
		time.Sleep(200 * time.Millisecond)
		console.Info("")
		console.Info("[goctl-env]: looking up %q", e.name)
		if e.exists {
			console.Success("[goctl-env]: %q is installed", e.name)
			continue
		}
		console.Warning("[goctl-env]: %q is not found in PATH", e.name)
		if install {
			install := func() {
				console.Info("[goctl-env]: preparing to install %q", e.name)
				path, err := e.get(env.Get(env.GoctlCache))
				if err != nil {
					console.Error("[goctl-env]: an error interrupted the installation: %+v", err)
					pending = false
				} else {
					console.Success("[goctl-env]: %q is already installed in %q", e.name, path)
				}
			}
			if force {
				install()
				continue
			}
			console.Info("[goctl-env]: do you want to install %q [y: YES, n: No]", e.name)
			for {
				var in string
				fmt.Scanln(&in)
				var brk bool
				switch {
				case strings.EqualFold(in, "y"):
					install()
					brk = true
				case strings.EqualFold(in, "n"):
					pending = false
					console.Info("[goctl-env]: %q installation is ignored", e.name)
					brk = true
				default:
					console.Error("[goctl-env]: invalid input, input 'y' for yes, 'n' for no")
				}
				if brk {
					break
				}
			}
		} else {
			pending = false
		}
	}
	return nil
}
