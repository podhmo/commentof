package collect

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"sort"
)

type Collector struct {
	Fset  *token.FileSet
	Dot   string
	Sharp string
}

func (c *Collector) CollectFromPackage(p *Package, t *ast.Package) error {
	b := &PackageBuilder{Package: p}

	filenames := make([]string, 0, len(t.Files))
	for filename := range t.Files {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		ft := t.Files[filename]
		f := NewFile()
		if err := c.CollectFromFile(f, ft); err != nil {
			return fmt.Errorf("collect file: %s: %w", filename, err)
		}
		b.AddFile(f, filename)
	}
	return nil
}

func (c *Collector) CollectFromFile(f *File, t *ast.File) error {
	for _, decl := range t.Decls {
		switch decl := decl.(type) {
		case *ast.BadDecl:
		case *ast.FuncDecl:
			if err := c.CollectFromFuncDecl(f, t, decl); err != nil {
				return err
			}
		case *ast.GenDecl:
			if err := c.CollectFromGenDecl(f, decl); err != nil {
				return err
			}
		default:
			log.Printf("unexpected decl: %T?", decl)
			continue
		}
	}
	return nil
}

func (c *Collector) CollectFromFuncDecl(f *File, t *ast.File, decl *ast.FuncDecl) error {
	recv := ""
	if decl.Recv != nil && decl.Recv.List != nil {
		if v, ok := typeString(decl.Recv.List[0].Type); ok {
			recv = v
		}
	}

	name := decl.Name.Name
	id := name
	if recv != "" {
		id = recv + c.Sharp + name
	}

	f.Names = append(f.Names, id)

	paramNames := []string{}
	params := map[string]*Field{}
	var comments []*ast.CommentGroup
	idx := 0
	{
		start := decl.Type.Params.Opening
		end := decl.Type.Params.Closing
		for i, cg := range t.Comments {
			if cg.End() < start {
				continue
			}
			if end < cg.Pos() {
				idx = i
				break
			}
			comments = append(comments, cg)
		}
	}

	for i, x := range decl.Type.Params.List {
		name := ""
		id := ""
		if len(x.Names) > 0 {
			name = x.Names[0].Name
			id = name
		} else {
			id = fmt.Sprintf("param#%d", i)
		}

		doc := ""
		for _, cg := range comments {
			// fmt.Println(f.Names[len(f.Names)-1], id, "@@", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
			if x.Pos() < cg.Pos() && cg.End() < x.End() {
				doc += cg.Text()
				// fmt.Println(f.Names[len(f.Names)-1], id, "-#", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
				continue
			}
			if x.End() < cg.Pos() {
				// fmt.Println(f.Names[len(f.Names)-1], id, "--", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
				doc += cg.Text()
				break
			}
		}

		paramNames = append(paramNames, id)
		params[id] = &Field{
			Name:    name,
			Comment: doc,
		}
	}

	returnNames := []string{}
	returns := map[string]*Field{}
	comments = nil
	if decl.Type.Results != nil {
		{
			start := decl.Type.Results.Opening
			end := decl.Type.Results.Closing
			for _, cg := range t.Comments[idx:] {
				if cg.End() < start {
					continue
				}
				if end < cg.Pos() {
					break
				}
				comments = append(comments, cg)
			}
		}

		for i, x := range decl.Type.Results.List {
			name := ""
			id := ""
			if len(x.Names) > 0 {
				name = x.Names[0].Name
				id = name
			} else {
				id = fmt.Sprintf("ret#%d", i)
			}

			returnNames = append(returnNames, id)
			doc := ""
			for _, cg := range comments {
				// fmt.Println(f.Names[len(f.Names)-1], id, "@@", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
				if x.Pos() < cg.Pos() && cg.End() < x.End() {
					doc += cg.Text()
					// fmt.Println(f.Names[len(f.Names)-1], id, "-#", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
					continue
				}
				if x.End() < cg.Pos() {
					// fmt.Println(f.Names[len(f.Names)-1], id, "--", x.Pos(), x.End(), "@", cg.Pos(), cg.End(), "--", strings.TrimSpace(cg.Text()))
					doc += cg.Text()
					break
				}
			}
			returns[id] = &Field{
				Name:    name,
				Comment: doc,
			}
		}
	}
	f.Functions[id] = &Func{
		Name:        name,
		Recv:        recv,
		Doc:         decl.Doc.Text(),
		Params:      params,
		ParamNames:  paramNames,
		Returns:     returns,
		ReturnNames: returnNames,
	}
	return nil
}

func (c *Collector) CollectFromGenDecl(f *File, decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.ImportSpec, *ast.ValueSpec:
		case *ast.TypeSpec:
			if err := c.CollectFromTypeSpec(f, decl, spec); err != nil {
				return err
			}
		default:
			log.Printf("unexpected decl: %T, spec: %T?", decl, spec)
			continue
		}
	}
	return nil
}

