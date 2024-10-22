package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rorua/otelemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

var telemetry otelemetry.Telemetry

func main() {

	ctx := context.Background()

	cfg := otelemetry.Config{
		Service: otelemetry.Service{
			Name:      "my-service",
			Namespace: "my-namespace",
			Version:   "1.0.0",
		},
		Collector: otelemetry.Collector{
			Host: "0.0.0.0",
			Port: "4317",
		},
		ResourceOptions: []resource.Option{
			resource.WithHost(),
			//resource.WithProcess(),
			resource.WithTelemetrySDK(),
		},
	}

	var err error
	telemetry, err = otelemetry.New(cfg)
	if err != nil {
		panic(err)
	}

	defer telemetry.Shutdown(ctx)

	mux := http.NewServeMux()
	mux.Handle("/hello", otelhttp.NewHandler(http.HandlerFunc(handler), "/hello"))
	server := &http.Server{
		Addr:              ":7080",
		Handler:           mux,
		ReadHeaderTimeout: 20 * time.Second,
	}
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func handler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	ctx, span := telemetry.StartSpan(ctx, "handler")
	defer span.End()

	sleep := rng.Int63n(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	requestCount, err := telemetry.Meter().Float64Counter("example_request_count")
	if err != nil {
		panic(err)
	}

	requestCount.Add(ctx, 1)

	telemetry.Log().Info(ctx, "Request processed "+fmt.Sprintf("Sleep: %dms", sleep))

	span.AddEvent("sleep event", attribute.Int64("sleep", sleep))

	if _, err := w.Write([]byte(fmt.Sprintf("Sleep: %dms", sleep))); err != nil {
		http.Error(w, "write operation failed.", http.StatusInternalServerError)
		return
	}
}
