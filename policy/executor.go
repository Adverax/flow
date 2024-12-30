package policy

import (
	"context"
	"github.com/adverax/core"
)

type Executor interface {
	Execute(ctx context.Context, action Action)
}

type BaseExecutor struct {
	policy Policy
	errors core.ErrorHandler
}

func NewBaseExecutor(
	policy Policy,
	errors core.ErrorHandler,
) *BaseExecutor {
	return &BaseExecutor{
		policy: policy,
		errors: errors,
	}
}

func (that *BaseExecutor) Execute(
	ctx context.Context,
	action Action,
) {
	err := that.policy.Execute(ctx, action)
	if err != nil {
		that.errors.HandleError(ctx, err)
	}
}

func NewDefaultExecutor() Executor {
	return NewBaseExecutor(
		NewDefaultPolicy(),
		core.NewDefaultErrorHandler(),
	)
}
