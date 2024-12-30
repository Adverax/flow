package pipelines

import (
	"context"
	"github.com/adverax/flow/policy"
	"sync"
)

// WorkerPool is pool of worker for handle queue.
type WorkerPool struct {
	count    int
	wg       sync.WaitGroup
	queue    <-chan policy.Action
	executor policy.Executor
}

func NewWorkerPool(
	queue <-chan policy.Action,
	executor policy.Executor,
	count int,
) *WorkerPool {
	return &WorkerPool{
		queue:    queue,
		count:    count,
		executor: executor,
	}
}

func (that *WorkerPool) Start(ctx context.Context) {
	that.wg.Add(that.count)
	for i := 0; i < that.count; i++ {
		go that.serve(ctx)
	}
}

func (that *WorkerPool) serve(ctx context.Context) {
	defer that.wg.Done()

	for action := range that.queue {
		that.executor.Execute(ctx, action)
	}
}
