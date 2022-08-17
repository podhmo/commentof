package commentof

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/commentof/collect"
)

func File(fset *token.FileSet, t *ast.File) (*collect.File, error) {
	c := &collect.Collector{Fset: fset, Dot: "."}
	f := &collect.File{Structs: map[string]*collect.Struct{}}
	return f, c.CollectFromFile(f, t)
}
