package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"

	"github.com/k0kubun/pp"
	"github.com/podhmo/commentof"
)

func main() {
	if err := run(); err != nil {
		log.Printf("!! %+v", err)
	}
}

func run() error {
	fset := token.NewFileSet()
	filename := "./testdata/fixture/struct.go"
	// filename := "./testdata/fixture/const.go"
	// filename := "./testdata/fixture/embedded.go"

	tree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	f, err := commentof.FileAST(fset, tree)
	if err != nil {
		return fmt.Errorf("collect: file=%s, %w", filename, err)
	}

	pp.Println(f)
	return nil
}
