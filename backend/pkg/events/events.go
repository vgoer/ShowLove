// Package events provides NATS JetStream publish/subscribe encapsulation.
package events

import (
	"context"
	"encoding/json"
)

// Event represents a domain event published across services.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Publisher sends events to the message broker.
type Publisher interface {
	// Publish sends an event to the given subject.
	Publish(ctx context.Context, subject string, event Event) error
}

// Subscriber receives events from the message broker.
type Subscriber interface {
	// Subscribe listens for events on the given subject.
	Subscribe(ctx context.Context, subject string, handler func(Event) error) error
}

// NoOpPubSub is a no-op implementation for testing and development.
type NoOpPubSub struct {
	Subscribers map[string][]func(Event) error
}

// NewNoOpPubSub creates a new no-op pub/sub.
func NewNoOpPubSub() *NoOpPubSub {
	return &NoOpPubSub{
		Subscribers: make(map[string][]func(Event) error),
	}
}

// Publish implements Publisher. Delivers events to registered local subscribers.
func (n *NoOpPubSub) Publish(ctx context.Context, subject string, event Event) error {
	for _, handler := range n.Subscribers[subject] {
		if err := handler(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe implements Subscriber. Registers a local handler.
func (n *NoOpPubSub) Subscribe(ctx context.Context, subject string, handler func(Event) error) error {
	n.Subscribers[subject] = append(n.Subscribers[subject], handler)
	return nil
}
