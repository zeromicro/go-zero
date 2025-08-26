package env

import (
	"fmt"
	"strings"
	"time"

	"github.com/lerity-yao/go-zero/tools/cztctl/pkg/env"
	"github.com/lerity-yao/go-zero/tools/cztctl/pkg/protoc"
	"github.com/lerity-yao/go-zero/tools/cztctl/pkg/protocgengo"
	"github.com/lerity-yao/go-zero/tools/cztctl/pkg/protocgengogrpc"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/console"
	"github.com/spf13/cobra"
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

func check(_ *cobra.Command, _ []string) error {
	return Prepare(boolVarInstall, boolVarForce, boolVarVerbose)
}

func Prepare(install, force, verbose bool) error {
	log := console.NewColorConsole(verbose)
	pending := true
	log.Info("[cztctl-env]: preparing to check env")
	defer func() {
		if p := recover(); p != nil {
			log.Error("%+v", p)
			return
		}
		if pending {
			log.Success("\n[cztctl-env]: congratulations! your cztctl environment is ready!")
		} else {
			log.Error(`
[cztctl-env]: check env finish, some dependencies is not found in PATH, you can execute
command 'cztctl env check --install' to install it, for details, please execute command 
'cztctl env check --help'`)
		}
	}()
	for _, e := range bins {
		time.Sleep(200 * time.Millisecond)
		log.Info("")
		log.Info("[cztctl-env]: looking up %q", e.name)
		if e.exists {
			log.Success("[cztctl-env]: %q is installed", e.name)
			continue
		}
		log.Warning("[cztctl-env]: %q is not found in PATH", e.name)
		if install {
			install := func() {
				log.Info("[cztctl-env]: preparing to install %q", e.name)
				path, err := e.get(env.Get(env.GoctlCache))
				if err != nil {
					log.Error("[cztctl-env]: an error interrupted the installation: %+v", err)
					pending = false
				} else {
					log.Success("[cztctl-env]: %q is already installed in %q", e.name, path)
				}
			}
			if force {
				install()
				continue
			}
			console.Info("[cztctl-env]: do you want to install %q [y: YES, n: No]", e.name)
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
					console.Info("[cztctl-env]: %q installation is ignored", e.name)
					brk = true
				default:
					console.Error("[cztctl-env]: invalid input, input 'y' for yes, 'n' for no")
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
