package policy

import (
	"context"
)

type PolicyWithContextValue struct {
	Policy
	key   interface{}
	value interface{}
}

func (that *PolicyWithContextValue) Execute(ctx context.Context, action Action) error {
	ctx = context.WithValue(ctx, that.key, that.value)
	return that.Policy.Execute(ctx, action)
}

func NewPolicyWithContextValue(
	policy Policy,
	key interface{},
	value interface{},
) *PolicyWithContextValue {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithContextValue{
		Policy: policy,
		key:    key,
		value:  value,
	}
}
