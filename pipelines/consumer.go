package pipelines

import (
	"context"
	"github.com/adverax/core"
)

func Consume[T any](
	ctx context.Context,
	income <-chan T,
	aggregator core.Aggregator[T],
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-income:
			if !ok {
				return nil
			}

			err := aggregator.Aggregate(ctx, val)
			if err != nil {
				return err
			}
		}
	}
}
