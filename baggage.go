package otelemetry

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

// GetBaggage retrieves the baggage from the context.
func GetBaggage(ctx context.Context) baggage.Baggage {
	return baggage.FromContext(ctx)
}

// AddBaggageItem adds a key-value pair to the baggage.
func AddBaggageItem(ctx context.Context, key, value string) context.Context {
	b := baggage.FromContext(ctx)
	m, _ := baggage.NewMember(key, value)
	b, _ = b.SetMember(m)
	return baggage.ContextWithBaggage(ctx, b)
}

// AddBaggageItems adds multiple key-value pairs to the baggage.
func AddBaggageItems(ctx context.Context, items map[string]string) context.Context {
	b := baggage.FromContext(ctx)
	for key, value := range items {
		m, _ := baggage.NewMember(key, value)
		b, _ = b.SetMember(m)
	}
	return baggage.ContextWithBaggage(ctx, b)
}

// GetBaggageItem retrieves the value of a baggage item by key.
func GetBaggageItem(ctx context.Context, key string) string {
	b := baggage.FromContext(ctx)
	member := b.Member(key)
	return member.Value()
}

// RemoveBaggageItem removes a key-value pair from the baggage.
func RemoveBaggageItem(ctx context.Context, key string) context.Context {
	b := baggage.FromContext(ctx)
	members := b.Members()
	var newMembers []baggage.Member
	for _, member := range members {
		if member.Key() != key {
			newMembers = append(newMembers, member)
		}
	}
	newBaggage, _ := baggage.New(newMembers...)
	return baggage.ContextWithBaggage(ctx, newBaggage)
}
