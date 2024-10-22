package otelemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
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

	OnStart func(context.Context) error
	OnStop  func(context.Context) error
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
	meterProvider, err := newMeterProvider(ctx, otelAgentAddr, res)
	handleErr(err, "failed to create the collector metric exporter or provider")
	otel.SetMeterProvider(meterProvider)

	// traces
	tracerProvider, traceExporter, err := newTraceProvider(ctx, otelAgentAddr, res)
	handleErr(err, "failed to create the collector trace exporter or provider")
	_ = traceExporter

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	// logger provider
	loggerProvider, err := newLoggerProvider(ctx, otelAgentAddr, res)
	handleErr(err, "failed to create the logger provider")

	// Set the logger provider globally
	global.SetLoggerProvider(loggerProvider)

	//params.Lifecycle.Append(
	//	fx.Hook{
	//		OnStart: func(ctx context.Context) error {
	//			params.Logger.Info("Telemetry started")
	//			return nil
	//		},
	//		OnStop: func(ctx context.Context) error {
	//			cxt, cancel := context.WithTimeout(ctx, time.Second)
	//			defer cancel()
	//
	//			// pushes any last exports to the receiver
	//			if err := traceExporter.Shutdown(cxt); err != nil {
	//				otel.Handle(err)
	//			}
	//
	//			if err := meterProvider.Shutdown(cxt); err != nil {
	//				otel.Handle(err)
	//			}
	//
	//			if err := loggerProvider.Shutdown(cxt); err != nil {
	//				otel.Handle(err)
	//			}
	//
	//			params.Logger.Debug("Telemetry stopped")
	//			return err
	//		},
	//	},
	//)

	return &telemetry{
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		loggerProvider: loggerProvider,

		tracer: tracerProvider.Tracer(serviceName),
		meter:  meterProvider.Meter(serviceName),
		logger: loggerProvider.Logger(serviceName),

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
