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
	fset := token.NewFileSet()
	for _, filename := range os.Args[1:] {
		if filename == "-" {
			continue
		}
		
		stat, err := os.Stat(filename)
		if err != nil {
			log.Printf("skip %+v", err)
			continue
		}

		if stat.IsDir() {
			if err := runDir(fset, filename); err != nil {
				log.Printf("!! %+v", err)
			}
		} else {
			if err := runFile(fset, filename); err != nil {
				log.Printf("!! %+v", err)
			}
		}
	}
}

func runDir(fset *token.FileSet, dirname string) error {
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
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	}
	return nil
}

func runFile(fset *token.FileSet, filename string) error {
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
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
