package policy

import (
	"context"
	"github.com/adverax/log"
)

type ErrorHandlerWithLogging struct {
	logger log.Logger
}

func (that *ErrorHandlerWithLogging) HandleError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	that.logger.Error(ctx, err.Error())
}

func NewErrorHandlerWithLogging(logger log.Logger) *ErrorHandlerWithLogging {
	return &ErrorHandlerWithLogging{
		logger: logger,
	}
}
