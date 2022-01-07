/*
Copyright © 2021 Remco Schoeman <remco.schoeman@logiqs.nl>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	run()
}

func run() {
	v := newGenerator()
	v.parseInputFile()
	v.generateOutputFile()

	log.Print(v.args.String())
}

type (
	generator struct {
		fileset *token.FileSet
		file    *ast.File

		interfaceName string
		args          args

		out io.WriteCloser
	}
)

func newGenerator() *generator {
	v := &generator{

		fileset: token.NewFileSet(),
	}

	return v
}
func (v *generator) parseInputFile() error {
	cwd := v.cwd()
	pkgs, err := parser.ParseDir(v.fileset, cwd, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing files: %w", err)
	}

	gofile := v.goFile()
	for pkgName := range pkgs {
		pkg := pkgs[pkgName]

		// we make a new package, so that all the unresolved identifiers are resolved, hopefully.
		ast.NewScope(nil)

		pkg, _ = ast.NewPackage(v.fileset, pkg.Files, nil, nil)

		// ignore the error, we only want to resolve identifiers in
		// the package itself, not of the imported packages
		pkgs[pkgName] = pkg
		for fileName := range pkg.Files {
			f := pkg.Files[fileName]
			if fileName == gofile {
				v.file = f
				return nil
			}
		}
	}
	return fmt.Errorf("could not find parse result for file '%s' in dir '%s'", v.Arg("GOGILE"), cwd)
}

func (v *generator) generateOutputFile() error {
	if v.file == nil {
		return fmt.Errorf("inputfile not parsed yet")
	}

	ast.Print(v.fileset, v.file)

	if f, err := os.Create(v.outputFilePath()); err != nil {
		return err
	} else {
		v.out = f
		v.printf("package %s\n\n", v.outputPackage())
		ast.Walk(v, v.file)
		return nil
	}
}

func (v *generator) cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory: %w", err))
	}
	return cwd
}

// goFile returns the full path of the file in the GOFILE environment variable
func (v *generator) goFile() string {
	return filepath.Join(v.cwd(), v.Arg("GOFILE"))
}

func (v *generator) outputFilePath() string {
	pkg := v.outputPackage()
	return filepath.Join(v.outputDir(), pkg+".go")
}

func (v *generator) outputDir() string {
	return filepath.Join(v.cwd(), v.outputPackage())
}

func (v *generator) outputPackage() string {
	return v.Arg("package")
}

func (v *generator) Arg(key string) string {
	return v.Args()[key]
}

func (v *generator) Args() args {
	if v.args != nil {
		return v.args
	}
	aa := make(args)

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		aa[kv[0]] = kv[1]
	}

	for _, arg := range os.Args[1:] {
		kv := strings.SplitN(arg, "=", 2)
		aa[kv[0]] = kv[1]
	}

	v.args = aa
	return aa
}

func (v *generator) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if imports, ok := n.(*ast.GenDecl); ok {
		if imports.Tok.String() == "import" {
			newImports := &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: append(imports.Specs, &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"github.com/LogiqsAgro/rmq/web\"",
					},
				}),
			}
			v.printNode(newImports)
			v.println("\n")

		}
	}

	if typeSpec, ok := n.(*ast.TypeSpec); ok {
		if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			v.interfaceName = typeSpec.Name.String()
			methods := iface.Methods.List
			for i := 0; i < len(methods); i++ {

				method := methods[i]
				name := method.Names[0].Name
				if !ast.IsExported(name) {
					continue
				}

				pv := newMethodVisitor(v)
				ast.Walk(pv, method)
				comment := strings.TrimSuffix(method.Doc.Text(), "\n")
				v.println("// ", strings.ReplaceAll(comment, "\n", "\n// "))
				v.print("func ")
				v.printIdent(method.Names...)
				v.print("(")
				pv.writeParameterDeclarations()
				v.printf(") func(%s.%s) {\n", v.file.Name.Name, v.interfaceName)

				v.printf("	return func(r %s.%s) { r.%s(", v.file.Name.Name, v.interfaceName, name)
				pv.writeArgumentList()
				v.printf(") }\n")
				v.print("}\n\n")
			}
		}
	}

	return v
}

func (v *generator) printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(v.out, format, args...)
}

func (v *generator) print(s ...string) (int, error) {
	nn := 0
	for i := 0; i < len(s); i++ {
		n, err := fmt.Fprint(v.out, s[i])
		nn += n
		if err != nil {
			return nn, err
		}
	}
	return nn, nil
}

func (v *generator) println(s ...string) (int, error) {
	n, err := v.print(s...)
	if err != nil {
		return n, err
	}
	n2, err := v.print("\n")
	return n + n2, err
}

func (v *generator) printNode(node interface{}) {
	if id, ok := isTypeIdent(node); ok {
		v.printNode(&ast.SelectorExpr{
			X:   &ast.Ident{Name: v.file.Name.Name},
			Sel: &ast.Ident{Name: id.Name},
		})
		return
	}
	panicIf(printer.Fprint(v.out, v.fileset, node))
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func isTypeIdent(node interface{}) (*ast.Ident, bool) {
	if id, ok := node.(*ast.Ident); ok {
		if id.Obj == nil {
			return nil, false
		} else if id.Obj.Kind == ast.Typ {
			return id, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (v *generator) printIdent(names ...*ast.Ident) {
	for i := 0; i < len(names); i++ {
		if i > 0 {
			v.out.Write([]byte(", "))
		}
		name := names[i]
		v.printNode(name)
	}
}

type (
	methodVisitor struct {
		parameters *ast.FieldList
		v          *generator
	}
)

func newMethodVisitor(v *generator) *methodVisitor {
	return &methodVisitor{v: v}
}

func (v *methodVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if f, ok := n.(*ast.FuncType); ok {
		v.parameters = f.Params
		return nil
	}

	return v
}

func (v *methodVisitor) writeParameterDeclarations() {
	pp := v.parameters.List
	for i := 0; i < len(pp); i++ {
		if i > 0 {
			v.v.print(", ")
		}
		p := pp[i]
		t := p.Type
		v.v.printIdent(p.Names...)
		v.v.print(" ")
		if e, ok := t.(*ast.Ellipsis); ok {
			v.v.print("...")
			v.v.printNode(e.Elt)
		} else {
			v.v.printNode(t)
		}
	}
}

func (v *methodVisitor) writeArgumentList() {
	pp := v.parameters.List

	for i := 0; i < len(pp); i++ {
		if i > 0 {
			v.v.print(", ")
		}

		v.v.printIdent(pp[i].Names...)

		switch pp[i].Type.(type) {
		case *ast.Ellipsis:
			v.v.print("...")
		default:
			// nothing
		}
	}

}

type args map[string]string

func (a args) String() string {
	keys := make([]string, 0, len(a))
	for key := range a {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buf := &bytes.Buffer{}
	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(a[key])
		buf.WriteString("\n")
	}
	return buf.String()
}
