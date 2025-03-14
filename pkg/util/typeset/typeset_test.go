package typeset

import (
	"reflect"
	"testing"
)
import "github.com/stretchr/testify/assert"

type Value1 struct {
	V int
}
type Value2 struct {
	V string
}
type Value3 struct {
	Value bool
}
type ValueStr string

func TestMemberFunctions(t *testing.T) {
	s := Make()

	// add some values
	s.Put(Value1{42})
	s.Put(Value2{"hello"})
	s.Put(ValueStr("world"))

	{
		// check has functionality
		assert.True(t, s.Has(Value1{}), "should contain Value1")
		assert.True(t, s.Has(Value2{}), "should contain Value2")
		assert.False(t, s.Has(Value3{}), "should not contain Value3")
		assert.True(t, s.Has(ValueStr("")), "should contain ValueStr")
	}
	{
		// check get functionality
		v, ok := s.Get(Value1{})
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, Value1{42}, v)
		v, ok = s.Get(Value2{})
		assert.True(t, ok, "Value2 should exist")
		assert.Equal(t, Value2{"hello"}, v)
		v, ok = s.Get(Value3{})
		assert.False(t, ok, "Value3 should not exist")
		v, ok = s.Get(ValueStr(""))
		assert.True(t, ok, "ValueStr should exist")
		assert.Equal(t, ValueStr("world"), v)
	}
	{
		// check replace functionality
		s.Put(Value1{99})
		v, ok := s.Get(Value1{})
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, Value1{99}, v)
		assert.Equal(t, 3, s.Len(), "length should be 3")
	}
	{
		// check delete functionality
		s.Del(Value1{})
		assert.False(t, s.Has(Value1{}), "Value1 should be removed")
		assert.True(t, s.Has(Value2{}), "Value2 should still exist")
		assert.Equal(t, 2, s.Len(), "length should be 2")
	}
	{
		// re-add deleted value
		s.Put(Value1{67})
		v, ok := s.Get(Value1{})
		assert.True(t, ok, "Value1 should exist again")
		assert.Equal(t, Value1{67}, v)
		assert.Equal(t, 3, s.Len(), "length should be 3")
	}
	{
		// check all items functionality
		a := s.Items()
		assert.Equal(t, a, map[ValType]any{
			reflect.TypeOf(Value1{}):     Value1{67},
			reflect.TypeOf(Value2{}):     Value2{"hello"},
			reflect.TypeOf(ValueStr("")): ValueStr("world"),
		})
	}
	{
		a := s.Values()
		assert.Equal(t, 3, len(a), "length should be 3")
		assert.Contains(t, a, Value1{67})
		assert.NotContains(t, a, Value1{99})
		assert.Contains(t, a, Value2{"hello"})
		assert.Contains(t, a, ValueStr("world"))
	}
	{
		// check reflect type assertions
		k := reflect.TypeFor[Value1]()
		assert.Panics(t, func() { s.Has(k) }, "using Has with reflect.Type should panic")
		assert.Panics(t, func() { s.Get(k) }, "using Get with reflect.Type should panic")
		assert.Panics(t, func() { s.Put(k) }, "using Put with reflect.Type should panic")
		assert.Panics(t, func() { s.Del(k) }, "using Del with reflect.Type should panic")
	}
}
