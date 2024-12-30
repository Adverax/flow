package policy

import (
	"context"
)

type Action interface {
	Execute(ctx context.Context) error
}

type ActionFunc func(ctx context.Context) error

func (fn ActionFunc) Execute(ctx context.Context) error {
	return fn(ctx)
}

type Policy interface {
	Execute(ctx context.Context, action Action) error
}
