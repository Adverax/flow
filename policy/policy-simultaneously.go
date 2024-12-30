package policy

import (
	"context"
	"github.com/adverax/core"
)

type PolicyWithAsyncExecution struct {
	Executor
	control core.Control
}

func newPolicyWithAsyncExecution(
	executor Executor,
	control core.Control,
) *PolicyWithAsyncExecution {
	if control == nil {
		control = core.NewDummyWaitGroup()
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}
	return &PolicyWithAsyncExecution{
		control:  control,
		Executor: executor,
	}
}

func (that *PolicyWithAsyncExecution) Execute(ctx context.Context, action Action) error {
	that.control.Enter()
	go func() {
		defer that.control.Leave()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

type PolicyWithPoolExecution struct {
	Executor
	pool    chan struct{}
	control core.Control
}

func newPolicyWithPoolExecution(
	executor Executor,
	control core.Control,
	size int,
) *PolicyWithPoolExecution {
	if control == nil {
		control = core.NewDummyWaitGroup()
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}

	pool := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		pool <- struct{}{}
	}

	return &PolicyWithPoolExecution{
		pool:     pool,
		control:  control,
		Executor: executor,
	}
}

func (that *PolicyWithPoolExecution) Execute(ctx context.Context, action Action) error {
	that.control.Enter()

	select {
	case <-that.pool:
	case <-ctx.Done():
		return ctx.Err()
	}

	go func() {
		defer func() {
			that.pool <- struct{}{}
			that.control.Leave()
		}()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

func NewSimultaneouslyPolicy(
	executor Executor,
	control core.Control,
	size int,
) Policy {
	if size == 0 {
		return newPolicyWithAsyncExecution(executor, control)
	}

	return newPolicyWithPoolExecution(executor, control, size)
}
