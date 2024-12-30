package policy

import (
	"context"
	"github.com/adverax/log"
)

type Loggable interface {
	Named
	Log(ctx context.Context, logger log.Logger)
}

type PolicyWithLoggging struct {
	Policy
	logger log.Logger
}

func NewPolicyWithLogging(policy Policy, logger log.Logger) *PolicyWithLoggging {
	return &PolicyWithLoggging{
		Policy: policy,
		logger: logger,
	}
}

func (that *PolicyWithLoggging) Execute(ctx context.Context, action Action) error {
	if a, ok := action.(Loggable); ok {
		a.Log(ctx, that.logger)

		err := that.Policy.Execute(ctx, action)
		if err != nil {
			that.logger.WithError(err).Errorf(ctx, "error executing action %s", a.Name())
			return err
		}

		return nil
	}

	err := that.Policy.Execute(ctx, action)
	if err != nil {
		return err
	}

	return nil
}
