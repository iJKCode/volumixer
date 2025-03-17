package typeset

import (
	"maps"
	"reflect"
)

type ValType = reflect.Type

type TypeSet struct {
	values map[ValType]any
}

func Make() TypeSet {
	return TypeSet{
		values: make(map[ValType]any),
	}
}

func (ts *TypeSet) Has(val any) bool {
	assertNotReflectType(val, "Has")
	typ := reflect.TypeOf(val)
	return ts.HasType(typ)
}

func (ts *TypeSet) HasType(typ ValType) bool {
	_, ok := ts.values[typ]
	return ok
}

func (ts *TypeSet) Get(val any) (any, bool) {
	assertNotReflectType(val, "Get")
	typ := reflect.TypeOf(val)
	return ts.GetType(typ)
}

func (ts *TypeSet) GetType(typ ValType) (any, bool) {
	val, ok := ts.values[typ]
	return val, ok
}

func (ts *TypeSet) Put(val any) {
	assertNotReflectType(val, "Put")
	typ := reflect.TypeOf(val)
	ts.values[typ] = val
}

func (ts *TypeSet) Swap(val any) (any, bool) {
	assertNotReflectType(val, "Pop")
	typ := reflect.TypeOf(val)
	return ts.SwapType(typ, val)
}

func (ts *TypeSet) SwapType(key ValType, val any) (any, bool) {
	old, ok := ts.values[key]
	ts.values[key] = val
	return old, ok
}

func (ts *TypeSet) Del(val any) {
	assertNotReflectType(val, "Del")
	typ := reflect.TypeOf(val)
	ts.DelType(typ)
}

func (ts *TypeSet) DelType(key ValType) {
	delete(ts.values, key)
}

func (ts *TypeSet) Pop(val any) (any, bool) {
	assertNotReflectType(val, "Pop")
	typ := reflect.TypeOf(val)
	return ts.PopType(typ)
}

func (ts *TypeSet) PopType(key ValType) (any, bool) {
	val, ok := ts.values[key]
	delete(ts.values, key)
	return val, ok
}

func (ts *TypeSet) Len() int {
	return len(ts.values)
}

func (ts *TypeSet) Values() []any {
	values := make([]any, 0, len(ts.values))
	for _, val := range ts.values {
		values = append(values, val)
	}
	return values
}

func (ts *TypeSet) Items() map[ValType]any {
	return maps.Clone(ts.values)
}

func (ts *TypeSet) UnsafeRef() map[ValType]any {
	return ts.values
}

func assertNotReflectType(val any, fn string) {
	_, ok := val.(reflect.Type)
	if ok {
		panic("use the %sType method when performing operations using a reflect.Type")
	}
}
