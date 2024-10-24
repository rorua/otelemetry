package otelemetry

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func resourceOpts(cfg Config, opts []resource.Option) []resource.Option {
	attrOption := resource.WithAttributes(
		// the telemetry name used to display traces in backends
		semconv.ServiceNameKey.String(cfg.Service.Name),
		semconv.ServiceNamespaceKey.String(cfg.Service.Namespace),
		semconv.ServiceVersionKey.String(cfg.Service.Version),
	)

	if len(opts) == 0 {
		return []resource.Option{
			attrOption,
		}
	}

	opts = append(opts, attrOption)

	return opts
}
