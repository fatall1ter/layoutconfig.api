package domain

import "time"

//go:generate easyjson -all mall.go

type Direction int

func (d Direction) String() string {
	switch d {
	case ForwardDirection:
		return "forward"
	case ReverseDirection:
		return "reverse"
	default:
		return "undefined"
	}
}

const (
	ForwardDirection Direction = 1
	ReverseDirection Direction = 2
)

// Mall is Layout of mall kind
type Mall struct {
	LayoutID   string     `json:"layout_id,omitempty"`
	Kind       string     `json:"kind,omitempty"`
	Title      string     `json:"title,omitempty"`
	Languages  string     `json:"languages,omitempty"`
	CRMKey     string     `json:"crm_key,omitempty"`
	Options    string     `json:"options,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	IsActive   bool       `json:"is_active"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
	ReadOnly   bool       `json:"read_only"`
}

type Malls []Mall

// MallEntrance is enter/exit entity and its properties
type MallEntrance struct {
	EntranceID string      `json:"entrance_id,omitempty"`
	LayoutID   string      `json:"layout_id,omitempty"`
	FloorID    string      `json:"floor_id,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Title      string      `json:"title,omitempty"`
	Options    string      `json:"options,omitempty"`
	Notes      string      `json:"notes,omitempty"`
	ValidFrom  *time.Time  `json:"valid_from,omitempty"`
	ValidTo    *time.Time  `json:"valid_to,omitempty"`
	IsActive   bool        `json:"is_active"`
	Creator    string      `json:"creator,omitempty"`
	CreatedAt  *time.Time  `json:"created_at,omitempty"`
	Modifier   *string     `json:"modifier,omitempty"`
	ModifiedAt *time.Time  `json:"modified_at,omitempty"`
	Sensors    MallSensors `json:"sensors,omitempty"`
}

type MallEntrances []MallEntrance

// MallZone is area of some physical territory,
// gellary, floor, renter room, cash...
type MallZone struct {
	ZoneID     string               `bson:"zone_id,omitempty" json:"zone_id,omitempty"`
	ParentID   *string              `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	LayoutID   string               `bson:"layout_id,omitempty" json:"layout_id,omitempty"`
	Kind       string               `bson:"kind,omitempty" json:"kind,omitempty"`
	Title      string               `bson:"title,omitempty" json:"title,omitempty"`
	Area       float64              `json:"area,omitempty"`
	Options    string               `bson:"options,omitempty" json:"options,omitempty"`
	Notes      string               `bson:"notes,omitempty" json:"notes,omitempty"`
	ValidFrom  *time.Time           `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo    *time.Time           `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	IsActive   bool                 `json:"is_active"`
	Creator    string               `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt  *time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier   *string              `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt *time.Time           `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Entrances  BindingsEntranceZone `json:"entrances,omitempty"`
	Sensors    MallSensors          `json:"sensors,omitempty"`
}

type MallZones []MallZone

// Renter arendator of the mall
type Renter struct {
	RenterID       string     `bson:"renter_id,omitempty" json:"renter_id,omitempty"`
	Title          string     `bson:"title,omitempty" json:"title,omitempty"`
	LayoutID       string     `bson:"layout_id,omitempty" json:"layout_id,omitempty"`
	CategoryID     string     `bson:"category_id,omitempty" json:"category_id,omitempty"`
	PriceSegmentID string     `bson:"price_segment_id,omitempty" json:"price_segment_id,omitempty"`
	TimeOpen       *time.Time `bson:"time_open,omitempty" json:"time_open,omitempty"`
	TimeClose      *time.Time `bson:"time_close,omitempty" json:"time_close,omitempty"`
	Contract       string     `bson:"contract,omitempty" json:"contract,omitempty"`
	Options        string     `bson:"options,omitempty" json:"options,omitempty"`
	Notes          string     `bson:"notes,omitempty" json:"notes,omitempty"`
	ValidFrom      *time.Time `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo        *time.Time `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	IsActive       bool       `json:"is_active"`
	Creator        string     `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt      *time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier       *string    `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt     *time.Time `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Zones          MallZones  `bson:"zones,omitempty" json:"zones,omitempty"`
}

//easyjson:json
type Renters []Renter

// MallDevice properties of device,
// controller, navigator, etc...
type MallDevice struct {
	DeviceID string  `json:"device_id,omitempty"`
	LayoutID string  `json:"layout_id,omitempty"`
	FloorID  string  `json:"floor_id,omitempty"`
	MasterID *string `json:"master_id,omitempty"`
	Kind     string  `json:"kind,omitempty"`
	Title    string  `json:"title,omitempty"`
	IsActive bool    `json:"is_active"`
	IP       string  `json:"ip,omitempty"`
	Port     string  `json:"port,omitempty"`
	SN       string  `json:"sn,omitempty"`
	// Mode: single/master/slave
	Mode string `json:"mode,omitempty"`
	// DCMode data collector mode: active - device transmit data to server;
	// passive - server request data from device
	DCMode     string      `json:"dcmode,omitempty"`
	Login      string      `json:"login,omitempty"`
	Password   string      `json:"password,omitempty"`
	Options    string      `json:"options,omitempty"`
	Notes      string      `json:"notes,omitempty"`
	ValidFrom  *time.Time  `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo    *time.Time  `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	Creator    string      `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt  *time.Time  `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier   *string     `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt *time.Time  `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Sensors    MallSensors `json:"sensors,omitempty"`
	Delay      *DelayPoint `json:"delay,omitempty"`
}

//easyjson:json
type MallDevices []MallDevice

// MallSensor is Sensor entity which make the measurements
type MallSensor struct {
	SensorID   string     `json:"sensor_id,omitempty"`
	DeviceID   string     `json:"device_id,omitempty"`
	LayoutID   string     `bson:"layout_id,omitempty" json:"layout_id,omitempty"`
	ExternalID string     `bson:"external_id,omitempty" json:"external_id,omitempty"`
	Kind       string     `bson:"kind,omitempty" json:"kind,omitempty"`
	Title      string     `bson:"title,omitempty" json:"title,omitempty"`
	Options    string     `bson:"options,omitempty" json:"options,omitempty"`
	Notes      string     `bson:"notes,omitempty" json:"notes,omitempty"`
	ValidFrom  *time.Time `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo    *time.Time `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	IsActive   bool       `json:"is_active"`
	Creator    string     `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt  *time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier   *string    `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt *time.Time `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
}

//easyjson:json
type MallSensors []MallSensor

//
type BindingEntranceZone struct {
	EntranceID   string `json:"entrance_id,omitempty"`
	EntranceName string `json:"entrance_name,omitempty"`
	ZoneID       string `json:"zone_id,omitempty"`
	KindZone     string `json:"kind_zone,omitempty"`
	Direction    string `json:"direction,omitempty"`
	MetaData
}

//easyjson:json
type BindingsEntranceZone []BindingEntranceZone

type MetaData struct {
	Options    string     `json:"options,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}
