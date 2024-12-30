package scheduler

import (
	"fmt"
	"github.com/adverax/core"
	"github.com/adverax/flow/policy"
	"time"
)

type Builder struct {
	*core.Builder
	scheduler *Engine
}

func NewBuilder() *Builder {
	return &Builder{
		Builder: core.NewBuilder("scheduler"),
		scheduler: &Engine{
			interval:    time.Minute,
			tickOnStart: false,
			control:     core.NewDummyWaitGroup(),
			executor:    policy.NewDefaultExecutor(),
		},
	}
}

func (that *Builder) WithInterval(interval time.Duration) *Builder {
	that.scheduler.interval = interval
	return that
}

func (that *Builder) WithTickOnStart(tickOnStart bool) *Builder {
	that.scheduler.tickOnStart = tickOnStart
	return that
}

func (that *Builder) WithExecutor(executor policy.Executor) *Builder {
	that.scheduler.executor = executor
	return that
}

func (that *Builder) WithTask(task Task) *Builder {
	that.scheduler.task = task
	return that
}

func (that *Builder) WithControl(control Control) *Builder {
	that.scheduler.control = control
	return that
}

func (that *Builder) Build() (*Engine, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.scheduler, nil
}

func (that *Builder) checkRequiredFields() error {
	that.RequiredField(that.scheduler.task, ErrFieldTaskRequired)
	return that.ResError()
}

var (
	ErrFieldTaskRequired = fmt.Errorf("Field task is required")
)
