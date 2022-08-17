package fixture

import "context"

// F is function @FUN0
func F(x1 int, y1 string, args ...interface{}) (string, error) {
	// inner function :IGNORED:
	return "", nil
} // F is function @FUN1 :IGNORED:

// F2 is function @FUN2
func F2(
	x1 int, // x1 is int @arg1 :IGNORED:
	y1 string, // y2 is int @arg2 :IGNORED:
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
func F4(x1 int /* x1 of F4 @arg4 :IGNORED:*/ /* y1 of F4 @arg5 :IGNORED:*/, y1 /* y1 of F4 @arg6 :IGNORED:*/ string /* y1 of F4 @arg7 :IGNORED:*/, args ...interface{}) ( /* result if F4 @ret3 :IGNORED */ string /* result if F4 @ret4 :IGNORED */ /* err of F4 @ret5 :IGNORED */, error /* err of F4 @ret6 :IGNORED */) {
	return "", nil
}
