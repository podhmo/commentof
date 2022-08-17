package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Printf("!! %+v", err)
	}
}

func run() error {
	fset := token.NewFileSet()
	filename := "./testdata/fixture/fixture.go"

	t, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	return Parse(fset, t)
}

func Parse(fset *token.FileSet, t *ast.File) error {
	for _, cg := range t.Comments {
		fmt.Println(strings.TrimSpace(cg.Text()))
	}
	fmt.Println("----------------------------------------")
	for _, decl := range t.Decls {
		switch decl.(type) {
		case *ast.FuncDecl, *ast.BadDecl:
		case *ast.GenDecl:
			fmt.Println("decl", decl)
		default:
			log.Printf("unexpected decl: %T", decl)
			continue
		}
	}
	return nil
}
