package pubsub

import (
	"fmt"
	"github.com/adverax/core"
	"github.com/adverax/flow/policy"
)

type Builder[T any] struct {
	*core.Builder
	pubsub *PubSub[T]
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		Builder: core.NewBuilder("PubSub"),
		pubsub:  &PubSub[T]{},
	}
}

func (that *Builder[T]) Subject(subject string) *Builder[T] {
	that.pubsub.subject = subject
	return that
}

func (that *Builder[T]) Executor(executor policy.Executor) *Builder[T] {
	that.pubsub.executor = executor
	return that
}

func (that *Builder[T]) Build() (*PubSub[T], error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	if err := that.updateDefaultFields(); err != nil {
		return nil, err
	}

	return that.pubsub, nil
}

func (that *Builder[T]) checkRequiredFields() error {
	that.RequiredField(that.pubsub.subject, ErrFieldSubjectIsRequired)

	return that.ResError()
}

func (that *Builder[T]) updateDefaultFields() error {
	if that.pubsub.executor == nil {
		that.pubsub.executor = policy.NewDefaultExecutor()
	}

	return that.ResError()
}

var (
	ErrFieldSubjectIsRequired = fmt.Errorf("Field 'subject' is required")
)
