package typemap

import (
	"maps"
	"reflect"
)

type KeyType = reflect.Type

type TypeMap[V any] struct {
	values map[KeyType]V
}

func Make[V any]() TypeMap[V] {
	return TypeMap[V]{
		values: make(map[KeyType]V),
	}
}

func (tm *TypeMap[V]) Has(key any) bool {
	assertNotReflectType(key, "Has")
	typ := reflect.TypeOf(key)
	return tm.HasType(typ)
}

func (tm *TypeMap[V]) HasType(key KeyType) bool {
	_, ok := tm.values[key]
	return ok
}

func (tm *TypeMap[V]) Get(key any) (V, bool) {
	assertNotReflectType(key, "Get")
	typ := reflect.TypeOf(key)
	return tm.GetType(typ)
}

func (tm *TypeMap[V]) GetType(key KeyType) (V, bool) {
	val, ok := tm.values[key]
	return val, ok
}

func (tm *TypeMap[V]) Put(key any, val V) {
	assertNotReflectType(key, "Put")
	typ := reflect.TypeOf(key)
	tm.PutType(typ, val)
}

func (tm *TypeMap[V]) PutType(key KeyType, val V) {
	tm.values[key] = val
}

func (tm *TypeMap[V]) Del(key any) {
	assertNotReflectType(key, "Del")
	typ := reflect.TypeOf(key)
	tm.DelType(typ)
}

func (tm *TypeMap[V]) DelType(key KeyType) {
	delete(tm.values, key)
}

func (tm *TypeMap[V]) Len() int {
	return len(tm.values)
}

func (tm *TypeMap[V]) Values() []V {
	values := make([]V, 0, len(tm.values))
	for _, val := range tm.values {
		values = append(values, val)
	}
	return values
}

func (tm *TypeMap[V]) Items() map[KeyType]V {
	return maps.Clone(tm.values)
}

func (tm *TypeMap[V]) UnsafeRef() map[KeyType]V {
	return tm.values
}

func assertNotReflectType(val any, fn string) {
	_, ok := val.(reflect.Type)
	if ok {
		panic("use the %sType method when performing operations using a reflect.Type")
	}
}
