package domain

//go:generate easyjson -all behavior.go

const (
	IxCashiersActivities IndexKind = "cashiers_activities"
	IxQueueLength        IndexKind = "queue_length"
	IxWorkTime           IndexKind = "work_time"
	OpMultiply           OpEnum    = "*"
	OpPlus               OpEnum    = "+"
)

type IndexKind string

type OpEnum string

type Behavior struct {
	Open           string          `json:"Open,omitempty"`
	Close          string          `json:"Close,omitempty"`
	TimeZone       string          `json:"time_zone,omitempty"`
	BehaviorConfig *BehaviorConfig `json:"behavior_config,omitempty"`
}

//easyjson:json
type Behaviors []Behavior

type BehaviorConfig struct {
	Queue           *Queue           `json:"queue,omitempty"`
	Recommendations *Recommendations `json:"recommendations,omitempty"`
	QueueThresholds *QueueThresholds `json:"queue_thresholds,omitempty"`
}

type Queue struct {
	Layouts         []QueueLayout           `json:"layouts,omitempty"`
	Stores          []QueueStore            `json:"stores,omitempty"`
	ServiceChannels []ServiceChannelElement `json:"service_channels,omitempty"`
}

type QueueLayout struct {
	LayoutID       string                `json:"layout_id,omitempty"`
	Title          string                `json:"title,omitempty"`
	Threshold      float64               `json:"threshold,omitempty"`
	ServiceChannel *LayoutServiceChannel `json:"service_channel,omitempty"`
}

type QueueStore struct {
	StoreID        string                `json:"store_id,omitempty"`
	Title          string                `json:"title,omitempty"`
	Threshold      float64               `json:"threshold,omitempty"`
	ServiceChannel *LayoutServiceChannel `json:"service_channel,omitempty"`
}

type LayoutServiceChannel struct {
	Indexes []Index `json:"indexes,omitempty"`
}

type Index struct {
	Kind   IndexKind `json:"kind,omitempty"`
	Weight int64     `json:"weight,omitempty"`
	Op     OpEnum    `json:"op,omitempty"`
}

type ServiceChannelElement struct {
	ServiceChannelID string                `json:"service_channel_id,omitempty"`
	Title            string                `json:"title,omitempty"`
	Threshold        float64               `json:"threshold,omitempty"`
	ServiceChannel   *LayoutServiceChannel `json:"service_channel,omitempty"`
}

type Recommendations struct {
	Layouts []RecommendationsLayout `json:"layouts,omitempty"`
	Stores  []Store                 `json:"stores,omitempty"`
}

type RecommendationsLayout struct {
	LayoutID             string  `json:"layout_id"`
	Title                string  `json:"title"`
	StdCoef              float64 `json:"std_coef"`
	QueueMultiplier      float64 `json:"queue_multiplier"`
	PredMinutes          uint    `json:"pred_minutes"`
	HistMinutes          uint    `json:"hist_minutes"`
	CheckoutProductivity float64 `json:"checkout_productivity"`
}

type Store struct {
	StoreID              string  `json:"store_id,omitempty"`
	Title                string  `json:"title,omitempty"`
	StdCoef              float64 `json:"std_coef,omitempty"`
	QueueMultiplier      float64 `json:"queue_multiplier,omitempty"`
	PredMinutes          uint    `json:"pred_minutes,omitempty"`
	HistMinutes          uint    `json:"hist_minutes,omitempty"`
	CheckoutProductivity float64 `json:"checkout_productivity"`
}

type QueueThresholds struct {
	Layouts               []QTLayout             `json:"layouts,omitempty"`
	Stores                []QTStore              `json:"stores,omitempty"`
	BlocksServiceChannels []BlocksServiceChannel `json:"blocks_service_channels,omitempty"`
}

type BlocksServiceChannel struct {
	BlockServiceChannelsID string  `json:"block_service_channels_id,omitempty"`
	Title                  string  `json:"title,omitempty"`
	Threshold              float64 `json:"threshold,omitempty"`
	SequenceLength         uint    `json:"sequence_length,omitempty"`
}

type QTLayout struct {
	LayoutID       string  `json:"layout_id,omitempty"`
	Title          string  `json:"title,omitempty"`
	Threshold      float64 `json:"threshold,omitempty"`
	SequenceLength uint    `json:"sequence_length,omitempty"`
}

type QTStore struct {
	StoreID        string  `json:"store_id,omitempty"`
	Title          string  `json:"title,omitempty"`
	Threshold      float64 `json:"threshold,omitempty"`
	SequenceLength uint    `json:"sequence_length,omitempty"`
}
