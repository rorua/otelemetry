package main

import (
	"context"

	"github.com/rorua/otelemetry"
	"go.opentelemetry.io/otel/metric"
)

func main() {

	ctx := context.Background()

	cfg := otelemetry.Config{}
	cfg.Service.Name = "my-service"
	cfg.Service.Namespace = "my-namespace"
	cfg.Service.Version = "1.0.0"
	cfg.Collector.Host = "0.0.0.0"
	cfg.Collector.Port = "4317"

	//This is a placeholder for the main function.
	telemetry, err := otelemetry.New(cfg)
	if err != nil {
		panic(err)
	}

	defer telemetry.Shutdown(ctx)

	ctx, span := telemetry.StartSpan(ctx, "main")
	defer span.End()

	// Do something with telemetry
	telemetry.Log().Info(ctx, "hello world")

	requestCount, err := telemetry.Meter().Int64Counter(
		"example_counter",
		metric.WithDescription("The number of requests received"),
	)

	requestCount.Add(ctx, 1)

}
