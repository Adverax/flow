package pubsub

import (
	"context"
	"github.com/adverax/flow/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type myEvent struct {
	name string
}

func TestPubSub(t *testing.T) {
	ctx := context.Background()

	ps, err := newTestPubSub[*myEvent]()
	require.NoError(t, err)
	defer ps.Close(ctx)

	var event *Event[*myEvent]
	ps.Subscribe(
		NewSubscription[*myEvent](
			HandlerFunc[*myEvent](func(ctx context.Context, e *Event[*myEvent]) {
				time.Sleep(100 * time.Millisecond)
				event = e
			}),
		),
	)

	ev := &myEvent{name: "event"}
	err = ps.Publish(ctx, ev).WaitWithContext(ctx)
	require.NoError(t, err)

	require.NotNil(t, event)
	assert.Equal(t, ev, event.Entity())
}

func newTestPubSub[T any]() (*PubSub[T], error) {
	return NewBuilder[T]().
		Subject("test").
		Executor(policy.NewDefaultExecutor()).
		Build()
}
