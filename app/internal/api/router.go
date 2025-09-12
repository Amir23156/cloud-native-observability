package api

import "github.com/gorilla/mux"
import "github.com/Amir23156/cloud-native-observability/app/internal/metrics"

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/orders", orders).Methods("GET")
	r.Handle("/metrics", metrics.Handler()).Methods("GET")
	return r
}
