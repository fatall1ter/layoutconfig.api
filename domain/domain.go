/*Package domain describe models and business rules

Author: Karpov Artem, mailto: karpov@watcom.ru
Date: 2019-10-14
*/
package domain

import (
	"context"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain/reference"
)

//go:generate defimpl domain.go
//go:generate easyjson -all domain.go

// Layout is like a object:
// mall, retail, airport, railway station, transport company, store, vehicle etc...
type Layout struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Kind     string      `json:"kind"`
	Owner    CRMCustomer `json:"owner"`
	IsActive bool        `json:"is_active"`
}

//easyjson:json
type Layouts []Layout

// ACL return only accessible list ids
func (los Layouts) ACL(listIDs string) (string, int) {
	const sep string = ","
	slList := strings.Split(listIDs, sep)
	slRes := make([]string, 0)
	if listIDs == "" || listIDs == "*" {
		for _, has := range los {
			slRes = append(slRes, has.ID)
		}
		return strings.Join(slRes, sep), len(slRes)
	}
	for _, chk := range slList {
		for _, has := range los {
			if has.ID == chk {
				slRes = append(slRes, chk)
			}
		}
	}
	return strings.Join(slRes, sep), len(slRes)
}

// CRMCustomer properties of customers from CRM system
type CRMCustomer struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// DeviceDelayData delay data by device
type DeviceDelayData struct {
	DeviceID     string `json:"device_id,omitempty"`
	DelayMinutes int    `json:"delay_minutes,omitempty"`
}

type FParam struct {
	Name  string
	Value string
}

//easyjson:json
type DeviceDelayDatas []DeviceDelayData

type ACLRepoInterface interface {
	FindStoresByCities(ctx context.Context, layoutID, list string) ([]string, error)
	FindStoresByRegions(ctx context.Context, layoutID, list string) ([]string, error)
	FindStoresByCountries(ctx context.Context, layoutID, list string) ([]string, error)
	FindEntrances(ctx context.Context, layoutID, storeIDs string) ([]string, error)
}

