package api

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/Amir23156/cloud-native-observability/app/internal/metrics"
	"go.opentelemetry.io/otel"
)

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"status": "healthy"})
}

func orders(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("go-orders-api")
	_, span := tr.Start(r.Context(), "orders_handler")
	defer span.End()

	metrics.IncOrders()
	writeJSON(w, map[string]any{
		"message":          "Order placed successfully!",
		"orders_processed": rand.Intn(100) + 1,
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
