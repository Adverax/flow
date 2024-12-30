package trace

import (
	"context"
	"github.com/adverax/log"
)

const (
	FieldKeyTraceID = "trace_id"
)

type LoggerHook struct {
}

func NewLoggerHook() *LoggerHook {
	return &LoggerHook{}
}

func (that *LoggerHook) Fire(ctx context.Context, entry *log.Entry) error {
	traceId := GetId(ctx)
	entry.Data[FieldKeyTraceID] = traceId
	return nil
}
