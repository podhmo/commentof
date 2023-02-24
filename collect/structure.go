package collect

import "go/token"

type Package struct {
	Files      map[string]*File   `json:"-"`
	Interfaces map[string]*Object `json:"interfaces"`
	Functions  map[string]*Func   `json:"functions"`
	Types      map[string]*Object `json:"types"`

	FileNames []string `json:"filenames"`
	Names     []string `json:"names"`
}

func NewPackage() *Package {
	return &Package{
		Files: map[string]*File{}, FileNames: []string{},
		Interfaces: map[string]*Object{},
		Functions:  map[string]*Func{},
		Types:      map[string]*Object{},
		Names:      []string{},
	}
}

type File struct {
	Interfaces map[string]*Object `json:"interfaces"`
	Functions  map[string]*Func   `json:"functions"`
	Types      map[string]*Object `json:"types"`
	Names      []string           `json:"names"`
}

func NewFile() *File {
	return &File{
		Interfaces: map[string]*Object{},
		Functions:  map[string]*Func{},
		Types:      map[string]*Object{},
		Names:      []string{},
	}
}

type Func struct {
	Name string    `json:"name"`
	Pos  token.Pos `json:"-"`

	Recv string `json:"recv,omitempty"`

	Params     map[string]*Field `json:"params"`
	ParamNames []string          `json:"paramnames"`

	Returns     map[string]*Field `json:"returns"`
	ReturnNames []string          `json:"returnnames"`

	Doc string `json:"doc"` // associated documentation; or nil (decl or spec?)
}

type Object struct {
	Name   string      `json:"name"`
	Pos    token.Pos   `json:"-"`
	Token  token.Token `json:"-"`
	Parent *Object     `json:"-"`

	Fields     map[string]*Field `json:"fields,omitempty"`
	FieldNames []string          `json:"fieldnames,omitempty"`

	Methods     map[string]*Func `json:"methods,omitempty"`
	MethodNames []string         `json:"methodnames,omitempty"`

	Doc     string `json:"doc"`     // associated documentation; or nil (decl or spec?)
	Comment string `json:"comment"` // line comments; or nil
}

type Field struct {
	Name      string    `json:"name"`
	Pos       token.Pos `json:"-"`
	Embedded  bool      `json:"embedded"`
	Anonymous *Object   `json:"annonymous,omitempty"`

	Doc     string `json:"doc"`     // associated documentation; or nil
	Comment string `json:"comment"` // line comments; or nil
}
