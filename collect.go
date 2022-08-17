package commentof

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
)

func FileAST(fset *token.FileSet, t *ast.File) (*File, error) {
	// for _, cg := range t.Comments {
	// 	log.Println(strings.TrimSpace(cg.Text()))
	// }
	// log.Println("----------------------------------------")
	c := &collector{fset: fset, dot: "."}
	f := &File{Structs: map[string]*Struct{}}
	return f, c.CollectFromFile(f, t)
}

type collector struct {
	fset *token.FileSet
	dot  string
}

func (c *collector) CollectFromFile(f *File, t *ast.File) error {
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

func (c *collector) CollectFromGenDecl(f *File, decl *ast.GenDecl) error {
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

func (c *collector) CollectFromTypeSpec(f *File, decl *ast.GenDecl, spec *ast.TypeSpec) error {
	name := spec.Name.Name
	f.Names = append(f.Names, name)
	s := &Struct{
		Name:    name,
		Doc:     spec.Doc.Text(),
		Comment: spec.Comment.Text(),
		Fields:  []*Field{},
	}
	if s.Doc == "" && decl.Doc != nil {
		s.Doc = decl.Doc.Text()
	}

	f.Structs[name] = s
	switch typ := spec.Type.(type) {
	case *ast.Ident:
		// type <S> <S>
		// type <S> = <S>
	case *ast.StructType:
		// type <S> struct { ... }
		if err := c.CollectFromStructType(f, s, decl, spec, typ); err != nil {
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
	default:
		return "", false
	}
}

func (c *collector) CollectFromStructType(f *File, s *Struct, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.StructType) error {
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

		s.Fields = append(s.Fields, &Field{
			Name:     name,
			Doc:      field.Doc.Text(),
			Comment:  field.Comment.Text(),
			Embedded: anonymous,
		})

		switch typ := field.Type.(type) {
		case *ast.Ident:
		case *ast.StructType:
			// type <S> struct { ... }
			name := s.Name + c.dot + name
			f.Names = append(f.Names, name)
			anonymous := &Struct{
				Name:    name,
				Parent:  s,
				Doc:     field.Doc.Text(),     // xxx
				Comment: field.Comment.Text(), // xxx
				Fields:  []*Field{},
			}
			s.Fields[len(s.Fields)-1].Anonymous = anonymous
			if err := c.CollectFromStructType(f, anonymous, decl, spec, typ); err != nil {
				return err
			}
		default:
			log.Printf("unexpected decl: %T, spec: %T, type: %T?, field=%s", decl, spec, typ, name)

		}
	}
	return nil
}

type File struct {
	Structs map[string]*Struct `json:"structs"`
	Names   []string           `json:"names"`
}

type Struct struct {
	Name   string   `json:"name"`
	Parent *Struct  `json:"-"`
	Fields []*Field `json:"fields"`

	Doc     string `json:"doc"`     // associated documentation; or nil (decl or spec?)
	Comment string `json:"comment"` // line comments; or nil
}

// func (s *Struct) MarshalJSON() ([]byte, error) {
// 	type T Struct
// 	inner := (T)(*s)
// 	inner.Parent = nil
// 	return json.Marshal(inner)
// }

type Field struct {
	Name      string  `json:"name"`
	Embedded  bool    `json:"embedded"`
	Anonymous *Struct `json:"annonymous,omitempty"`

	Doc     string `json:"doc"`     // associated documentation; or nil
	Comment string `json:"comment"` // line comments; or nil
	// TODO: tag
}
