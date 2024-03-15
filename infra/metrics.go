package infra

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	common = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_up",
			Help: "Common metric of state of available service and subservices",
		},
		[]string{"scope", "destination", "version", "githash", "build"},
	)

	api_websocket_connections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_websocket_connections",
			Help: "WebSocket connections by url",
		},
		[]string{"path"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_request_duration_seconds",
			Help: "Histogram of the api req-resp durations in seconds by backets",
		},
		[]string{"url", "code", "method"},
	)
)
