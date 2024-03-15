package infra

import "time"

// RequestEvent parameter for event channel filtering
type RequestEvent struct {
	Key      string     `json:"key,omitempty"`
	Kind     string     `json:"kind,omitempty"`
	Severity string     `json:"severity,omitempty"`
	LayoutID string     `json:"layout_id,omitempty"`
	StoreID  string     `json:"store_id,omitempty"`
	From     *time.Time `json:"from,omitempty"`
}
