package policy

import (
	"context"
)

type Tracer interface {
	NewTrace(ctx context.Context, info string) context.Context
}

type PolicyWithTraceId struct {
	Policy
	tracer Tracer
	info   string
}

func (that *PolicyWithTraceId) Execute(ctx context.Context, action Action) error {
	ctx = that.tracer.NewTrace(ctx, that.info)
	return that.Policy.Execute(ctx, action)
}

func NewPolicyWithTraceId(
	policy Policy,
	tracer Tracer,
	info string,
) *PolicyWithTraceId {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithTraceId{
		Policy: policy,
		tracer: tracer,
		info:   info,
	}
}
