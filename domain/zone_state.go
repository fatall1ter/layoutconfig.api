package domain

import "time"

// ZoneState describes detectionPoint's state online/offline
type ZoneState struct {
	ZoneID     string     `json:"zone_id"`
	StoreID    string     `json:"store_id,omitempty"`
	LayoutID   string     `json:"layout_id,omitempty"`
	State      string     `json:"state,omitempty"`
	ValidFrom  time.Time  `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
	Comment    string     `json:"comment,omitempty"`
}

//easyjson:json
type ZoneStates []ZoneState
