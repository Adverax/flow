package scheduler

import (
	"context"
	"github.com/adverax/flow/policy"
	"time"
)

type Task interface {
	policy.Action
	IsActive() bool
}

type Control interface {
	Enter()
	Leave()
}

type Engine struct {
	interval    time.Duration
	tickOnStart bool
	task        Task
	control     Control
	executor    policy.Executor
}

func (that *Engine) Start(ctx context.Context) {
	that.control.Enter()

	go func() {
		defer that.control.Leave()

		that.run(ctx)
	}()
}

func (that *Engine) run(ctx context.Context) {
	ticker := time.NewTicker(that.interval)
	defer ticker.Stop()

	if that.tickOnStart {
		that.step(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			that.step(ctx)
		}
	}
}

func (that *Engine) step(ctx context.Context) {
	if !that.task.IsActive() {
		that.executor.Execute(ctx, that.task)
	}
}
