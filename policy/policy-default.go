package policy

import "context"

type DefaultPolicy struct{}

func NewDefaultPolicy() *DefaultPolicy {
	return &DefaultPolicy{}
}

func (c *DefaultPolicy) Execute(ctx context.Context, action Action) error {
	return action.Execute(ctx)
}

var dummyPolicy = NewDefaultPolicy()
