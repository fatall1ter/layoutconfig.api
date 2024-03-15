package domain

import "time"

//go:generate easyjson -all chain_data.go

// DataInside contains data about number of people inside zone at the moment
type DataInside struct {
	ZoneID string     `json:"zone_id"`
	Points DataPoints `json:"points"`
}

//easyjson:json
type DatasInside []DataInside

func (di *DatasInside) AddPoint(zoneID string, point DataPoint) {
	for i, edi := range *di {
		if edi.ZoneID == zoneID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// zone not found
	newDataInside := DataInside{
		ZoneID: zoneID,
		Points: DataPoints{point},
	}
	*di = append(*di, newDataInside)
}

type DataPoint struct {
	Time  time.Time `json:"time"`
	Value int32     `json:"value"`
}

type DelayPoint struct {
	Value int `json:"value"`
}

//easyjson:json
// DataPoints slice of measurements of DataPoint
type DataPoints []DataPoint

// DATA QUEUE

// QueueDataPoint queue measurement in specified moment of the time
type QueueDataPoint struct {
	// datetime of measurement
	Time time.Time `json:"time"`
	// count of people in the queue total (include service channels only in the online mode)
	Value int32 `json:"value"`
	// count of people in the queue total (include service channels in the offline and online mode)
	ValueTotal int32 `json:"value_total"`
	// number of service channels in the online mode at measurement moment
	CountChannels int `json:"count_channels"`
	// verified number of service channels in the online mode at measurement moment
	CountChannelsVerified int `json:"count_channels_verified"`
	// flag existsing verified number of service channels
	HasVerified bool `json:"has_verified"`
	// how many "people" income to queue from prev interval
	CashIncomeFlow int32 `json:"cash_income_flow"`
	// summa of the income flow values by window interval
	SumCashIncomeFlowByWindow int32 `json:"sum_cash_income_flow_by_window"`
}

//easyjson:json
type QueueDataPoints []QueueDataPoint

type StoreDataQueue struct {
	StoreID string          `json:"store_id"`
	Points  QueueDataPoints `json:"points"`
}

//easyjson:json
type StoresDataQueue []StoreDataQueue

func (di *StoresDataQueue) AddPoint(storeID string, point QueueDataPoint) {
	for i, edi := range *di {
		if edi.StoreID == storeID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// store not found
	newStoreDataQueue := StoreDataQueue{
		StoreID: storeID,
		Points:  QueueDataPoints{point},
	}
	*di = append(*di, newStoreDataQueue)
}

type ZoneDataQueue struct {
	ZoneID string          `json:"zone_id"`
	Points QueueDataPoints `json:"points"`
}

//easyjson:json
type ZonesDataQueue []ZoneDataQueue

func (di *ZonesDataQueue) AddPoint(zoneID string, point QueueDataPoint) {
	for i, edi := range *di {
		if edi.ZoneID == zoneID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// zone not found
	newZoneDataQueue := ZoneDataQueue{
		ZoneID: zoneID,
		Points: QueueDataPoints{point},
	}
	*di = append(*di, newZoneDataQueue)
}

// Attendance

// AttendanceDataPoint simple attendance data unit
type AttendanceDataPoint struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	SumIn     int32     `json:"sum_in"`
	SumOut    int32     `json:"sum_out"`
	PassingBy int32     `json:"passing_by"`
}

//easyjson:json AttendanceDataPoints slice of measurements of AttendanceDataPoint
type AttendanceDataPoints []AttendanceDataPoint

type ZoneAttendance struct {
	ZoneID string               `json:"zone_id"`
	Points AttendanceDataPoints `json:"points"`
}

//easyjson:json
type ZonesAttendance []ZoneAttendance

func (di *ZonesAttendance) AddPoint(zoneID string, point AttendanceDataPoint) {
	for i, edi := range *di {
		if edi.ZoneID == zoneID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// zone not found
	newDataPoint := ZoneAttendance{
		ZoneID: zoneID,
		Points: AttendanceDataPoints{point},
	}
	*di = append(*di, newDataPoint)
}

type StoreAttendance struct {
	StoreID string               `json:"store_id"`
	Points  AttendanceDataPoints `json:"points"`
}

//easyjson:json
type StoresAttendance []StoreAttendance

func (di *StoresAttendance) AddPoint(storeID string, point AttendanceDataPoint) {
	for i, edi := range *di {
		if edi.StoreID == storeID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// store not found
	newDataPoint := StoreAttendance{
		StoreID: storeID,
		Points:  AttendanceDataPoints{point},
	}
	*di = append(*di, newDataPoint)
}

type EntranceAttendance struct {
	EntranceID string               `json:"entrance_id"`
	Points     AttendanceDataPoints `json:"points"`
}

//easyjson:json
type EntrancesAttendance []EntranceAttendance

func (di *EntrancesAttendance) AddPoint(entranceID string, point AttendanceDataPoint) {
	for i, edi := range *di {
		if edi.EntranceID == entranceID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// entrance not found
	newDataPoint := EntranceAttendance{
		EntranceID: entranceID,
		Points:     AttendanceDataPoints{point},
	}
	*di = append(*di, newDataPoint)
}

type RenterAttendance struct {
	RenterID string               `json:"renter_id"`
	Points   AttendanceDataPoints `json:"points"`
}

//easyjson:json
type RentersAttendance []RenterAttendance

func (di *RentersAttendance) AddPoint(renterID string, point AttendanceDataPoint) {
	for i, edi := range *di {
		if edi.RenterID == renterID {
			edi.Points = append(edi.Points, point)
			(*di)[i] = edi
			return
		}
	}
	// entrance not found
	newDataPoint := RenterAttendance{
		RenterID: renterID,
		Points:   AttendanceDataPoints{point},
	}
	*di = append(*di, newDataPoint)
}
