package fixture

// Base is struct @S10
type Base struct {
	// ExportedString is exported string @F10
	ExportedString string
}

// S10 is struct @S10
type S10 struct {
	Base

	// ExportedString2 is exported string @F11
	ExportedString2 string
}
