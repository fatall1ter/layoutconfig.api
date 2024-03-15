package domain

//go:generate easyjson -all zone_data_evaluation.go

import "time"

type ZoneDataEvaluation struct {
	LayoutID              string    `json:"layout_id"`
	StoreID               string    `json:"store_id"`
	ServiceChannelBlockID string    `json:"service_channel_block_id"`
	RecordTime            time.Time `json:"record_time"`
	IsFull                bool      `json:"is_full"`
	Comment               string    `json:"comment"`
}

//easyjson:json
type ZoneDataEvaluations []ZoneDataEvaluation
