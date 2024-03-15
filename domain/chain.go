package domain

import "time"

//go:generate easyjson -all chain.go

const (
	FilterByStores    string = "stores"
	FilterByCities    string = "cities"
	FilterByRegions   string = "regions"
	FilterByCountries string = "countries"
)

// Chain is Layout of chain kind
type Chain struct {
	LayoutID   string      `json:"layout_id,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Title      string      `json:"title,omitempty"`
	Languages  string      `json:"languages,omitempty"`
	CRMKey     string      `json:"crm_key,omitempty"`
	Brands     string      `json:"brands,omitempty"`
	Currency   string      `json:"currency,omitempty"`
	Options    string      `json:"options,omitempty"`
	Notes      string      `json:"notes,omitempty"`
	ValidFrom  *time.Time  `json:"valid_from,omitempty"`
	ValidTo    *time.Time  `json:"valid_to,omitempty"`
	IsActive   bool        `json:"is_active"`
	Creator    string      `json:"creator,omitempty"`
	CreatedAt  *time.Time  `json:"created_at,omitempty"`
	Modifier   *string     `json:"modifier,omitempty"`
	ModifiedAt *time.Time  `json:"modified_at,omitempty"`
	ReadOnly   bool        `json:"read_only"`
	Stores     ChainStores `json:"stores,omitempty"`
}

type Chains []Chain

// ChainStore is retail object abstraction of store/shop/renter
type ChainStore struct {
	StoreID    string         `json:"store_id,omitempty"`
	LayoutID   string         `json:"layout_id,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Title      string         `json:"title,omitempty"`
	CRMKey     string         `json:"crm_key,omitempty"`
	Brands     string         `json:"brands,omitempty"`
	Statistics string         `json:"statistics,omitempty"`
	LocationID string         `json:"location_id,omitempty"`
	Area       float64        `json:"area,omitempty"`
	Currency   string         `json:"currency,omitempty"`
	Options    string         `json:"options,omitempty"`
	Notes      string         `json:"notes,omitempty"`
	ValidFrom  *time.Time     `json:"valid_from,omitempty"`
	ValidTo    *time.Time     `json:"valid_to,omitempty"`
	IsActive   bool           `json:"is_active"`
	Creator    string         `json:"creator,omitempty"`
	CreatedAt  *time.Time     `json:"created_at,omitempty"`
	Modifier   *string        `json:"modifier,omitempty"`
	ModifiedAt *time.Time     `json:"modified_at,omitempty"`
	Entrances  ChainEntrances `json:"entrances,omitempty"`
	Zones      ChainZones     `json:"zones,omitempty"`
	Devices    ChainDevices   `json:"devices,omitempty"`
}

type ChainStores []ChainStore

// ChainEntrance is enter/exit entity and its properties
type ChainEntrance struct {
	EntranceID string       `json:"entrance_id,omitempty"`
	LayoutID   string       `bson:"layout_id,omitempty" json:"layout_id,omitempty"`
	StoreID    string       `bson:"store_id,omitempty" json:"store_id,omitempty"`
	Kind       string       `bson:"kind,omitempty" json:"kind,omitempty"`
	Title      string       `bson:"title,omitempty" json:"title,omitempty"`
	Options    string       `bson:"options,omitempty" json:"options,omitempty"`
	Notes      string       `bson:"notes,omitempty" json:"notes,omitempty"`
	ValidFrom  *time.Time   `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo    *time.Time   `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	IsActive   bool         `json:"is_active"`
	Creator    string       `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt  *time.Time   `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier   *string      `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt *time.Time   `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Sensors    ChainSensors `json:"sensors,omitempty"`
}

type ChainEntrances []ChainEntrance

