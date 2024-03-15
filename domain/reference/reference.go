package reference

// RefCategory renter category
type RefCategory struct {
	Ref
}

//easyjson:json
type RefCategories []RefCategory

// RefCategory renter category
type RefPrice struct {
	Ref
}

//easyjson:json
type RefPrices []RefPrice

// RefCategory renter category
type RefKindZone struct {
	Ref
}

//easyjson:json
type RefKindZones []RefKindZone

// RefKindEnter enter category
type RefKindEnter struct {
	Ref
}

//easyjson:json
type RefKindEnters []RefKindEnter

type Ref struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Kind  string `json:"kind,omitempty"`
}
