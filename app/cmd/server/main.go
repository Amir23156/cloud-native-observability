package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Amir23156/cloud-native-observability/app/internal/api"
	"github.com/Amir23156/cloud-native-observability/app/internal/metrics"
	"github.com/Amir23156/cloud-native-observability/app/internal/observability"
)

func main() {
	ctx := context.Background()

	// Tracing
	shutdown, err := observability.InitTracer(ctx)
	if err != nil {
		log.Fatalf("tracing init failed: %v", err)
	}
	defer func() { _ = shutdown(ctx) }()

	// Metrics
	metrics.Register()

	// Router
	r := api.NewRouter()

	addr := ":" + env("PORT", "5000")
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Signal handling for graceful shutdown ............
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
