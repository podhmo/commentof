package commentof

// S is struct
type S struct {
	// ExportedString is exported string
	ExportedString string

	ExportedString2 string // ExportedString2 is exported string

	// ExportedString3 is exported string
	ExportedString3 string // ExportedString3 is exported string

	Nested struct {
		// ExportedString is exported string
		ExportedString string // ExportedString is exported string
	}

	// unexportedString is unexported string
	unexportedString string
}

const (
	// CONSTANT_STRING is constant string
	CONSTNAT_STRING = ""

	CONSTNAT_STRING2 = "" // CONSTANT_STRING2 is constant string

	// CONSTANT_STRING3 is constant string
	CONSTNAT_STRING3 = "" // CONSTANT_STRING3 is constant string
)

// CONSTANT_STRING4 is constant string
const CONSTNAT_STRING4 = ""

const CONSTNAT_STRING5 = "" // CONSTANT_STRING5 is constant string
