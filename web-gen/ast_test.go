package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestAstPrinter(t *testing.T) {
	fileset := token.NewFileSet()

	createFuncType := func() *ast.FuncType {
		return &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		}
	}

	f := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "Test",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: createFuncType(),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.FuncLit{
							Type: createFuncType(),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "nil",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	printer.Fprint(buf, fileset, f)
	t.Error("\n" + buf.String())
}
func TestAstParser(t *testing.T) {
	fileset := token.NewFileSet()
	source := `
package opt

import (
	"io"
	urlpkg "net/url"

	"github.com/LogiqsAgro/rmq/api/rq"
)

// Interface Request

func Method(method string) func(rq.Request) {
	return func(r rq.Request) { r.Method(method) }
}
`
	file, err := parser.ParseFile(fileset, "", source, parser.AllErrors|parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ast.Print(fileset, file)

}

func TestParseRequestInterface(t *testing.T) {

	wd, err := os.Getwd()
	t.Logf("cwd: %s, %v", wd, err)

	filename := filepath.Join(wd, "../web/requestinterface.go")

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	ast.Fprint(buf, fileset, file, nil)
	fmt.Fprintln(buf)
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		if ident, ok := n.(*ast.Ident); ok {
			pos := fileset.Position(ident.NamePos)

			fmt.Fprintf(buf, "%s %s", pos, ident.Name)
			obj := ident.Obj
			if obj != nil {
				fmt.Fprintf(buf, " kind:%v, name: %s", obj.Kind, obj.Name)
			}
			fmt.Fprintln(buf)

		}

		return true
	})

	t.Error("\n" + buf.String())

}
