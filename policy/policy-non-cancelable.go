package policy

import (
	"context"
	"github.com/adverax/core"
)

type PolicyNonCancelable struct {
	Policy
}

func (that *PolicyNonCancelable) Execute(ctx context.Context, action Action) error {
	ctx = core.NewNonCancelableContext(ctx)
	return that.Policy.Execute(ctx, action)
}

func NewPolicyNonCancelable(
	policy Policy,
) *PolicyNonCancelable {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyNonCancelable{
		Policy: policy,
	}
}
