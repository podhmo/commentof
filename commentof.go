package commentof

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/commentof/collect"
)

func Package(fset *token.FileSet, t *ast.Package) (*collect.Package, error) {
	c := &collect.Collector{Fset: fset, Dot: ".", Sharp: "#"}
	p := collect.NewPackage()
	if err := c.CollectFromPackage(p, t); err != nil {
		return p, err
	}
	b := &collect.PackageBuilder{
		Package:           p,
		EnableMergeMethod: true,
	}
	return b.Build(), nil
}

func File(fset *token.FileSet, t *ast.File) (*collect.Package, error) {
	c := &collect.Collector{Fset: fset, Dot: ".", Sharp: "#"}
	f := collect.NewFile()
	if err := c.CollectFromFile(f, t); err != nil {
		return nil, err
	}
	b := &collect.PackageBuilder{
		Package:           collect.NewPackage(),
		EnableMergeMethod: true,
	}
	b.AddFile(f, fset.File(t.Pos()).Name())
	return b.Build(), nil
}