type LayoutRepo interface {
	ACLRepoInterface
	FindLayouts(ctx context.Context, loc, dt string, offset, limit int64) (Layouts, int64, error)
	FindLayoutByID(context.Context, string, string, string) (*Layout, error)

	// chains
	FindChains(ctx context.Context, loc, date, layoutID string, offset, limit int64) (Chains, int64, error)
	FindChainByID(ctx context.Context, loc, layoutID, dt string) (*Chain, error)
	AddChain(context.Context, Chain) (string, error)
	DelChain(context.Context, string) (int64, error)
	UpdChain(context.Context, Chain, bool, string) (int64, error)

	// chain/stores
	FindChainStores(ctx context.Context,
		loc, date, layoutID, crmKey string, offset, limit int64, filters string) (ChainStores, int64, error)
	FindChainStoreByID(ctx context.Context, loc, storeID, dt string) (*ChainStore, error)
	AddChainStore(context.Context, ChainStore) (string, error)
	DelChainStore(context.Context, string) (int64, error)
	UpdChainStore(context.Context, ChainStore, bool, string) (int64, error)

	// chain/entrances
	FindChainEntrances(ctx context.Context,
		loc, date, layoutID, storeID, kind, entranceIDs string, offset, limit int64) (ChainEntrances, int64, error)
	FindChainEntranceByID(ctx context.Context, loc, entranceID, dt string) (*ChainEntrance, error)
	AddChainEntrance(context.Context, ChainEntrance) (string, error)
	DelChainEntrance(context.Context, string) (int64, error)
	UpdChainEntrance(context.Context, ChainEntrance, bool, string) (int64, error)

	// chain/zones
	FindChainZones(ctx context.Context, loc, date, layoutID, storeID, kind, parentID, isOnline, isActive string,
		offset, limit int64) (ChainZones, int64, error)
	FindChainZoneByID(ctx context.Context, loc, zoneID, dt string) (*ChainZone, error)
	AddChainZone(context.Context, ChainZone) (string, error)
	DelChainZone(context.Context, string) (int64, error)
	UpdChainZone(context.Context, ChainZone, bool, string) (int64, error)

	// chain/devices
	FindChainDevices(ctx context.Context, loc, date, layoutID, storeID, kind, isActive, dcMode, sn, mode string,
		offset, limit int64) (ChainDevices, int64, error)
	FindChainDeviceByID(ctx context.Context, loc, deviceID, dt string) (*ChainDevice, error)
	AddChainDevice(context.Context, ChainDevice) (string, error)
	DelChainDevice(context.Context, string) (int64, error)
	UpdChainDevice(context.Context, ChainDevice, bool, string) (int64, error)
	FindChainDeviceDelays(context.Context, string, []string, []string) (DeviceDelayDatas, error)
	// device/tracks
	FindChainDeviceTracks(ctx context.Context, layoutID, storeID, deviceID string, from, to time.Time,
		offset, limit int64) (Tracks, int64, error)
	FindChainDeviceTracksAt(ctx context.Context,
		layoutID, storeID, deviceID string, at time.Time, accuracy time.Duration) (Tracks, error)

	// chain/sensors
	FindChainSensors(ctx context.Context,
		loc, date, layoutID, storeID, deviceID, kind string, offset, limit int64) (ChainSensors, int64, error)
	FindChainSensorByID(ctx context.Context, loc, sensorID, dt string) (*ChainSensor, error)
	AddChainSensor(context.Context, ChainSensor) (string, error)
	DelChainSensor(context.Context, string) (int64, error)
	UpdChainSensor(context.Context, ChainSensor, bool, string) (int64, error)

	// bindings
	BindChainSensorEntrance(context.Context, BindingChainSensorEntrance) error
	BindChainSensorZone(context.Context, BindingChainSensorZone) error
	BindChainEntranceZone(context.Context, BindingChainEntranceZone) error
	BindChainEntranceStore(context.Context, BindingChainEntranceStore) error
	FindBindChainSensorZone(context.Context, string, string, int64, int64) (BindingsChainSensorZone, int64, error)
	DelChainBindSensorEntrance(context.Context, string, string) (int64, error)
	DelChainBindSensorZone(context.Context, string, string) (int64, error)
	DelChainBindEntranceZone(context.Context, string, string) (int64, error)
	DelChainBindEntranceStore(context.Context, string, string) (int64, error)
	UpdBindChainSensorEntrance(context.Context, BindingChainSensorEntrance, string, bool, string) error
	UpdBindChainSensorZone(context.Context, BindingChainSensorZone, string, bool, string) error
	UpdBindChainEntranceZone(context.Context, BindingChainEntranceZone, string, bool, string) error
	UpdBindChainEntranceStore(context.Context, BindingChainEntranceStore, string, bool, string) error
	//

	// malls
	FindMalls(ctx context.Context, loc, date string, offset, limit int64) (Malls, int64, error)
	FindMallByID(ctx context.Context, loc, layoutID, dt string) (*Mall, error)

	// mall/entrances
	FindMallEntrances(ctx context.Context,
		loc, date, layoutID, floorID, kind, entranceIDs string, offset, limit int64) (MallEntrances, int64, error)
	FindMallEntranceByID(ctx context.Context, loc, entranceID, dt string) (*MallEntrance, error)

	// mall/zones
	FindMallZones(ctx context.Context, loc, date, layoutID, kind, parentID, isOnline, isActive string,
		offset, limit int64) (MallZones, int64, error)
	FindMallZoneByID(ctx context.Context, loc, zoneID, dt string) (*MallZone, error)
	FindMallZonesByRenter(ctx context.Context, loc, date, renterID string) (MallZones, error)

	// mall/renters
	FindRenters(ctx context.Context, loc, date, layoutID, categorID, priceSegmentID, contract string,
		offset, limit int64) (Renters, int64, error)
	FindRenterByID(ctx context.Context, loc, renterID, dt string) (*Renter, error)

	// mall/devices
	FindMallDevices(ctx context.Context, loc, date, layoutID, kind, isActive, dcMode, sn, mode string,
		offset, limit int64) (MallDevices, int64, error)
	FindMallDeviceByID(ctx context.Context, loc, deviceID, dt string) (*MallDevice, error)
	FindMallDeviceDelays(context.Context, string, []string) (DeviceDelayDatas, error)

	// mall/sensors
	FindMallSensors(ctx context.Context, loc, date, layoutID, deviceID, kind string,
		offset, limit int64) (MallSensors, int64, error)
	FindMallSensorByID(ctx context.Context, loc, sensorID, dt string) (*MallSensor, error)

	// data
	// counting/inside
	FindZoneDataInsideDay(context.Context, *string, *time.Time) (DatasInside, error)
	FindZoneDataInsideNow(context.Context, *string) (DatasInside, error)
	FindZoneDataInsideRange(context.Context, time.Time, time.Time, *string, *time.Time) (DatasInside, error)
	// attendance
	FindChainStoresDataAttendance(ctx context.Context, from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filterIDs string) (StoresAttendance, error)
	FindChainZonesDataAttendance(ctx context.Context, from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filteredZoneIDs string) (ZonesAttendance, error)
	FindChainEntrancesDataAttendance(ctx context.Context,
		from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filteredEnterIDs string) (EntrancesAttendance, error)
	//
	FindMallZonesDataAttendance(ctx context.Context, from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filteredZoneIDs string) (ZonesAttendance, error)
	FindMallEntrancesDataAttendance(ctx context.Context,
		from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filteredEnterIDs string) (EntrancesAttendance, error)
	FindRenterDataAttendance(ctx context.Context, from time.Time, to time.Time, groupBy, layoutID, useRawData string,
		filteredList string) (RentersAttendance, error)
	// queue
	FindChainStoresDataQueue(ctx context.Context, from time.Time, to time.Time,
		storeID *string, groupBy, groupFunc string, window int) (StoresDataQueue, error)
	FindChainStoresDataQueueNow(ctx context.Context, storeID *string) (StoresDataQueue, error)
	FindChainZonesDataQueue(ctx context.Context,
		from time.Time, to time.Time, storeID, zoneID string) (ZonesDataQueue, error)
	FindChainZonesDataQueueNow(ctx context.Context, zoneID *string) (ZonesDataQueue, error)
	// zone states
	FindChainZonesStates(ctx context.Context, layoutID, storeID, zoneID string, from, to time.Time,
		offset, limit int64) (ZoneStates, int64, error)
	FindChainZonesStatesLast(ctx context.Context, layoutID, storeID, zoneID string) (ZoneStates, error)
	FindChainZoneStateAtTime(ctx context.Context, layoutID, zoneID string, t time.Time) (*ZoneState, error)
	// prediction data
	FindChainPredictionQueue(ctx context.Context, from time.Time, to time.Time, storeID *string,
		splitInterval time.Duration) (PredictionsQueue, error)
	// data_evaluation
	FindZoneDataEvaluation(ctx context.Context,
		from, to time.Time, layoutID, storeID, scBlockID, isFull string) (ZoneDataEvaluations, error)
	// behavior
	FindBehaviors(ctx context.Context, offset, limit int64) (Behaviors, int64, error)
	FindBehaviorByLayoutID(ctx context.Context, layoutID string) (*Behavior, error)
	UpdBehavior(ctx context.Context, layoutID string, bhv Behavior) (int64, error)

	// mall/chains common
	FindCrossesZoneEnter(ctx context.Context, zoneID, enterID string) (BindingsEntranceZone, error)

	// reports
	FindReports(ctx context.Context, layoutID, isSent string) (Reports, error)
	FindReportFiles(ctx context.Context,
		layoutID, reportID, isSent string, offset, limit int64) (ReportFiles, int64, error)
	FindReportFileByID(ctx context.Context, fileID string) (*ReportFile, error)

	// references
	GetReferences(ctx context.Context) ([]string, error)
	GetRefCategories(ctx context.Context) (reference.RefCategories, error)
	GetRefPriceSegments(ctx context.Context) (reference.RefPrices, error)
	GetRefKindZones(ctx context.Context) (reference.RefKindZones, error)
	GetRefKindEnters(ctx context.Context) (reference.RefKindEnters, error)

	// common
	GetEntities(ctx context.Context, loc, entitykey, parentkey, kind string, offset, limit int64) (Entities, int64, error)
	GetRepos(ctx context.Context) (map[string]LayoutRepo, error)
	GetSrvPortDB() string
	Health(ctx context.Context) error
	Scope() string
	Dest() string
}
