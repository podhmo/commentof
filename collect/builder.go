package collect

import "strings"

type PackageBuilder struct {
	Package *Package

	EnableMergeMethod bool
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
