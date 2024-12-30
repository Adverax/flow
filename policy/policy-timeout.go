package policy

import (
	"context"
	"time"
)

type PolicyWithTimeout struct {
	Policy
	timeout time.Duration
}

func (that *PolicyWithTimeout) Execute(ctx context.Context, action Action) error {
	ctx2, cancel := context.WithTimeout(ctx, that.timeout)
	defer cancel()

	return that.Policy.Execute(ctx2, action)
}

func NewPolicyWithTimeout(policy Policy, timeout time.Duration) *PolicyWithTimeout {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithTimeout{
		Policy:  policy,
		timeout: timeout,
	}
}
