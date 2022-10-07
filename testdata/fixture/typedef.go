package fixture

import (
	"context"
	"io"
)

// EmitFunc is function
type EmitFunc func(ctx context.Context, w io.Writer) error

// MyInt is new type
type MyInt int

// IntAlias is alias
type IntAlias = int
