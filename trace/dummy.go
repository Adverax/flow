package trace

import "context"

type DummyTracer struct {
}

func (that *DummyTracer) NewTrace(ctx context.Context, traceId string, info string) context.Context {
	return ctx
}

func (that *DummyTracer) NewTraceEx(ctx context.Context, traceId string, info string) (context.Context, string) {
	return ctx, ""
}

func NewDummy() *DummyTracer {
	return &DummyTracer{}
}
