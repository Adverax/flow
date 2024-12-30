package promise

import (
	"context"
	"fmt"
	"github.com/adverax/core"
	"github.com/adverax/flow/policy"
)

type Builder[T any] struct {
	*core.Builder
	promise  *Promise[T]
	action   Action[T]
	executor policy.Executor
	ctx      context.Context
}

func New[T any]() *Builder[T] {
	return &Builder[T]{
		Builder: core.NewBuilder("promise"),
		promise: &Promise[T]{
			ch: make(chan struct{}),
		},
		executor: DefaultExecutor,
	}
}

func (that *Builder[T]) WithAction(action Action[T]) *Builder[T] {
	that.action = action
	return that
}

func (that *Builder[T]) WithConsumer(consumer Consumer[T]) *Builder[T] {
	that.promise.consumer = consumer
	return that
}

func (that *Builder[T]) WithExecutor(executor policy.Executor) *Builder[T] {
	that.executor = executor
	return that
}

func (that *Builder[T]) WithContext(ctx context.Context) *Builder[T] {
	that.ctx = ctx
	return that
}

func (that *Builder[T]) Build() (*Promise[T], error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	that.executor.Execute(
		that.ctx,
		policy.ActionFunc(func(ctx context.Context) error {
			val, err := that.action(ctx)
			if err != nil {
				that.promise.reject(ctx, err)
			} else {
				that.promise.resolve(ctx, val)
			}
			return nil
		}),
	)

	return that.promise, nil
}

func (that *Builder[T]) checkRequiredFields() error {
	if that.action == nil {
		return ErrFieldActionRequired
	}
	return nil
}

var (
	ErrFieldActionRequired = fmt.Errorf("field action is required")
)
