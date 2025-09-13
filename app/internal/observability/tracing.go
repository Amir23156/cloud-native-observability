package observability

import (
    "context"
    "log"
    "os"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/semconv/v1.24.0"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// InitTracer configures OTLP exporter -> OTel Collector (Tempo backend).
// If no endpoint is provided or connection fails, tracing is disabled but the app still runs.
func InitTracer(ctx context.Context) (func(context.Context) error, error) {
    endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if endpoint == "" {
        log.Println("tracing disabled: no OTLP endpoint configured")
        return func(context.Context) error { return nil }, nil
    }

    dctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(
        dctx, endpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        log.Printf("tracing disabled: failed to connect to %s: %v", endpoint, err)
        return func(context.Context) error { return nil }, nil
    }

    exp, err := otlptrace.New(dctx,
        otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn)),
    )
    if err != nil {
        log.Printf("tracing disabled: exporter init error: %v", err)
        return func(context.Context) error { return nil }, nil
    }

    res, _ := resource.New(dctx,
        resource.WithAttributes(
            semconv.ServiceName("go-orders-api"),
            semconv.DeploymentEnvironment(os.Getenv("ENVIRONMENT")),
        ),
    )

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(res),
    )

    otel.SetTracerProvider(tp)
    log.Printf("tracing initialized: exporting to %s", endpoint)

    return tp.Shutdown, nil
}
