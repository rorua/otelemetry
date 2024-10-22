package otelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func newLoggerProvider(ctx context.Context, otelAgentAddr string, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	//f, err := os.OpenFile("otel_log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	return nil, err
	//}
	//
	//
	//exporter, err := stdoutlog.New(
	//	stdoutlog.WithWriter(f),
	//	//stdoutlog.WithPrettyPrint(),
	//)
	//if err != nil {
	//	return nil, err
	//}

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(otelAgentAddr),
	)
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	return provider, nil
}

const (
	Debug = "DEBUG"
	Info  = "INFO"
	Warn  = "WARN"
	Error = "ERROR"
	Fatal = "FATAL"
)

type Log interface {
	Debug(ctx context.Context, msg string, kv ...log.KeyValue)
	Info(ctx context.Context, msg string, kv ...log.KeyValue)
	Warning(ctx context.Context, msg string, kv ...log.KeyValue)
	Error(ctx context.Context, msg string, kv ...log.KeyValue)
	Fatal(ctx context.Context, msg string, kv ...log.KeyValue)
}

type otellog struct {
	log log.Logger
}

func (l *otellog) Debug(ctx context.Context, msg string, kv ...log.KeyValue) {
	record := getRecord(msg, log.SeverityDebug, Debug, kv...)
	l.log.Emit(ctx, record)
}

func (l *otellog) Info(ctx context.Context, msg string, kv ...log.KeyValue) {
	record := getRecord(msg, log.SeverityInfo, Info, kv...)
	l.log.Emit(ctx, record)
}

func (l *otellog) Warning(ctx context.Context, msg string, kv ...log.KeyValue) {
	record := getRecord(msg, log.SeverityWarn, Warn, kv...)
	l.log.Emit(ctx, record)
}

func (l *otellog) Error(ctx context.Context, msg string, kv ...log.KeyValue) {
	record := getRecord(msg, log.SeverityError, Error, kv...)
	l.log.Emit(ctx, record)
}

func (l *otellog) Fatal(ctx context.Context, msg string, kv ...log.KeyValue) {
	record := getRecord(msg, log.SeverityFatal, Fatal, kv...)
	l.log.Emit(ctx, record)
}

func getRecord(msg string, severity log.Severity, sevName string, kv ...log.KeyValue) log.Record {
	var record log.Record
	record.SetBody(log.StringValue(msg))
	record.SetSeverity(severity)
	record.SetSeverityText(sevName)
	record.SetTimestamp(time.Now())
	record.AddAttributes(kv...)
	return record
}
