package fixture

// toplevel comment 3  :IGNORED:

// I is interface @I0
type I interface {
	// Exported is exported method @IF0
	Exported() string

	Exported2() string // Exported2 is exported method  @IF1

	// Exported3 is exported method
	Exported3() string // Exported2 is exported method  @IF2

	// unexported is unexported method @IUF0 :IGNORED:
	unexported() string
} // I is interface @I1

// I2 is interface @I2
type I2 interface {
	// embedded I @IF3
	I // embedded I @IF4
}

// I3 is interface @I3
type I3 interface {
	I

	// embedded anonymous @IUF1 :IGNORED:
	interface {
		// Nested is exported method @IFF0
		Nested() string

		Nested2() string // Nested is exported method @IFF1

	} // embedded anonymous @IF5
}
