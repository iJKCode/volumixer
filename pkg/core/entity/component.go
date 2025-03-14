package entity

import (
	"ijkcode.tech/volumixer/pkg/util/typeset"
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

func GetComponent[C any](e *Entity) (any, bool) {
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
		e.components.Put(component)
	}
}

func SetComponent[C any](e *Entity, component C) {
	e.mut.Lock()
	defer e.mut.Unlock()
	typeset.Put[C](&e.components, component)
}

func (e *Entity) RemoveComponent(component any) {
	e.RemoveComponents(component)
}

func (e *Entity) RemoveComponents(components ...any) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, component := range components {
		e.components.Del(component)
	}
}

func (e *Entity) RemoveComponentType(component ComponentType) {
	e.RemoveComponentTypes(component)
}

func (e *Entity) RemoveComponentTypes(components ...ComponentType) {
	e.mut.Lock()
	defer e.mut.Unlock()
	for _, typ := range components {
		e.components.DelType(typ)
	}
}

func RemoveComponent[C any](e *Entity) {
	e.mut.Lock()
	defer e.mut.Unlock()
	typeset.Del[C](&e.components)
}

func (e *Entity) Components() []any {
	e.mut.RLock()
	defer e.mut.RUnlock()
	return e.components.Values()
}
