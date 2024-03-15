package domain

import (
	"reflect"
	"testing"
	"time"
)

var (
	di1 DataInside = DataInside{
		ZoneID: "111",
		Points: DataPoints{
			DataPoint{
				Time:  makeTime("2020-01-01T10:00:00+03:00"),
				Value: 10,
			},
			DataPoint{
				Time:  makeTime("2020-01-01T11:00:00+03:00"),
				Value: 11,
			},
		},
	}
	di1_2 DataInside = DataInside{
		ZoneID: "111",
		Points: DataPoints{
			DataPoint{
				Time:  makeTime("2020-01-01T10:00:00+03:00"),
				Value: 10,
			},
			DataPoint{
				Time:  makeTime("2020-01-01T11:00:00+03:00"),
				Value: 11,
			},
			DataPoint{
				Time:  makeTime("2020-01-01T12:00:00+03:00"),
				Value: 12,
			},
		},
	}
	di2 DataInside = DataInside{
		ZoneID: "222",
		Points: DataPoints{
			DataPoint{
				Time:  makeTime("2020-01-01T10:00:00+03:00"),
				Value: 10,
			},
			DataPoint{
				Time:  makeTime("2020-01-01T11:00:00+03:00"),
				Value: 11,
			},
		},
	}
	di3 DataInside = DataInside{
		ZoneID: "333",
		Points: DataPoints{
			DataPoint{
				Time:  makeTime("2020-01-01T10:00:00+03:00"),
				Value: 10,
			},
			DataPoint{
				Time:  makeTime("2020-01-01T11:00:00+03:00"),
				Value: 11,
			},
		},
	}
	dss     DatasInside = DatasInside{di1, di2}
	expDss  DatasInside = DatasInside{di1_2, di2}
	expDssN DatasInside = DatasInside{di1_2, di2, di3}
)

func makeTime(tsi string) time.Time {
	ts, err := time.Parse(time.RFC3339, tsi)
	if err != nil {
		ts, _ = time.Parse(time.RFC3339, "2020-04-01T00:00:00+03:00")
		return ts
	}
	return ts
}

func TestAddPoint(t *testing.T) {
	point := DataPoint{
		Time:  makeTime("2020-01-01T12:00:00+03:00"),
		Value: 12,
	}
	zoneID := "111"
	dss.AddPoint(zoneID, point)
	if !reflect.DeepEqual(dss, expDss) {
		t.Errorf("expected dss after add\n%+v\nbut got:\n%+v\n", expDss, dss)
	}
	// add not exists zone points
	point2 := DataPoint{
		Time:  makeTime("2020-01-01T10:00:00+03:00"),
		Value: 10,
	}
	point3 := DataPoint{
		Time:  makeTime("2020-01-01T11:00:00+03:00"),
		Value: 11,
	}
	zoneID2 := "333"
	dss.AddPoint(zoneID2, point2)
	dss.AddPoint(zoneID2, point3)
	if !reflect.DeepEqual(dss, expDssN) {
		t.Errorf("expected dss after add\n%+v\nbut got:\n%+v\n", expDssN, dss)
	}
}
