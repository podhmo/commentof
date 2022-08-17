package commentof

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"
)

func FileAST(fset *token.FileSet, t *ast.File) (*File, error) {
	for _, cg := range t.Comments {
		log.Println(strings.TrimSpace(cg.Text()))
	}
	log.Println("----------------------------------------")
	c := &collector{fset: fset, dot: "."}
	f := &File{structMap: map[string]*Struct{}}
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
		Doc:     spec.Doc,
		Comment: spec.Comment,
	}
	if s.Doc == nil && decl.Doc != nil {
		s.Doc = decl.Doc
	}

	f.structMap[name] = s
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

func (c *collector) CollectFromStructType(f *File, s *Struct, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.StructType) error {
	for _, field := range typ.Fields.List {
		name := ""
		anonymous := false
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else {
			name = fmt.Sprintf("??%T", field.Type) // TODO: NG:embedded
			anonymous = true
		}

		s.Fields = append(s.Fields, &Field{
			Name:     name,
			Doc:      field.Doc,
			Comment:  field.Comment,
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
				Doc:     field.Doc,     // xxx
				Comment: field.Comment, // xxx
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
	structMap map[string]*Struct
	Names     []string
}

type Struct struct {
	Name   string
	Fields []*Field

	Doc     *ast.CommentGroup // decl and spec?
	Comment *ast.CommentGroup
}

type Field struct {
	Name      string
	Embedded  bool
	Anonymous *Struct

	Doc     *ast.CommentGroup // associated documentation; or nil
	Comment *ast.CommentGroup // line comments; or nil
	// TODO: tag
}
