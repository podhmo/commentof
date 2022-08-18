package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

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
			stdSrcFilename := filepath.Join(runtime.GOROOT(), "src", filename)
			stat, err = os.Stat(stdSrcFilename)
			if err != nil {
				log.Printf("skip %+v", err)
				continue
			}
			filename = stdSrcFilename
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

	names := make([]string, 0, len(tree))
	for name := range tree {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		p := tree[name]
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
