package promise

import (
	"context"
	"github.com/adverax/flow/policy"
	"sync"
)

type Action[T any] func(ctx context.Context) (T, error)

type Consumer[T any] func(ctx context.Context, val T, err error)

type Promise[T any] struct {
	value    T
	err      error
	ch       chan struct{}
	once     sync.Once
	consumer Consumer[T]
}

func (that *Promise[T]) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-that.ch:
		return that.err
	}
}

func (that *Promise[T]) Await(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var val T
		return val, ctx.Err()
	case <-that.ch:
		return that.value, that.err
	}
}

func (that *Promise[T]) resolve(ctx context.Context, value T) {
	that.once.Do(
		func() {
			that.value = value
			close(that.ch)
			if that.consumer != nil {
				that.consumer(ctx, value, nil)
			}
		},
	)
}

func (that *Promise[T]) reject(ctx context.Context, err error) {
	that.once.Do(
		func() {
			that.err = err
			close(that.ch)
			if that.consumer != nil {
				var val T
				that.consumer(ctx, val, err)
			}
		},
	)
}

var DefaultExecutor = policy.NewDefaultExecutor()

type Waiter interface {
	Wait(ctx context.Context) error
}

func WaitAll(ctx context.Context, waiter ...Waiter) error {
	for _, w := range waiter {
		err := w.Wait(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
