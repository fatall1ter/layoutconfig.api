package domain

import "time"

type DevConfig struct {
	SN      string     `json:"sn"`
	Sensors DevSensors `json:"sensors"`
}

func (dc *DevConfig) GetSensorByID(id string) *DevSensor {
	for _, s := range dc.Sensors {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

type DevSensors []DevSensor

type DevSensor struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// tracks

// Track describe position of the customer in vision zone of the counting device
type Track struct {
	LayoutID   string    `json:"layout_id"`
	StoreID    string    `json:"store_id"`
	DeviceID   string    `json:"device_id"`
	CustomerID string    `json:"customer_id"`
	TrackTime  time.Time `json:"track_time"`
	XPos       float64   `json:"x"`
	YPos       float64   `json:"y"`
	Height     float64   `json:"h"`
	Creator    string    `json:"creator"`
	CreatedAt  time.Time `json:"created_at"`
	Comment    string    `json:"comment"`
}

//easyjson:json
type Tracks []Track
