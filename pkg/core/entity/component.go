package entity

import (
	"github.com/ijkcode/volumixer/pkg/util/typeset"
	"reflect"
)

func (e *Entity) HasComponent(component any) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.Has(component)
}

func (e *Entity) HasComponentType(typ ComponentType) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.HasType(typ)
}

func HasComponent[C any](e *Entity) bool {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return typeset.Has[C](&e.components)
}

func (e *Entity) GetComponent(component any) (any, bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.Get(component)
}

func (e *Entity) GetComponentType(component ComponentType) (any, bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.GetType(component)
}

func GetComponent[C any](e *Entity) (C, bool) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return typeset.Get[C](&e.components)
}

func (e *Entity) SetComponent(component any) {
	e.SetComponents(component)
}

func (e *Entity) SetComponents(components ...any) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, component := range components {
		if component == nil {
			continue
		}
		old, ok := e.components.Swap(component)
		if ok && reflect.DeepEqual(old, component) {
			continue
		}
		e.publishEvent(ComponentUpdatedEvent{
			Entity:    e,
			Component: component,
		})
	}
}

func SetComponent[C any](e *Entity, component C) {
	e.mut.Lock()
	defer e.mut.Unlock()
	old, ok := typeset.Swap[C](&e.components, component)
	if ok && reflect.DeepEqual(old, component) {
		return
	}
	e.publishEvent(ComponentUpdatedEvent{
		Entity:    e,
		Component: component,
	})
}

func (e *Entity) RemoveComponent(component any) {
	e.RemoveComponents(component)
}

func (e *Entity) RemoveComponents(components ...any) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, component := range components {
		if component == nil {
			continue
		}
		value, ok := e.components.Pop(component)
		if ok {
			e.publishEvent(ComponentRemovedEvent{
				Entity:    e,
				Component: value,
			})
		}
	}
}

func (e *Entity) RemoveComponentType(component ComponentType) {
	e.RemoveComponentTypes(component)
}

func (e *Entity) RemoveComponentTypes(components ...ComponentType) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, typ := range components {
		if typ == nil {
			continue
		}
		value, ok := e.components.PopType(typ)
		if ok {
			e.publishEvent(ComponentRemovedEvent{
				Entity:    e,
				Component: value,
			})
		}
	}
}

func RemoveComponent[C any](e *Entity) {
	e.mut.Lock()
	defer e.mut.Unlock()
	value, ok := typeset.Pop[C](&e.components)
	if ok {
		e.publishEvent(ComponentRemovedEvent{
			Entity:    e,
			Component: value,
		})
	}
}

func (e *Entity) Components() []any {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.Values()
}
