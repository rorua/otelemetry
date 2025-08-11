package otelemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddBaggageItem(t *testing.T) {
	ctx := context.Background()
	ctx = AddBaggageItem(ctx, "user", "alice")
	ctx = AddBaggageItem(ctx, "tenant", "acme")

	b := GetBaggage(ctx)
	assert.Equal(t, "alice", b.Member("user").Value())
	assert.Equal(t, "acme", b.Member("tenant").Value())
}

func TestAddBaggageItems(t *testing.T) {
	ctx := context.Background()
	ctx = AddBaggageItem(ctx, "user", "alice")
	ctx = AddBaggageItems(ctx, map[string]string{"tenant": "acme", "role": "admin"})

	b := GetBaggage(ctx)
	assert.Equal(t, "alice", b.Member("user").Value())
	assert.Equal(t, "acme", b.Member("tenant").Value())
	assert.Equal(t, "admin", b.Member("role").Value())
}
