package pubsub

import (
	"context"
	"github.com/adverax/core"
	"github.com/adverax/flow/policy"
	"sync"
)

type Handler[T any] interface {
	Handle(ctx context.Context, event *Event[T])
}

type HandlerFunc[T any] func(ctx context.Context, event *Event[T])

func (fn HandlerFunc[T]) Handle(ctx context.Context, event *Event[T]) {
	fn(ctx, event)
}

type Subscriber[T any] interface {
	Handler[T]
	ID() string
	Close(ctx context.Context)
}

type PubSub[T any] struct {
	mx       sync.RWMutex
	subs     []Subscriber[T]
	executor policy.Executor
	subject  string
	closed   bool
}

func (that *PubSub[T]) Subject() string {
	return that.subject
}

func (that *PubSub[T]) Close(ctx context.Context) {
	that.mx.Lock()
	defer that.mx.Unlock()

	if !that.closed {
		that.closed = true
		for _, sub := range that.subs {
			sub.Close(ctx)
		}
	}
}

func (that *PubSub[T]) Subscribe(sub Subscriber[T]) {
	that.mx.Lock()
	defer that.mx.Unlock()

	that.subs = append(that.subs, sub)
}

func (that *PubSub[T]) Unsubscribe(ctx context.Context, subID string) {
	that.mx.Lock()
	defer that.mx.Unlock()

	for i, sub := range that.subs {
		if sub.ID() == subID {
			sub.Close(ctx)
			that.subs[i] = nil
			that.subs = append(that.subs[:i], that.subs[i+1:]...)
			break
		}
	}
}

func (that *PubSub[T]) Publish(ctx context.Context, entity T) core.Waiter {
	wg := core.NewWaitGroup()

	event := &Event[T]{
		ctx:     ctx,
		subject: that.subject,
		entity:  entity,
	}

	event.Enter()
	go that.publish(ctx, event)

	return wg
}

func (that *PubSub[T]) publish(ctx context.Context, event *Event[T]) {
	defer event.Leave()

	that.mx.RLock()
	defer that.mx.RUnlock()

	if that.closed {
		return
	}

	for _, sub := range that.subs {
		that.pub(ctx, sub, event)
	}
}

func (that *PubSub[T]) pub(ctx context.Context, sub Subscriber[T], event *Event[T]) {
	event.Enter()

	that.executor.Execute(
		ctx,
		policy.ActionFunc(
			func(ctx context.Context) error {
				defer event.Leave()

				sub.Handle(ctx, event)
				return nil
			},
		),
	)
}
