package fixture

import (
	"context"
	"io"
)

// F is function @FUN0
func F(x int, y string, args ...interface{}) (string, error) {
	// inner function :IGNORED:
	return "", nil
} // F is function @FUN1 :IGNORED:

// F2 is function @FUN2
func F2(
	x int, // x is int @arg1 :IGNORED:
	y string, // y is int @arg2 :IGNORED:
	args ...interface{}, // args is int @arg3 :IGNORED:
) (string, // result of F2 @ret1 :IGNORED:
	error, // error of F2 @ret2 :IGNORED:
) {
	return "", nil
}

// F3 is function @FUN3
func F3(
	context.Context,
	string,
	...interface{},
) (result string, err error) {
	return "", nil
}

// F4 is function @FUN4
func F4(x int /* x of F4 @arg4 :IGNORED:*/ /* x of F4 @arg5 :IGNORED:*/, y /* y of F4 @arg6 :IGNORED:*/ string /* y of F4 @arg7 :IGNORED:*/, args ...interface{} /* arg of F4 @arg8 :IGNORED:*/) ( /* result if F4 @ret3 :IGNORED */ string /* result if F4 @ret4 :IGNORED */ /* ret of F4 @ret5 :IGNORED */ /* err of F4 @ret6 :IGNORED */, error /* err of F4 @ret7 :IGNORED */) {
	return "", nil
}

// F5 is function @FUN5
func F5() {

}

// EmitFunc is function @FUN6
type EmitFunc func(ctx context.Context, w io.Writer) error
