package gogen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func updateHandlerComments(filename string, handlerName string, comments []string) error {
	// Parse the code file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Find the UserSelfHandler function
	var fn *ast.FuncDecl
	for i := range f.Decls {
		if fd, ok := f.Decls[i].(*ast.FuncDecl); ok && fd.Name.Name == handlerName {
			fn = fd
			break
		}
	}
	if fn == nil {
		return fmt.Errorf("handler %s not found", handlerName)
	}

	// Update the function's comments
	var list []*ast.Comment
	for i, comment := range comments {
		if i == 0 {
			list = append(list, &ast.Comment{Slash: fn.Pos() - 1, Text: comment})
		} else {
			list = append(list, &ast.Comment{Text: comment})
		}
	}
	if fn.Doc == nil {
		fn.Doc = &ast.CommentGroup{}
	}
	fn.Doc.List = list

	// Format and write the updated code to file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
