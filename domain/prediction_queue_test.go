package domain

import (
	"testing"
	"time"
)

func TestAddPointPrediction(t *testing.T) {
	t0, _ := time.Parse(time.RFC3339, "2020-10-20T12:00:00+03:00")
	points := genPredPoints(t0, time.Minute, 10)
	pred := make(PredictionsQueue, 0, len(points))
	fillInterval := 15 * time.Second
	for _, p := range points {
		pred.AddPoint("1", fillInterval, p)
	}
	for _, pq := range pred {
		t1 := t0
		for t1.Before(t0.Add(9 * time.Minute)) {
			exists := false
			for _, p := range pq.Points {
				if p.Time.Equal(t1) {
					exists = true
				}
			}
			if !exists {
				t.Errorf("point at time=%s not found in %v", t1, pq.Points)
			}
			t1 = t1.Add(fillInterval)
		}
	}
}

func genPredPoints(t0 time.Time, step time.Duration, n int) []PredictionDataPoint {
	result := make([]PredictionDataPoint, n)
	for i := 0; i < n; i++ {
		result[i] = PredictionDataPoint{
			Time:               t0,
			CashIncomeFlow:     1.0,
			CashIncomeSum:      2.0,
			TotalCustomerCount: 3.0,
			QueueLength:        []float64{1.0, 0.0},
		}
		t0 = t0.Add(step)
	}
	return result
}
