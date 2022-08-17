package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Printf("!! %+v", err)
	}
}

func run() error {
	fset := token.NewFileSet()
	// filename := "./testdata/fixture/struct.go"
	// filename := "./testdata/fixture/const.go"
	filename := "./testdata/fixture/embedded.go"

	tree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	f, err := Collect(fset, tree)
	if err != nil {
		return fmt.Errorf("collect: file=%s, %w", filename, err)
	}
	fmt.Println(f)
	return nil
}

func Collect(fset *token.FileSet, t *ast.File) (*File, error) {
	for _, cg := range t.Comments {
		log.Println(strings.TrimSpace(cg.Text()))
	}
	log.Println("----------------------------------------")
	c := &collector{fset: fset}
	f := &File{structMap: map[string]*Struct{}}
	return f, c.CollectFromFile(f, t)
}

type collector struct {
	fset *token.FileSet
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
		Name: name,
		Decl: decl,
		Spec: spec,
	}
	f.structMap[name] = s
	switch typ := spec.Type.(type) {
	case *ast.Ident:
		// type <S> <S>
		// type <S> = <S>
	case *ast.StructType:
		// type <S> struct { ... }
		if err := c.CollectFromStructType(s, decl, spec, typ); err != nil {
			return err
		}
	default:
		log.Printf("unexpected decl: %T, spec: %T, type: %T?", decl, spec, typ)
	}

	return nil
}

func (c *collector) CollectFromStructType(s *Struct, decl *ast.GenDecl, spec *ast.TypeSpec, typ *ast.StructType) error {
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
			Name:      name,
			Field:     field,
			Anonymous: anonymous,
		})
	}
	return nil
}

type File struct {
	Names     []string
	structMap map[string]*Struct
}

type Struct struct {
	Name   string
	Fields []*Field
	Decl   *ast.GenDecl
	Spec   *ast.TypeSpec
}

// TODO: see spec
func (s *Struct) Doc() *ast.CommentGroup {
	return s.Decl.Doc
}
func (s *Struct) Comment() *ast.CommentGroup {
	return nil
}

type Field struct {
	Name      string
	Field     *ast.Field
	Anonymous bool
}

// TODO: see spec
func (s *Field) Doc() *ast.CommentGroup {
	return s.Field.Doc
}
func (s *Field) Comment() *ast.CommentGroup {
	return s.Field.Comment
}

type Target interface {
	HasDoc
	HasComment
}
type HasDoc interface {
	Doc() *ast.CommentGroup
}
type HasComment interface {
	Comment() *ast.CommentGroup
}

var _ Target = (*Struct)(nil)
