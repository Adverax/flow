package policy

import (
	"context"
	"errors"
	"time"
)

type RetryPolicyErrorChecker interface {
	IsRetryableError(err error) bool
}

type RetryPolicyMetrics interface {
	IncSuccess()
	IncFailure()
	IncAttempts()
}

type RetryPolicyOptions struct {
	InitialInterval       time.Duration
	BackoffCoefficient    float64
	MaximumInterval       time.Duration
	MaximumAttempts       int
	RetryableErrorChecker RetryPolicyErrorChecker
	Metrics               RetryPolicyMetrics
}

type retryState struct {
	interval time.Duration
	attempts int
}

type PolicyWithRetry struct {
	Policy
	options RetryPolicyOptions
}

func NewPolicyWithRetry(policy Policy, options RetryPolicyOptions) *PolicyWithRetry {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithRetry{
		options: options,
		Policy:  policy,
	}
}

func (that *PolicyWithRetry) Execute(ctx context.Context, action Action) error {
	err := that.Policy.Execute(ctx, action)
	if err == nil {
		that.success()
		return nil
	}

	if that.options.MaximumAttempts < 0 {
		that.failure()
		return err
	}

	err = that.retry(ctx, action, err)
	if err == nil {
		that.success()
		return nil
	}

	that.failure()
	return err
}

func (that *PolicyWithRetry) retry(ctx context.Context, action Action, err error) error {
	state := retryState{
		interval: that.options.InitialInterval,
	}
	for that.canAttempt(err, &state) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(state.interval):
		}

		state.attempts++
		that.attempt()

		err = that.Policy.Execute(ctx, action)
		if err == nil {
			return nil
		}
	}
	return err
}

func (that *PolicyWithRetry) canAttempt(err error, state *retryState) bool {
	if !that.IsRetryableError(err) {
		return false
	}

	if state.attempts >= that.options.MaximumAttempts && that.options.MaximumAttempts != 0 {
		return false
	}

	interval := time.Duration(that.options.BackoffCoefficient * float64(state.interval))
	if interval > that.options.MaximumInterval {
		state.interval = that.options.MaximumInterval
	} else {
		state.interval = interval
	}

	return true
}

func (that *PolicyWithRetry) success() {
	if that.options.Metrics != nil {
		that.options.Metrics.IncSuccess()
	}
}

func (that *PolicyWithRetry) failure() {
	if that.options.Metrics != nil {
		that.options.Metrics.IncFailure()
	}
}

func (that *PolicyWithRetry) attempt() {
	if that.options.Metrics != nil {
		that.options.Metrics.IncAttempts()
	}
}

func (that *PolicyWithRetry) IsRetryableError(err error) bool {
	return that.options.RetryableErrorChecker == nil || that.options.RetryableErrorChecker.IsRetryableError(err)
}

type retryErrorChecker struct {
	nonRetryableErrors []error
}

func NewRetryErrorChecker(nonRetryableErrors []error) RetryPolicyErrorChecker {
	return &retryErrorChecker{
		nonRetryableErrors: nonRetryableErrors,
	}
}

func (that *retryErrorChecker) IsRetryableError(err error) bool {
	for _, nonRetryableError := range that.nonRetryableErrors {
		if errors.Is(err, nonRetryableError) {
			return false
		}
	}
	return true
}
