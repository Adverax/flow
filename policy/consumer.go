package policy

import (
	"context"
	"github.com/adverax/core"
)

type Consumer[T any] interface {
	Consume(ctx context.Context, entity T)
}

type ConsumerFunc[T any] func(context.Context, T)

func (fn ConsumerFunc[T]) Consume(ctx context.Context, entity T) {
	fn(ctx, entity)
}

type BaseConsumer[T any] struct {
	executor Executor
	handler  core.Handler[T]
}

func NewConsumer[T any](executor Executor, handler core.Handler[T]) *BaseConsumer[T] {
	return &BaseConsumer[T]{
		executor: executor,
		handler:  handler,
	}
}

func (that *BaseConsumer[T]) Consume(ctx context.Context, entity T) {
	that.executor.Execute(
		ctx,
		ActionFunc(func(ctx context.Context) error {
			return that.handler.Handle(ctx, entity)
		}),
	)
}