func (c *Collector) CollectFromTypeSpec(f *File, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	name := spec.Name.Name
	f.Names = append(f.Names, name)
	s := &Object{
		Name:       name,
		Doc:        spec.Doc.Text(),
		Comment:    spec.Comment.Text(),
		FieldNames: []string{},
		Fields:     map[string]*Field{},
		Methods:    map[string]*Func{},
	}
	if s.Doc == "" && decl.Doc != nil {
		s.Doc = decl.Doc.Text()
	}

	switch typ := spec.Type.(type) {
	case *ast.Ident:
		// type <S> <S>
		// type <S> = <S>
		f.Types[name] = s
	case *ast.StructType:
		// type <S> struct { ... }
		f.Types[name] = s
		if err := c.CollectFromStructType(f, s, decl, spec, typ); err != nil {
			return err
		}
	case *ast.InterfaceType:
		// type <S> interface { ... }
		f.Interfaces[name] = s
		if err := c.CollectFromInterfaceType(f, s, decl, spec, typ); err != nil {
			return err
		}
	case *ast.FuncType:
		f.Types[name] = s
	default:
		log.Printf("unexpected decl: %T, spec: %T, type: %T?", decl, spec, typ)
	}

	return nil
}

func typeString(typ ast.Expr) (string, bool) {
	switch t := typ.(type) {
	case *ast.Ident:
		return t.String(), true
	case *ast.SelectorExpr:
		name, ok := typeString(t.X)
		return name + "." + t.Sel.String(), ok
	case *ast.InterfaceType, *ast.StructType:
		return "", true
	case *ast.StarExpr:
		name, ok := typeString(t.X)
		return "*" + name, ok
	default:
		return "", false
	}
}

func (c *Collector) CollectFromStructType(f *File, s *Object, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.StructType) error {
	s.Token = token.STRUCT
	for i, field := range typ.Fields.List {
		name := ""
		anonymous := false
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else {
			anonymous = true
			if typename, ok := typeString(field.Type); ok {
				name = typename
			} else {
				name = fmt.Sprintf("??%T", field.Type) // TODO: NG
				log.Printf("unexpected embedded field type: %T, spec: %T, struct: %T, field:%v", decl, spec, typ, field.Type)
			}
		}
		id := name
		if id == "" {
			id = fmt.Sprintf("anon#%d", i)
		}

		s.FieldNames = append(s.FieldNames, id)
		fieldof := &Field{
			Name:     name,
			Doc:      field.Doc.Text(),
			Comment:  field.Comment.Text(),
			Embedded: anonymous,
		}
		s.Fields[id] = fieldof

		switch typ := field.Type.(type) {
		case *ast.Ident, *ast.FuncType, *ast.SelectorExpr:
		case *ast.StructType:
			// struct { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),
				Comment:    field.Comment.Text(),
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromStructType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		case *ast.InterfaceType:
			// interface { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),
				Comment:    field.Comment.Text(),
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromInterfaceType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		case *ast.BadExpr, *ast.Ellipsis, *ast.BasicLit, *ast.FuncLit, *ast.CompositeLit,
			*ast.ParenExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr, *ast.CallExpr,
			*ast.StarExpr, *ast.UnaryExpr, *ast.BinaryExpr, *ast.KeyValueExpr,
			*ast.ArrayType, *ast.MapType, *ast.ChanType:
		default:
			log.Printf("unexpected decl: %T, spec: %T, type: %T?, field=%s", decl, spec, typ, name)
		}
	}
	return nil
}

func (c *Collector) CollectFromInterfaceType(f *File, s *Object, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.InterfaceType) error {
	s.Token = token.INTERFACE
	for i, field := range typ.Methods.List {
		name := ""
		anonymous := false
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else {
			anonymous = true
			if typename, ok := typeString(field.Type); ok {
				name = typename
			} else {
				name = fmt.Sprintf("??%T", field.Type) // TODO: NG
				log.Printf("unexpected embedded field type: %T, spec: %T, struct: %T, field:%v", decl, spec, typ, field.Type)
			}
		}
		id := name
		if id == "" {
			id = fmt.Sprintf("anon#%d", i)
		}

		s.FieldNames = append(s.FieldNames, id)
		fieldof := &Field{
			Name:     name,
			Doc:      field.Doc.Text(),
			Comment:  field.Comment.Text(),
			Embedded: anonymous,
		}
		s.Fields[id] = fieldof

		switch typ := field.Type.(type) {
		case *ast.Ident, *ast.FuncType, *ast.SelectorExpr:
		case *ast.InterfaceType:
			// interface { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),
				Comment:    field.Comment.Text(),
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromInterfaceType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		case *ast.BadExpr, *ast.Ellipsis, *ast.BasicLit, *ast.FuncLit, *ast.CompositeLit,
			*ast.ParenExpr, *ast.IndexExpr, *ast.IndexListExpr, *ast.SliceExpr, *ast.TypeAssertExpr, *ast.CallExpr,
			*ast.StarExpr, *ast.UnaryExpr, *ast.BinaryExpr, *ast.KeyValueExpr,
			*ast.ArrayType, *ast.MapType, *ast.ChanType:
		default:
			log.Printf("unexpected decl: %T, spec: %T, type: %T?, field=%s", decl, spec, typ, name)
		}
	}
	return nil
}
