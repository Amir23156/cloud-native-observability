package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ordersTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders placed",
		},
	)
)

// Register must be called once at startup.
func Register() {
	prometheus.MustRegister(ordersTotal)
}

// IncOrders increments the orders counter.
func IncOrders() { ordersTotal.Inc() }

// Handler exposes Prometheus metrics at /metrics.
func Handler() http.Handler { return promhttp.Handler() }
