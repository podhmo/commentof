package commentof

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/commentof/collect"
)

func File(fset *token.FileSet, t *ast.File) (*collect.File, error) {
	c := &collect.Collector{Fset: fset, Dot: "."}
	f := &collect.File{
		Structs:    map[string]*collect.Object{},
		Interfaces: map[string]*collect.Object{},
		Functions:  map[string]*collect.Func{},
		Names:      []string{},
	}
	return f, c.CollectFromFile(f, t)
}

func Package(fset *token.FileSet, t *ast.Package) (*collect.Package, error) {
	c := &collect.Collector{Fset: fset, Dot: "."}
	p := &collect.Package{
		Files: map[string]*collect.File{}, FileNames: []string{},
		Structs:    map[string]*collect.Object{},
		Interfaces: map[string]*collect.Object{},
		Functions:  map[string]*collect.Func{},
		Names:      []string{},
	}
	return p, c.CollectFromPackage(p, t)
}
