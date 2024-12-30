package policy

import (
	"context"
	"fmt"
	"github.com/adverax/log"
)

type PolicyWithRecovery struct {
	Policy
	logger log.Logger
}

func (that *PolicyWithRecovery) Execute(ctx context.Context, action Action) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
			if that.logger != nil {
				that.logger.Fatal(ctx, err.Error())
			}
		}
	}()

	return that.Policy.Execute(ctx, action)
}

func NewPolicyWithRecovery(policy Policy, logger log.Logger) *PolicyWithRecovery {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithRecovery{
		Policy: policy,
		logger: logger,
	}
}
