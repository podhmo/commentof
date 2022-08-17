package collect

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
)

type Package struct {
	Files      map[string]*File   `json:"-"`
	Structs    map[string]*Object `json:"structs"`
	Interfaces map[string]*Object `json:"interfaces"`

	FileNames []string `json:"filenames"`
	Names     []string `json:"names"`
}

type File struct {
	Structs    map[string]*Object `json:"structs"`
	Interfaces map[string]*Object `json:"interfaces"`

	Names []string `json:"names"`
}

type Object struct {
	Name   string      `json:"name"`
	Token  token.Token `json:"-"`
	Parent *Object     `json:"-"`

	Fields     map[string]*Field `json:"fields"`
	FieldNames []string          `json:"fieldnames"`

	Doc     string `json:"doc"`     // associated documentation; or nil (decl or spec?)
	Comment string `json:"comment"` // line comments; or nil
}

type Field struct {
	Name      string  `json:"name"`
	Embedded  bool    `json:"embedded"`
	Anonymous *Object `json:"annonymous,omitempty"`

	Doc     string `json:"doc"`     // associated documentation; or nil
	Comment string `json:"comment"` // line comments; or nil
	// TODO: tag
}

type Collector struct {
	Fset *token.FileSet
	Dot  string
}

func (c *Collector) CollectFromPackage(p *Package, t *ast.Package) error {
	for filename, ft := range t.Files {
		p.FileNames = append(p.FileNames, filename)
		f := &File{
			Structs:    map[string]*Object{},
			Interfaces: map[string]*Object{},
			Names:      []string{},
		}
		p.Files[filename] = f
		if err := c.CollectFromFile(f, ft); err != nil {
			return fmt.Errorf("collect file: %s: %w", filename, err)
		}

		p.Names = append(p.Names, f.Names...)
		for name, s := range f.Structs {
			p.Structs[name] = s
		}
		for name, s := range f.Interfaces {
			p.Interfaces[name] = s
		}
	}
	return nil
}

func (c *Collector) CollectFromFile(f *File, t *ast.File) error {
	for _, decl := range t.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl, *ast.BadDecl:
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

func (c *Collector) CollectFromGenDecl(f *File, decl *ast.GenDecl) error {
	// decl.Tok == token.TYPE
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
	}
	if s.Doc == "" && decl.Doc != nil {
		s.Doc = decl.Doc.Text()
	}

	switch typ := spec.Type.(type) {
	case *ast.Ident:
		// type <S> <S>
		// type <S> = <S>
	case *ast.StructType:
		// type <S> struct { ... }
		f.Structs[name] = s
		if err := c.CollectFromStructType(f, s, decl, spec, typ); err != nil {
			return err
		}
	case *ast.InterfaceType:
		// type <S> struct { ... }
		f.Interfaces[name] = s
		if err := c.CollectFromInterfaceType(f, s, decl, spec, typ); err != nil {
			return err
		}
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
	default:
		return "", false
	}
}

func (c *Collector) CollectFromStructType(f *File, s *Object, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.StructType) error {
	s.Token = token.STRUCT
	for _, field := range typ.Fields.List {
		name := ""
		anonymous := false
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else {
			anonymous = true
			if typename, ok := typeString(field.Type); ok {
				name = typename
			} else {
				name = fmt.Sprintf("??%T", field.Type) // TODO: NG:embedded
				log.Printf("unexpected embedded field type: %T, spec: %T, struct: %T, field:%v", decl, spec, typ, field.Type)
			}
		}

		s.FieldNames = append(s.FieldNames, name)
		fieldof := &Field{
			Name:     name,
			Doc:      field.Doc.Text(),
			Comment:  field.Comment.Text(),
			Embedded: anonymous,
		}
		s.Fields[name] = fieldof

		switch typ := field.Type.(type) {
		case *ast.Ident, *ast.FuncType, *ast.SelectorExpr:
		case *ast.StructType:
			// type <S> struct { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),     // xxx
				Comment:    field.Comment.Text(), // xxx
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromStructType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		case *ast.InterfaceType:
			// type <S> struct { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),     // xxx
				Comment:    field.Comment.Text(), // xxx
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromInterfaceType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		default:
			log.Printf("unexpected decl: %T, spec: %T, type: %T?, field=%s", decl, spec, typ, name)
		}
	}
	return nil
}

func (c *Collector) CollectFromInterfaceType(f *File, s *Object, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.InterfaceType) error {
	s.Token = token.INTERFACE
	for _, field := range typ.Methods.List {
		name := ""
		anonymous := false
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else {
			anonymous = true
			if typename, ok := typeString(field.Type); ok {
				name = typename
			} else {
				name = fmt.Sprintf("??%T", field.Type) // TODO: NG:embedded
				log.Printf("unexpected embedded field type: %T, spec: %T, struct: %T, field:%v", decl, spec, typ, field.Type)
			}
		}

		s.FieldNames = append(s.FieldNames, name)
		fieldof := &Field{
			Name:     name,
			Doc:      field.Doc.Text(),
			Comment:  field.Comment.Text(),
			Embedded: anonymous,
		}
		s.Fields[name] = fieldof

		switch typ := field.Type.(type) {
		case *ast.Ident, *ast.FuncType, *ast.SelectorExpr:
		case *ast.InterfaceType:
			// type <S> struct { ... }
			name := s.Name + c.Dot + name
			f.Names = append(f.Names, name)
			anonymous := &Object{
				Name:       name,
				Parent:     s,
				Doc:        field.Doc.Text(),     // xxx
				Comment:    field.Comment.Text(), // xxx
				FieldNames: []string{},
				Fields:     map[string]*Field{},
			}
			fieldof.Anonymous = anonymous
			if err := c.CollectFromInterfaceType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		default:
			log.Printf("unexpected decl: %T, spec: %T, type: %T?, field=%s", decl, spec, typ, name)
		}
	}
	return nil
}
