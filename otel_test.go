package otelemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTelemetryWithValidConfig(t *testing.T) {
	cfg := Config{
		Service: Service{
			Name:      "test-service",
			Namespace: "test-namespace",
			Version:   "1.0.0",
		},
		Collector: Collector{
			Host: "localhost",
			Port: "4317",
		},
		WithMetrics:    true,
		WithLogs:       true,
		WithStdoutLogs: true,
	}

	tel, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, tel)
}

func TestNewTelemetryWithInvalidCollector(t *testing.T) {
	cfg := Config{
		Service: Service{
			Name:      "test-service",
			Namespace: "test-namespace",
			Version:   "1.0.0",
		},
		Collector: Collector{
			Host: "",
			Port: "",
		},
	}

	_, err := New(cfg)
	assert.Error(t, err)
}

func TestTelemetryShutdown(t *testing.T) {
	cfg := Config{
		Service: Service{
			Name:      "test-service",
			Namespace: "test-namespace",
			Version:   "1.0.0",
		},
		Collector: Collector{
			Host: "localhost",
			Port: "4317",
		},
	}

	tel, err := New(cfg)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = tel.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestTelemetryLog(t *testing.T) {
	cfg := Config{
		Service: Service{
			Name:      "test-service",
			Namespace: "test-namespace",
			Version:   "1.0.0",
		},
		Collector: Collector{
			Host: "localhost",
			Port: "4317",
		},
		WithLogs: true,
	}

	tel, err := New(cfg)
	assert.NoError(t, err)

	logger := tel.Log()
	assert.NotNil(t, logger)
}
