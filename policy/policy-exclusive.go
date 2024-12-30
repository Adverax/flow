package policy

import (
	"context"
	"sync"
)

type PolicyWithExclusiveExecution struct {
	Policy Policy
	sync.Mutex
}

func NewPolicyWithExclusiveExecution(policy Policy) *PolicyWithExclusiveExecution {
	return &PolicyWithExclusiveExecution{Policy: policy}
}

func (that *PolicyWithExclusiveExecution) Execute(ctx context.Context, action Action) error {
	that.Lock()
	defer that.Unlock()

	return that.Policy.Execute(ctx, action)
}
