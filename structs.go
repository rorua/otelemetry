package otelemetry

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Service
	Collector
	WithTraces      bool
	WithMetrics     bool
	WithLogs        bool
	ResourceOptions []resource.Option
	TracerOptions   TracerOptions
	LoggerOptions   LoggerOptions
	MetricOptions   MetricOptions
}

type Service struct {
	Name      string
	Namespace string
	Version   string
}

type Collector struct {
	Host string
	Port string
}

type LoggerOptions struct {
	ExporterOption []otlploggrpc.Option
	ProviderOption []sdklog.LoggerProviderOption

	LoggerOption []log.LoggerOption
}

type TracerOptions struct {
	ClientOption             []otlptracegrpc.Option
	ProviderOption           []sdktrace.TracerProviderOption
	BatchSpanProcessorOption []sdktrace.BatchSpanProcessorOption

	TracerOption []trace.TracerOption
}

type MetricOptions struct {
	ExporterOptions []otlpmetricgrpc.Option
	ProviderOptions []sdkmetric.Option

	MeterOptions []metric.MeterOption
}
