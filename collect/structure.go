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

type File struct {
	Structs    map[string]*Object `json:"structs"`
	Interfaces map[string]*Object `json:"interfaces"`
	Functions  map[string]*Func   `json:"functions"`

	Names []string `json:"names"`
}

type Func struct {
	Name string `json:"name"`

	Params     map[string]*Field `json:"params"`
	ParamNames []string          `json:"paramnames"`

	Returns     map[string]*Field `json:"returns"`
	ReturnNames []string          `json:"returnnames"`

	// 	Doc  *CommentGroup // associated documentation; or nil
	// 	Recv *FieldList    // receiver (methods); or nil (functions)
	// 	Name *Ident        // function/method name
	// 	Type *FuncType     // function signature: type and value parameters, results, and position of "func" keyword
	// 	Body *BlockStmt    // function body; or nil for external (non-Go) function
	// }

	Doc string `json:"doc"` // associated documentation; or nil (decl or spec?)
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
