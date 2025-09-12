package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Amir23156/cloud-native-observability/app/internal/api"
	"github.com/Amir23156/cloud-native-observability/app/internal/metrics"
	"github.com/Amir23156/cloud-native-observability/app/internal/observability"
)

func main() {
	ctx := context.Background()

	// Init tracing (OTLP -> OTel Collector -> Tempo)
	shutdown, err := observability.InitTracer(ctx)
	if err != nil {
		log.Fatalf("tracing init failed: %v", err)
	}
	defer func() { _ = shutdown(ctx) }()

	// Register Prometheus metrics
	metrics.Register()

	// Build router
	r := api.NewRouter()

	// Serve
	addr := ":" + env("PORT", "5000")
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
