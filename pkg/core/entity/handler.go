package entity

import (
	"ijkcode.tech/volumixer/pkg/util/typemap"
)

type Handler[T any] func(ent *Entity, cmd any) error

func (e *Entity) HasHandler(command any) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.handlers.Has(command)
}

func (e *Entity) HasHandlerType(command CommandType) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.handlers.HasType(command)
}

func HasHandler[C any](e *Entity) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return typemap.Has[C](&e.handlers)
}

func (e *Entity) GetHandler(command any) (Handler[any], bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.handlers.Get(command)
}

func (e *Entity) GetHandlerType(command CommandType) (Handler[any], bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.handlers.GetType(command)
}

func GetHandler[C any](e *Entity) (Handler[any], bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return typemap.Get[C](&e.handlers)
}

func (e *Entity) SetHandler(command any, handler Handler[any]) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.handlers.Put(command, handler)
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func (e *Entity) SetHandlerType(command CommandType, handler Handler[any]) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.handlers.PutType(command, handler)
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func SetHandler[C any](e *Entity, handler Handler[C]) {
	e.mut.Lock()
	defer e.mut.Unlock()
	wrapped := wrapHandler(handler)
	typemap.Put[C](&e.handlers, wrapped)
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func (e *Entity) RemoveHandler(command any) {
	e.RemoveHandlers(command)
}

func (e *Entity) RemoveHandlers(commands ...any) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, command := range commands {
		e.handlers.Del(command)
	}
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func (e *Entity) RemoveHandlerType(command CommandType) {
	e.RemoveHandlers(command)
}

func (e *Entity) RemoveHandlerTypes(commands ...CommandType) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, command := range commands {
		e.handlers.DelType(command)
	}
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func RemoveHandler[C any](e *Entity) {
	e.mut.Lock()
	defer e.mut.Unlock()
	typemap.Del[C](&e.handlers)
	e.publishEvent(HandlersUpdatedEvent{
		Entity: e,
	})
}

func (e *Entity) Handlers() []Handler[any] {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.handlers.Values()
}

func wrapHandler[T any](handler Handler[T]) Handler[any] {
	return func(ent *Entity, command any) error {
		return handler(ent, command)
	}
}
