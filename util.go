package otelemetry

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func handleErr(err error, s string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", s, err))
	}
}

func Attribute(k string, v any) attribute.KeyValue {
	return parseAttribute(k, v)
}

func parseAttribute(key string, value any) attribute.KeyValue {
	var attr attribute.KeyValue
	switch v := value.(type) {
	case string:
		attr = attribute.String(key, v)
	case []string:
		attr = attribute.StringSlice(key, v)
	case fmt.Stringer:
		attr = attribute.Stringer(key, v)
	case int:
		attr = attribute.Int(key, v)
	case []int:
		attr = attribute.IntSlice(key, v)
	case int64:
		attr = attribute.Int64(key, v)
	case []int64:
		attr = attribute.Int64Slice(key, v)
	case bool:
		attr = attribute.Bool(key, v)
	case []bool:
		attr = attribute.BoolSlice(key, v)
	case float64:
		attr = attribute.Float64(key, v)
	case []float64:
		attr = attribute.Float64Slice(key, v)
	default:
		attr = attribute.String(key, fmt.Sprintf("%#v", v))
	}
	return attr
}

func LogAttribute(k string, v any) log.KeyValue {
	return parseLogAttribute(k, v)
}

func parseLogAttribute(key string, value any) log.KeyValue {
	var attr log.KeyValue
	switch v := value.(type) {
	case string:
		attr = log.String(key, v)
	case int:
		attr = log.Int(key, v)
	case int64:
		attr = log.Int64(key, v)
	case bool:
		attr = log.Bool(key, v)
	case float64:
		attr = log.Float64(key, v)
	default:
		attr = log.String(key, fmt.Sprintf("%+v", v))
	}
	return attr
}
