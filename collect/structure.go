package collect

import "go/token"

type Package struct {
	Files      map[string]*File   `json:"-"`
	Structs    map[string]*Object `json:"structs"`
	Interfaces map[string]*Object `json:"interfaces"`
	Functions  map[string]*Func   `json:"functions"`

	FileNames []string `json:"filenames"`
	Names     []string `json:"names"`
}

func NewPackage() *Package {
	return &Package{
		Files: map[string]*File{}, FileNames: []string{},
		Structs:    map[string]*Object{},
		Interfaces: map[string]*Object{},
		Functions:  map[string]*Func{},
		Names:      []string{},
	}
}

type File struct {
	Structs    map[string]*Object `json:"structs"`
	Interfaces map[string]*Object `json:"interfaces"`
	Functions  map[string]*Func   `json:"functions"`

	Names []string `json:"names"`
}

func NewFile() *File {
	return &File{
		Structs:    map[string]*Object{},
		Interfaces: map[string]*Object{},
		Functions:  map[string]*Func{},
		Names:      []string{},
	}
}

type Func struct {
	Name string `json:"name"`

	Recv string `json:"recv,omitempty"`

	Params     map[string]*Field `json:"params"`
	ParamNames []string          `json:"paramnames"`

	Returns     map[string]*Field `json:"returns"`
	ReturnNames []string          `json:"returnnames"`

	Doc string `json:"doc"` // associated documentation; or nil (decl or spec?)
}

type Object struct {
	Name   string      `json:"name"`
	Token  token.Token `json:"-"`
	Parent *Object     `json:"-"`

	Fields     map[string]*Field `json:"fields"`
	FieldNames []string          `json:"fieldnames"`

	Methods     map[string]*Func `json:"methods,omitempty"`
	MethodNames []string         `json:"methodnames,omitempty"`

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
