package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	}, []string{"method", "path"})

	PVZCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "business_pvz_created_total",
		Help: "Total number of created PVZ",
	})

	ReceptionsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "business_receptions_created_total",
		Help: "Total number of created receptions",
	})

	ProductsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "business_products_added_total",
		Help: "Total number of added products",
	})
)
