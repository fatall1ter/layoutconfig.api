package domain

type Entity struct {
	EntityKey string `json:"entity_key,omitempty"`
	ParentKey string `json:"parent_key,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Title     string `json:"title,omitempty"`
	Notes     string `json:"notes,omitempty"`
	Options   string `json:"options,omitempty"`
}

type Entities []Entity
