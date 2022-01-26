package migrate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

const zeromicroVersion = "1.3.0"

var fset = token.NewFileSet()

func Migrate(c *cli.Context) error {
	verbose := c.Bool("verbose")
	version := c.String("version")
	if len(version) == 0 {
		version = zeromicroVersion
	}
	err := editMod(version, verbose)
	if err != nil {
		return err
	}

	err = rewriteImport(verbose)
	if err != nil {
		return err
	}

	err = tidy(verbose)
	if err != nil {
		return err
	}

	if verbose {
		console.Success("[OK] refactor finish, execute %q on project root to check status.", "go test -race ./...")
	}
	return nil
}

func rewriteImport(verbose bool) error {
	if verbose {
		console.Info("preparing to rewrite import ...")
		time.Sleep(200 * time.Millisecond)
	}

	var doneChan = syncx.NewDoneChan()
	defer func() {
		doneChan.Close()
	}()
	go func(dc *syncx.DoneChan) {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
		select {
		case <-c:
			console.Error(`
migrate failed, reason: "User Canceled"`)
			os.Exit(0)
		case <-dc.Done():
			return
		}
	}(doneChan)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	project, err := ctx.Prepare(wd)
	if err != nil {
		return err
	}
	root := project.Dir
	fsys := os.DirFS(root)
	var final []*ast.Package
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		if verbose {
			console.Info("walking to %q", path)
		}
		pkgs, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
			return strings.HasSuffix(info.Name(), ".go")
		}, parser.ParseComments)
		if err != nil {
			return err
		}

		err = rewriteFile(pkgs, verbose)
		if err != nil {
			return err
		}
		for _, v := range pkgs {
			final = append(final, v)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if verbose {
		console.Info("start to write files ... ")
	}
	return writeFile(final, verbose)
}

func rewriteFile(pkgs map[string]*ast.Package, verbose bool) error {
	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			var containsDeprecatedBuilderxPkg bool
			for _, imp := range file.Imports {
				if !strings.Contains(imp.Path.Value, deprecatedGoZeroMod) {
					continue
				}

				if verbose {
					console.Debug("[...] migrate %q ... ", filepath.Base(filename))
				}

				if strings.Contains(imp.Path.Value, deprecatedBuilderx) {
					containsDeprecatedBuilderxPkg = true
					var doNext bool
					refactorBuilderx(deprecatedBuilderx, replacementBuilderx, func(allow bool) {
						doNext = !allow
						if allow {
							newPath := strings.ReplaceAll(imp.Path.Value, deprecatedBuilderx, replacementBuilderx)
							imp.EndPos = imp.End()
							imp.Path.Value = newPath
						}
					})
					if !doNext {
						continue
					}
				}

				newPath := strings.ReplaceAll(imp.Path.Value, deprecatedGoZeroMod, goZeroMod)
				imp.EndPos = imp.End()
				imp.Path.Value = newPath
			}

			if containsDeprecatedBuilderxPkg {
				replacePkg(file)
			}
		}
	}
	return nil
}

func writeFile(pkgs []*ast.Package, verbose bool) error {
	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			var w = bytes.NewBuffer(nil)
			err := format.Node(w, fset, file)
			if err != nil {
				return fmt.Errorf("[rewriteImport] format file %s error: %+v", filename, err)
			}

			err = ioutil.WriteFile(filename, w.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("[rewriteImport] write file %s error: %+v", filename, err)
			}
			if verbose {
				console.Success("[OK] migrate %q success ", filepath.Base(filename))
			}
		}
	}
	return nil
}

func replacePkg(file *ast.File) {
	scope := file.Scope
	if scope == nil {
		return
	}
	obj := scope.Objects
	for _, v := range obj {
		decl := v.Decl
		if decl == nil {
			continue
		}
		vs, ok := decl.(*ast.ValueSpec)
		if !ok {
			continue
		}
		values := vs.Values
		if len(values) != 1 {
			continue
		}
		value := values[0]
		callExpr, ok := value.(*ast.CallExpr)
		if !ok {
			continue
		}
		fn := callExpr.Fun
		if fn == nil {
			continue
		}
		selector, ok := fn.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		x := selector.X
		sel := selector.Sel
		if x == nil || sel == nil {
			continue
		}
		ident, ok := x.(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Name == "builderx" {
			ident.Name = "builder"
			ident.NamePos = ident.End()
		}
		if sel.Name == "FieldNames" {
			sel.Name = "RawFieldNames"
			sel.NamePos = sel.End()
		}
	}
}

func refactorBuilderx(deprecated, replacement string, fn func(allow bool)) {
	msg := fmt.Sprintf(`Detects a deprecated package in the source code,
Deprecated package: %q
Replacement package: %q
It's recommended to use the replacement package, do you want to replace?
[input 'Y' for yes, 'N' for no]:`, deprecated, replacement)

	if runtime.GOOS != vars.OsWindows {
		msg = aurora.Yellow(msg).String()
	}
	fmt.Print(msg)
	var in string
	for {
		fmt.Scanln(&in)
		if len(in) == 0 {
			console.Warning("nothing input, please try again [input 'Y' for yes, 'N' for no]:")
			continue
		}
		if strings.EqualFold(in, "Y") {
			fn(true)
			return
		} else if strings.EqualFold(in, "N") {
			fn(false)
			return
		} else {
			console.Warning("invalid input, please try again [input 'Y' for yes, 'N' for no]:")
		}
	}
}
