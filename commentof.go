package commentof

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/commentof/collect"
)

func Package(fset *token.FileSet, t *ast.Package, options ...Option) (*collect.Package, error) {
	b := defaultBuilder()
	for _, opt := range options {
		opt(b)
	}

	c := &collect.Collector{Fset: fset, Dot: ".", Sharp: "#"}
	p := b.Package
	if err := c.CollectFromPackage(p, t); err != nil {
		return p, err
	}
	return b.Build(), nil
}

func File(fset *token.FileSet, t *ast.File, options ...Option) (*collect.Package, error) {
	b := defaultBuilder()
	for _, opt := range options {
		opt(b)
	}

	c := &collect.Collector{Fset: fset, Dot: ".", Sharp: "#"}
	f := collect.NewFile()
	if err := c.CollectFromFile(f, t); err != nil {
		return nil, err
	}
	b.AddFile(f, fset.File(t.Pos()).Name())
	return b.Build(), nil
}

func defaultBuilder() *collect.PackageBuilder {
	return &collect.PackageBuilder{
		Package:           collect.NewPackage(),
		EnableMergeMethod: true,
		IgnoreExported:    true,
	}
}

type Option func(*collect.PackageBuilder)

func WithIncludeUnexported(ok bool) Option {
	return func(b *collect.PackageBuilder) {
		b.IgnoreExported = !ok
	}
}
