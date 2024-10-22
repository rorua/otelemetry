package otelemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Service struct {
		Name      string
		Namespace string
		Version   string
	}

	Collector struct {
		Host string
		Port string
	}

	ResourceOptions []resource.Option

	TracerOptions TracerOptions
	LoggerOption  LoggerOptions
	MetricOptions MetricOptions
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

func New(cfg Config) (Telemetry, error) {

	var (
		ctx               = context.Background()
		serviceName       = cfg.Service.Name
		otelCollectorHost = cfg.Collector.Host
		otelCollectorPort = cfg.Collector.Port
	)

	// otel collector OTEL_COLLECTOR_HOST:OTEL_COLLECTOR_PORT_GRPC
	otelAgentAddr := fmt.Sprintf("%s:%s", otelCollectorHost, otelCollectorPort)

	// resource
	res, err := resource.New(ctx, resourceOpts(cfg, cfg.ResourceOptions)...)
	handleErr(err, "failed to create resource")

	// metrics
	meterProvider, err := newMeterProvider(ctx, otelAgentAddr, res, cfg.MetricOptions)
	handleErr(err, "failed to create the collector metric exporter or provider")
	otel.SetMeterProvider(meterProvider)

	// traces
	tracerProvider, err := newTraceProvider(ctx, otelAgentAddr, res, cfg.TracerOptions)
	handleErr(err, "failed to create the collector trace exporter or provider")

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	// logger provider
	loggerProvider, err := newLoggerProvider(ctx, otelAgentAddr, res)
	handleErr(err, "failed to create the logger provider")

	// Set the logger provider globally
	global.SetLoggerProvider(loggerProvider)

	return &telemetry{
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		loggerProvider: loggerProvider,

		tracer: tracerProvider.Tracer(serviceName, cfg.TracerOptions.TracerOption...),
		meter:  meterProvider.Meter(serviceName, cfg.MetricOptions.MeterOptions...),
		logger: loggerProvider.Logger(serviceName, cfg.LoggerOption.LoggerOption...),

		serviceName: serviceName,
	}, nil
}

func handleErr(err error, s string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", s, err))
	}
}

type Telemetry interface {
	StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span)
	SpanFromContext(ctx context.Context) Span

	Log() Log

	Shutdown(ctx context.Context) error

	providers
}

type providers interface {
	Tracer() trace.Tracer
	Meter() metric.Meter
	Logger() log.Logger
}

type telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	loggerProvider *sdklog.LoggerProvider

	tracer trace.Tracer
	meter  metric.Meter
	logger log.Logger

	serviceName string
}

func (t *telemetry) Tracer() trace.Tracer {
	return t.tracer
}

func (t *telemetry) Meter() metric.Meter {
	return t.meter
}

func (t *telemetry) Logger() log.Logger {
	return t.logger
}

func (t *telemetry) Shutdown(ctx context.Context) error {
	cxt, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// pushes any last exports to the receiver
	if err := t.tracerProvider.Shutdown(cxt); err != nil {
		otel.Handle(err)
	}

	if err := t.meterProvider.Shutdown(cxt); err != nil {
		otel.Handle(err)
	}

	if err := t.loggerProvider.Shutdown(cxt); err != nil {
		otel.Handle(err)
	}

	return nil
}

func (t *telemetry) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	ctx, span := t.Tracer().Start(ctx, name, opts...)
	return ctx, &otelspan{span: span}
}

func (t *telemetry) SpanFromContext(ctx context.Context) Span {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	return &otelspan{span: span}
}

func (t *telemetry) Log() Log {
	return &otellog{log: t.Logger()}
}

func resourceOpts(cfg Config, opts []resource.Option) []resource.Option {
	attrOption := resource.WithAttributes(
		// the telemetry name used to display traces in backends
		semconv.ServiceNameKey.String(cfg.Service.Name),
		semconv.ServiceNamespaceKey.String(cfg.Service.Version),
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
