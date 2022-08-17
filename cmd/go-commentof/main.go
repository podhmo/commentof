package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/podhmo/commentof"
)

func main() {
	if err := run(); err != nil {
		log.Printf("!! %+v", err)
	}
}

func run() error {
	fset := token.NewFileSet()
	dirname := "./testdata/fixture"
	tree, err := parser.ParseDir(fset, dirname, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse dir: %w", err)
	}

	for name, p := range tree {
		result, err := commentof.Package(fset, p)
		if err != nil {
			return fmt.Errorf("collect: dir=%s, name=%s, %w", dirname, name, err)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "	")
		fmt.Println(enc.Encode(result))
	}
	return nil
}

func RunFile() error {
	fset := token.NewFileSet()
	// filename := "./testdata/fixture/struct.go"
	// filename := "./testdata/fixture/const.go"
	filename := "./testdata/fixture/embedded.go"

	tree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	result, err := commentof.File(fset, tree)
	if err != nil {
		return fmt.Errorf("collect: file=%s, %w", filename, err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "	")
	fmt.Println(enc.Encode(result))

	// pp.Println(result)
	return nil
}
