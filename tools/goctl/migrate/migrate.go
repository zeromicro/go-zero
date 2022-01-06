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
	"path/filepath"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
	"github.com/urfave/cli"
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

	console.Success("[OK] refactor finish, execute %q on project root to check status.", "go test -race ./...")
	return nil
}

func rewriteImport(verbose bool) error {
	if verbose {
		console.Info("preparing to rewrite import ...")
		time.Sleep(200 * time.Millisecond)
	}

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
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
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

		return rewriteFile(pkgs, verbose)
	})
}

func rewriteFile(pkgs map[string]*ast.Package, verbose bool) error {
	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			for _, imp := range file.Imports {
				if !strings.Contains(imp.Path.Value, deprecatedGoZeroMod) {
					continue
				}
				newPath := strings.ReplaceAll(imp.Path.Value, deprecatedGoZeroMod, goZeroMod)
				imp.EndPos = imp.End()
				imp.Path.Value = newPath
			}

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
				console.Success("[OK] rewriting %q ... ", filepath.Base(filename))
			}
		}
	}
	return nil
}