// ChainZone is area of some physical territory,
// cash-group, cash...
type ChainZone struct {
	ZoneID     string         `json:"zone_id,omitempty" extensions:"x-order=0"`
	ParentID   *string        `json:"parent_id,omitempty" extensions:"x-order=1"`
	LayoutID   string         `bson:"layout_id,omitempty" json:"layout_id,omitempty" extensions:"x-order=2"`
	StoreID    string         `bson:"store_id,omitempty" json:"store_id,omitempty" extensions:"x-order=3"`
	Kind       string         `bson:"kind,omitempty" json:"kind,omitempty" extensions:"x-order=4"`
	Title      string         `bson:"title,omitempty" json:"title,omitempty" extensions:"x-order=5"`
	Area       float64        `json:"area,omitempty" extensions:"x-order=6"`
	Options    string         `bson:"options,omitempty" json:"options,omitempty" extensions:"x-order=7"`
	Notes      string         `bson:"notes,omitempty" json:"notes,omitempty" extensions:"x-order=8"`
	ValidFrom  *time.Time     `bson:"valid_from,omitempty" json:"valid_from,omitempty" extensions:"x-order=9"`
	ValidTo    *time.Time     `bson:"valid_to,omitempty" json:"valid_to,omitempty" extensions:"x-order=10"`
	IsActive   bool           `json:"is_active"`
	Creator    string         `bson:"creator,omitempty" json:"creator,omitempty" extensions:"x-order=11"`
	CreatedAt  *time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty" extensions:"x-order=12"`
	Modifier   *string        `bson:"modifier,omitempty" json:"modifier,omitempty" extensions:"x-order=13"`
	ModifiedAt *time.Time     `bson:"modified_at,omitempty" json:"modified_at,omitempty" extensions:"x-order=14"`
	Entrances  ChainEntrances `json:"entrances,omitempty" extensions:"x-order=15"`
	Sensors    ChainSensors   `json:"sensors,omitempty" extensions:"x-order=16"`
}

type ChainZones []ChainZone

// ChainDevice properties of device,
// controller, navigator, etc...
type ChainDevice struct {
	DeviceID string  `json:"device_id,omitempty"`
	LayoutID string  `json:"layout_id,omitempty"`
	StoreID  string  `json:"store_id,omitempty"`
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
	DCMode     string       `json:"dcmode,omitempty"`
	Login      string       `json:"login,omitempty"`
	Password   string       `json:"password,omitempty"`
	Options    string       `json:"options,omitempty"`
	Notes      string       `json:"notes,omitempty"`
	ValidFrom  *time.Time   `bson:"valid_from,omitempty" json:"valid_from,omitempty"`
	ValidTo    *time.Time   `bson:"valid_to,omitempty" json:"valid_to,omitempty"`
	Creator    string       `bson:"creator,omitempty" json:"creator,omitempty"`
	CreatedAt  *time.Time   `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Modifier   *string      `bson:"modifier,omitempty" json:"modifier,omitempty"`
	ModifiedAt *time.Time   `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Sensors    ChainSensors `json:"sensors,omitempty"`
	Delay      *DelayPoint  `json:"delay,omitempty"`
}

type ChainDevices []ChainDevice

// ChainSensor is Sensor entity which make the measurements
type ChainSensor struct {
	SensorID   string     `json:"sensor_id,omitempty"`
	DeviceID   string     `json:"device_id,omitempty"`
	LayoutID   string     `bson:"layout_id,omitempty" json:"layout_id,omitempty"`
	StoreID    string     `bson:"store_id,omitempty" json:"store_id,omitempty"`
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

type ChainSensors []ChainSensor

type BindingChainSensorEntrance struct {
	SensorID     string     `json:"sensor_id,omitempty"`
	EntranceID   string     `json:"entrance_id,omitempty"`
	KindEntrance string     `json:"kind_entrance,omitempty"`
	Direction    string     `json:"direction,omitempty"`
	KIn          float64    `json:"k_in,omitempty"`
	KOut         float64    `json:"k_out,omitempty"`
	Options      string     `json:"options,omitempty"`
	ValidFrom    *time.Time `json:"valid_from,omitempty"`
	ValidTo      *time.Time `json:"valid_to,omitempty"`
	Creator      string     `json:"creator,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Modifier     *string    `json:"modifier,omitempty"`
	ModifiedAt   *time.Time `json:"modified_at,omitempty"`
}

type BindingChainSensorZone struct {
	SensorID   string     `json:"sensor_id,omitempty"`
	ZoneID     string     `json:"zone_id,omitempty"`
	KindZone   string     `json:"kind_zone,omitempty"`
	Options    string     `json:"options,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}

type BindingsChainSensorZone []BindingChainSensorZone

type BindingChainEntranceZone struct {
	EntranceID string     `json:"entrance_id,omitempty"`
	ZoneID     string     `json:"zone_id,omitempty"`
	KindZone   string     `json:"kind_zone,omitempty"`
	Direction  string     `json:"direction,omitempty"`
	Options    string     `json:"options,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}

type BindingChainEntranceStore struct {
	EntranceID string     `json:"entrance_id,omitempty"`
	StoreID    string     `json:"store_id,omitempty"`
	KindStore  string     `json:"kind_store,omitempty"`
	Direction  string     `json:"direction,omitempty"`
	Options    string     `json:"options,omitempty"`
	ValidFrom  *time.Time `json:"valid_from,omitempty"`
	ValidTo    *time.Time `json:"valid_to,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Modifier   *string    `json:"modifier,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}
