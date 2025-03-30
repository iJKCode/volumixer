package event

import (
	"context"
	"github.com/ijkcode/volumixer/pkg/util/typemap"
	"github.com/ijkcode/volumixer/pkg/util/valset"
	"maps"
	"sync"
)

type handlerGroupAny interface {
	handleAny(ctx context.Context, event any)
}

type handlerGroup[E any] struct {
	mut      sync.RWMutex
	handlers valset.Set[Handler[E]]
}

func newHandlerGroup[E any]() *handlerGroup[E] {
	return &handlerGroup[E]{
		mut:      sync.RWMutex{},
		handlers: make(valset.Set[Handler[E]]),
	}
}

func getHandlerGroupCreate[E any](bus *Bus) *handlerGroup[E] {
	bus.mut.Lock()
	defer bus.mut.Unlock()

	groupAny, ok := typemap.Get[E](&bus.handlers)
	if ok {
		return groupAny.(*handlerGroup[E])
	}

	group := newHandlerGroup[E]()
	typemap.Put[E, handlerGroupAny](&bus.handlers, group)

	return group
}

func getHandlerGroupAny(bus *Bus, event any) handlerGroupAny {
	bus.mut.Lock()
	defer bus.mut.Unlock()

	groupAny, ok := bus.handlers.Get(event)
	if !ok {
		return nil
	}

	return groupAny
}

func (s *handlerGroup[E]) subscribe(handler Handler[E]) (unsubscribe func()) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.handlers.Put(handler)
	return func() {
		s.mut.Lock()
		defer s.mut.Unlock()
		s.handlers.Del(handler)
	}
}

func (s *handlerGroup[E]) handleAny(ctx context.Context, event any) {
	s.handle(ctx, event.(E))
}

func (s *handlerGroup[E]) handle(ctx context.Context, event E) {
	s.mut.RLock()
	handlers := maps.Clone(s.handlers)
	s.mut.RUnlock()

	for handler := range handlers {
		handler.Handle(ctx, event)
	}
}
