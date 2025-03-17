package typeset

import "reflect"

func Has[V any](ts *TypeSet) bool {
	typ := reflect.TypeFor[V]()
	return ts.HasType(typ)
}

func Get[V any](ts *TypeSet) (V, bool) {
	typ := reflect.TypeFor[V]()
	val, ok := ts.GetType(typ)
	if !ok {
		return zero[V](), ok
	}
	return val.(V), ok
}

func Put[V any](ts *TypeSet, val V) {
	typ := reflect.TypeFor[V]()
	ts.values[typ] = val
}

func Swap[V any](ts *TypeSet, val V) (V, bool) {
	typ := reflect.TypeFor[V]()
	old, ok := ts.SwapType(typ, val)
	if !ok {
		return zero[V](), ok
	}
	return old.(V), ok
}

func Del[V any](ts *TypeSet) {
	typ := reflect.TypeFor[V]()
	ts.DelType(typ)
}

func Pop[V any](ts *TypeSet) (V, bool) {
	typ := reflect.TypeFor[V]()
	val, ok := ts.PopType(typ)
	if !ok {
		return zero[V](), ok
	}
	return val.(V), ok
}

func zero[T any]() T {
	var value T
	return value
}
