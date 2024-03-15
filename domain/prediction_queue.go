package domain

import "time"

//go:generate easyjson -all prediction_queue.go

// PredictionQueue recommendations on the need to open service channels to prevent queuing
type PredictionQueue struct {
	StoreID string               `json:"store_id"`
	Points  PredictionDataPoints `json:"points"`
}

//easyjson:json
type ListQueueLength []float64

type PredictionDataPoint struct {
	Time                       time.Time       `json:"time,omitempty"`
	CashIncomeFlow             float64         `json:"cash_income_flow"`
	CashIncomeSum              float64         `json:"cash_income_sum"`
	TotalCustomerCount         float64         `json:"total_customer_count"`
	QueueLength                ListQueueLength `json:"queue_length"`
	RecommendedCheckoutsNumber int             `json:"recommended_checkouts_number"`
	PredictionQueueLength      float64         `json:"prediction_queue_length"`
}

//easyjson:json
type PredictionDataPoints []PredictionDataPoint

//easyjson:json
type PredictionsQueue []PredictionQueue

func (pqs *PredictionsQueue) AddPoint(storeID string, fillInterval time.Duration, point PredictionDataPoint) {
	for i, pq := range *pqs {
		if pq.StoreID == storeID {
			if len(pq.Points) > 0 { // fill by interval
				pre := pq.Points[len(pq.Points)-1]
				t := pre.Time.Add(fillInterval)
				for t.Before(point.Time) {
					pre.Time = t
					pq.Points = append(pq.Points, pre)
					t = t.Add(fillInterval)
				}
			}
			pq.Points = append(pq.Points, point)
			(*pqs)[i] = pq
			return
		}
	}
	// zone not found
	mewPQ := PredictionQueue{
		StoreID: storeID,
		Points:  PredictionDataPoints{point},
	}
	*pqs = append(*pqs, mewPQ)
}
