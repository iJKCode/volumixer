package typeset

import "reflect"

func Has[V any](ts *TypeSet) bool {
	typ := reflect.TypeFor[V]()
	return ts.HasType(typ)
}

func Get[V any](ts *TypeSet) (any, bool) {
	typ := reflect.TypeFor[V]()
	return ts.GetType(typ)
}

func Put[V any](ts *TypeSet, val V) {
	typ := reflect.TypeFor[V]()
	ts.values[typ] = val
}

func Del[V any](ts *TypeSet) {
	typ := reflect.TypeFor[V]()
	ts.DelType(typ)
}
