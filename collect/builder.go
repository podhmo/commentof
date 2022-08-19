package collect

import (
	"go/ast"
	"strings"
)

type PackageBuilder struct {
	Package *Package

	EnableMergeMethod bool
	IgnoreExported    bool
}

func (b *PackageBuilder) AddFile(f *File, filename string) {
	p := b.Package

	p.FileNames = append(p.FileNames, filename)
	p.Files[filename] = f

	p.Names = append(p.Names, f.Names...)
	for id, s := range f.Structs {
		p.Structs[id] = s
	}
	for id, s := range f.Interfaces {
		p.Interfaces[id] = s
	}
	for id, s := range f.Functions {
		p.Functions[id] = s
	}
}

func (b *PackageBuilder) Build() *Package {
	if b.EnableMergeMethod {
		mergeMethod(b.Package)
	}
	if b.IgnoreExported {
		ignoreExported(b.Package)
	}
	return b.Package
}

func mergeMethod(p *Package) {
	names := make([]string, 0, len(p.Names))
	for _, name := range p.Names {
		if !strings.Contains(name, "#") {
			names = append(names, name)
			continue
		}
		method := p.Functions[name]
		obName := strings.TrimPrefix(method.Recv, "*")
		if ob, ok := p.Structs[obName]; ok {
			ob.MethodNames = append(ob.MethodNames, method.Name)
			ob.Methods[method.Name] = method
			delete(p.Functions, name)
		}
	}
	p.Names = names
}

func ignoreExported(p *Package) {
	names := make([]string, 0, len(p.Names))
	for _, name := range p.Names {
		if ast.IsExported(name) {
			names = append(names, name)
			continue
		}

		if _, ok := p.Functions[name]; ok {
			delete(p.Functions, name)
			continue
		} else if _, ok := p.Structs[name]; ok {
			delete(p.Structs, name)
			continue
		} else if _, ok := p.Interfaces[name]; ok {
			delete(p.Interfaces, name)
			continue
		}
	}
	p.Names = names
	for _, f := range p.Files {
		ignoreExportedForFile(f)
	}
}

func ignoreExportedForFile(f *File) {
	names := make([]string, 0, len(f.Names))
	for _, name := range f.Names {
		if ast.IsExported(name) {
			names = append(names, name)
			if ob, ok := f.Structs[name]; ok {
				ignoreExportedForObject(ob)
			} else if ob, ok := f.Interfaces[name]; ok {
				ignoreExportedForObject(ob)
			} else {
				_ = false // function
			}
			continue
		}

		if _, ok := f.Structs[name]; ok {
			delete(f.Structs, name)
			continue
		} else if _, ok := f.Interfaces[name]; ok {
			delete(f.Interfaces, name)
			continue
		} else if _, ok := f.Functions[name]; ok {
			delete(f.Functions, name)
			continue
		}
	}
	f.Names = names
}

func ignoreExportedForObject(ob *Object) {
	if len(ob.Fields) > 0 {
		names := make([]string, 0, len(ob.FieldNames))
		for _, name := range ob.FieldNames {
			if ast.IsExported(name) {
				names = append(names, name)
				continue
			}
			delete(ob.Fields, name)
		}
		ob.FieldNames = names
	}
	if len(ob.Methods) > 0 {
		names := make([]string, 0, len(ob.MethodNames))
		for _, name := range ob.MethodNames {
			if ast.IsExported(name) {
				names = append(names, name)
				continue
			}
			delete(ob.Methods, name)
		}
		ob.MethodNames = names
	}
}
