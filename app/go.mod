module github.com/Amir23156/cloud-native-observability/app

go 1.22

require (
	github.com/gorilla/mux v1.8.1
	github.com/prometheus/client_golang v1.19.1
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/sdk/resource v1.28.0
	go.opentelemetry.io/otel/semconv v1.24.0
	google.golang.org/grpc v1.65.0
)
