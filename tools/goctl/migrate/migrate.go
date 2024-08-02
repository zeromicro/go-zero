package migrate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
)

const defaultMigrateVersion = "v1.3.0"

const (
	confirmUnknown = iota
	confirmAll
	confirmIgnore
)

var (
	fset            = token.NewFileSet()
	builderxConfirm = confirmUnknown
)

func migrate(_ *cobra.Command, _ []string) error {
	if len(stringVarVersion) == 0 {
		stringVarVersion = defaultMigrateVersion
	}
	err := editMod(stringVarVersion, boolVarVerbose)
	if err != nil {
		return err
	}

	err = rewriteImport(boolVarVerbose)
	if err != nil {
		return err
	}

	err = tidy(boolVarVerbose)
	if err != nil {
		return err
	}

	if boolVarVerbose {
		console.Success("[OK] refactor finish, execute %q on project root to check status.",
			"go test -race ./...")
	}

	return nil
}

func rewriteImport(verbose bool) error {
	if verbose {
		console.Info("preparing to rewrite import ...")
		time.Sleep(200 * time.Millisecond)
	}

	cancelOnSignals()

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
					console.Debug("[...] migrating %q ... ", filepath.Base(filename))
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
			w := bytes.NewBuffer(nil)
			err := format.Node(w, fset, file)
			if err != nil {
				return fmt.Errorf("[rewriteImport] format file %s error: %w", filename, err)
			}

			err = os.WriteFile(filename, w.Bytes(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("[rewriteImport] write file %s error: %w", filename, err)
			}
			if verbose {
				console.Success("[OK] migrated %q successfully", filepath.Base(filename))
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
	switch builderxConfirm {
	case confirmAll:
		fn(true)
		return
	case confirmIgnore:
		fn(false)
		return
	}

	msg := fmt.Sprintf(`Detects a deprecated package in the source code,
Deprecated package: %q
Replacement package: %q
It's recommended to use the replacement package, do you want to replace?
['Y' for yes, 'N' for no, 'A' for all, 'I' for ignore]: `,
		deprecated, replacement)

	fmt.Print(color.Yellow.Render(msg))

	for {
		var in string
		fmt.Scanln(&in)
		switch {
		case strings.EqualFold(in, "Y"):
			fn(true)
			return
		case strings.EqualFold(in, "N"):
			fn(false)
			return
		case strings.EqualFold(in, "A"):
			fn(true)
			builderxConfirm = confirmAll
			return
		case strings.EqualFold(in, "I"):
			fn(false)
			builderxConfirm = confirmIgnore
			return
		default:
			console.Warning("['Y' for yes, 'N' for no, 'A' for all, 'I' for ignore]: ")
		}
	}
}
