package typemap

import "reflect"

func Has[K any, V any](tm *TypeMap[V]) bool {
	typ := reflect.TypeFor[K]()
	return tm.HasType(typ)
}

func Get[K any, V any](tm *TypeMap[V]) (V, bool) {
	typ := reflect.TypeFor[K]()
	return tm.GetType(typ)
}

func Put[K any, V any](tm *TypeMap[V], val V) {
	typ := reflect.TypeFor[K]()
	tm.PutType(typ, val)
}

func Del[K any, V any](tm *TypeMap[V]) {
	typ := reflect.TypeFor[K]()
	tm.DelType(typ)
}
