package commentof

// toplevel comment 0  :IGNORED:

// S is struct @S0
type S struct {
	// in struct comment 0  :IGNORED:

	// ExportedString is exported string @F0
	ExportedString string

	// in struct comment 1  :IGNORED:

	ExportedString2 string // ExportedString2 is exported string @F1

	// ExportedString3 is exported string @F2
	ExportedString3 string // ExportedString3 is exported string @F3

	// Nested is struct @SS0
	Nested struct { // in struct comment 2  :IGNORED:

		// ExportedString is exported string @FF0
		ExportedString string // ExportedString is exported string @@FF1

		// in struct comment 3  :IGNORED:
	} // Nested is struct @SS1
	// in struct comment 4  :IGNORED:

	// unexportedString is unexported string @U1  :IGNORED:
	unexportedString string
} // S is struct @S1

const (
	// CONSTANT_STRING is constant string @C0
	CONSTNAT_STRING = ""

	CONSTNAT_STRING2 = "" // CONSTANT_STRING2 is constant string @C1

	// CONSTANT_STRING3 is constant string @C2
	CONSTNAT_STRING3 = "" // CONSTANT_STRING3 is constant string  @C3
)

// CONSTANT_STRING4 is constant string @C4
const CONSTNAT_STRING4 = ""

const CONSTNAT_STRING5 = "" // CONSTANT_STRING5 is constant string  @C5
