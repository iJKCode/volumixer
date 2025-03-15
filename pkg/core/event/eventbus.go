package event

import (
	"context"
	"ijkcode.tech/volumixer/pkg/util/typemap"
	"sync"
)

const DefaultQueueLength = 10

type Handler[E any] func(ctx context.Context, event E)

type Bus struct {
	mut      sync.Mutex
	events   chan any
	handlers typemap.TypeMap[handlerGroupAny]
	wildcard *handlerGroup[any]
}

func NewBus(options ...func(*Bus)) *Bus {
	bus := &Bus{
		mut:      sync.Mutex{},
		events:   nil,
		handlers: typemap.Make[handlerGroupAny](),
		wildcard: newHandlerGroup[any](),
	}
	if bus.events == nil {
		bus.events = make(chan any, DefaultQueueLength)
	}
	return bus
}

func WithChan(events chan any) func(*Bus) {
	return func(b *Bus) {
		b.events = events
	}
}

func WithQueueLength(length int) func(*Bus) {
	return func(b *Bus) {
		b.events = make(chan any, length)
	}
}

func Subscribe[E any](bus *Bus, handler Handler[E]) (unsubscribe func()) {
	if bus == nil || handler == nil {
		return func() {}
	}
	group := getHandlerGroupCreate[E](bus)
	return group.subscribe(handler)
}

func SubscribeAll(bus *Bus, handler Handler[any]) (unsubscribe func()) {
	if bus == nil || handler == nil {
		return func() {}
	}
	return bus.wildcard.subscribe(handler)
}

func Publish(bus *Bus, event any) {
	if bus == nil || event == nil {
		return
	}
	bus.events <- event
}

func (b *Bus) Chan() chan<- any {
	return b.events
}

func (b *Bus) Run(ctx context.Context) {
	for {
		select {
		case event := <-b.events:
			b.process(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bus) RunOnce(ctx context.Context) {
	select {
	case event := <-b.events:
		b.process(ctx, event)
	case <-ctx.Done():
		return
	}
}

func (b *Bus) process(ctx context.Context, event any) {
	if event == nil {
		return
	}
	group := getHandlerGroupAny(b, event)
	if group != nil {
		group.handleAny(ctx, event)
	}
	b.wildcard.handle(ctx, event)
}
