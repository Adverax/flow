package pubsub

import (
	"context"
	"sync"
)

type Event[T any] struct {
	ctx     context.Context
	wg      sync.WaitGroup
	subject string
	entity  T
}

func (that *Event[T]) Enter() {
	that.wg.Add(1)
}

func (that *Event[T]) Leave() {
	that.wg.Done()
}

func (that *Event[T]) Context() context.Context {
	return that.ctx
}

func (that *Event[T]) Subject() string {
	return that.subject
}

func (that *Event[T]) Entity() T {
	return that.entity
}
