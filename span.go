package otelemetry

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	Span() trace.Span // think about naming
	AddEvent(name string, kv ...attribute.KeyValue)
	AddErrorEvent(name string, err error, kv ...attribute.KeyValue)
	SetAttribute(kv ...attribute.KeyValue)
	End(opts ...trace.SpanEndOption) // aka Finish()
	RecordError(err error, kv ...attribute.KeyValue)
	TraceID() string
	SpanID() string
}

type otelspan struct {
	span trace.Span
}

func (s *otelspan) Span() trace.Span {
	return s.span
}

func (s *otelspan) AddEvent(name string, kv ...attribute.KeyValue) {
	s.span.AddEvent(name, trace.WithAttributes(kv...))
}

func (s *otelspan) AddErrorEvent(name string, err error, kv ...attribute.KeyValue) {
	s.span.SetStatus(codes.Error, err.Error())
	kv = append(kv, attribute.String("error.message", err.Error()))
	kv = append(kv, attribute.String("error.type", fmt.Sprintf("%T", err)))
	s.span.AddEvent(name, trace.WithAttributes(kv...))
}

func (s *otelspan) SetAttribute(kv ...attribute.KeyValue) {
	s.span.SetAttributes(kv...)
}

func (s *otelspan) RecordError(err error, kv ...attribute.KeyValue) {
	s.span.SetStatus(codes.Error, err.Error())
	s.span.RecordError(err, trace.WithAttributes(kv...))
}

func (s *otelspan) End(opts ...trace.SpanEndOption) {
	s.span.End(opts...)
}

func (s *otelspan) TraceID() string {
	return s.span.SpanContext().TraceID().String()
}

func (s *otelspan) SpanID() string {
	return s.span.SpanContext().SpanID().String()
}
