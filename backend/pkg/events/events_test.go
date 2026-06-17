package events

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoOpPubSub_PublishSubscribe(t *testing.T) {
	ps := NewNoOpPubSub()

	received := make(chan Event, 1)
	err := ps.Subscribe(context.Background(), "post.created", func(e Event) error {
		received <- e
		return nil
	})
	require.NoError(t, err)

	payload, _ := json.Marshal(map[string]string{"post_id": "123"})
	event := Event{
		Type:    "post.created",
		Payload: payload,
	}

	err = ps.Publish(context.Background(), "post.created", event)
	require.NoError(t, err)

	select {
	case got := <-received:
		assert.Equal(t, "post.created", got.Type)
		var data map[string]string
		json.Unmarshal(got.Payload, &data)
		assert.Equal(t, "123", data["post_id"])
	default:
		t.Fatal("expected to receive event")
	}
}

func TestNoOpPubSub_MultipleSubscribers(t *testing.T) {
	ps := NewNoOpPubSub()

	count := 0
	ps.Subscribe(context.Background(), "test.event", func(e Event) error {
		count++
		return nil
	})
	ps.Subscribe(context.Background(), "test.event", func(e Event) error {
		count++
		return nil
	})

	err := ps.Publish(context.Background(), "test.event", Event{Type: "test.event", Payload: []byte("{}")})
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}
