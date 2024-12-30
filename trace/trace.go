package trace

import (
	"context"
	"github.com/adverax/core"
	"github.com/adverax/log"
)

type typeTraceId int

var (
	traceIdKey typeTraceId = 1
)

type Engine struct {
	logger log.Logger
}

func New(logger log.Logger) *Engine {
	return &Engine{
		logger: logger,
	}
}

func (that *Engine) EnsureTrace(ctx context.Context, info string) context.Context {
	traceId := GetId(ctx)
	if traceId == "" {
		return that.NewTrace(ctx, info)
	}

	return ctx
}

func (that *Engine) NewTrace(ctx context.Context, info string) context.Context {
	ctx2, _ := that.NewTraceEx(ctx, "", info)
	return ctx2
}

func (that *Engine) NewTraceWithId(ctx context.Context, traceId string, info string) context.Context {
	ctx2, _ := that.NewTraceEx(ctx, traceId, info)
	return ctx2
}

func (that *Engine) NewTraceEx(ctx context.Context, traceId string, info string) (context.Context, string) {
	if traceId == "" {
		traceId = core.NewGUID()
	}

	if info != "" {
		that.logger.WithField("trace_id", traceId).Trace(ctx, "New trace: "+info)
	}

	return context.WithValue(ctx, traceIdKey, traceId), traceId
}

func GetId(ctx context.Context) string {
	traceId, _ := ctx.Value(traceIdKey).(string)
	return traceId
}
