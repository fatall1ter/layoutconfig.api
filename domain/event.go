package domain

import "time"

type Event struct {
	ID        string     `json:"id"`
	Key       string     `json:"key,omitempty"`
	EventTime time.Time  `json:"event_time"`
	Kind      string     `json:"kind,omitempty"`
	Message   string     `json:"message,omitempty"`
	Severity  string     `json:"severity,omitempty"`
	LayoutID  string     `json:"layout_id,omitempty"`
	StoreID   string     `json:"store_id,omitempty"`
	Source    *Source    `json:"source,omitempty"`
	Creator   string     `json:"creator,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Source struct {
	Kind   string  `json:"kind"`
	Params []Param `json:"params"`
}

func (s Source) GetParamByName(name string) string {
	for _, p := range s.Params {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//easyjson:json
type Events []Event

type IEventRepo interface {
	// events
	FindChainEvents(from, to time.Time, layoutID, storeID, key, kind, severity string,
		limit, offset int64) (Events, int64, error)
	FindConsumerChainEvents(subscriberID, layoutID, storeID, key, kind, severity string,
		from time.Time, cancel <-chan struct{}) chan Event
}
