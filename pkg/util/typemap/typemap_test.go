package typemap

import (
	"reflect"
	"testing"
)
import "github.com/stretchr/testify/assert"

type Value1 struct{}
type Value2 struct{}
type Value3 struct {
	V string
}
type ValueStr string

func TestMemberFunctions(t *testing.T) {
	m := Make[int]()

	// add some values
	m.Put(Value1{}, 42)
	m.Put(Value3{}, 99)
	m.Put(ValueStr(""), 69)

	{
		// check has functionality
		assert.True(t, m.Has(Value1{}), "should contain Value1")
		assert.False(t, m.Has(Value2{}), "should not contain Value2")
		assert.True(t, m.Has(Value3{}), "should contain Value3")
		assert.True(t, m.Has(ValueStr("")), "should contain ValueStr")
		assert.False(t, m.Has(""), "should not contain string")
	}
	{
		// check get functionality
		v, ok := m.Get(Value1{})
		assert.True(t, ok, "should contain Value1")
		assert.Equal(t, 42, v)
		v, ok = m.Get(Value2{})
		assert.False(t, ok, "should not contain Value2")
		v, ok = m.Get(Value3{})
		assert.True(t, ok, "should contain Value3")
		assert.Equal(t, 99, v)
		v, ok = m.Get(ValueStr(""))
		assert.True(t, ok, "should contain ValueStr")
		assert.Equal(t, 69, v)
		v, ok = m.Get("")
		assert.False(t, ok, "should not contain string")
	}
	{
		// check type with different contents
		v, ok := m.Get(Value3{})
		assert.True(t, ok, "should contain Value3")
		assert.Equal(t, 99, v)
		v, ok = m.Get(Value3{"test"})
		assert.True(t, ok, "should contain Value3")
		assert.Equal(t, 99, v)
	}
	{
		// check replace functionality
		m.Put(Value1{}, 101)
		v, ok := m.Get(Value1{})
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, 101, v)
		assert.Equal(t, 3, m.Len(), "length should be 3")
	}
	{
		// check delete functionality
		m.Del(Value1{})
		assert.False(t, m.Has(Value1{}), "Value1 should be removed")
		assert.True(t, m.Has(Value3{}), "Value3 should still exist")
		assert.Equal(t, 2, m.Len(), "length should be 2")
	}
	{
		// re-add deleted value
		m.Put(Value1{}, 67)
		v, ok := m.Get(Value1{})
		assert.True(t, ok, "Value1 should exist again")
		assert.Equal(t, 67, v)
		assert.Equal(t, 3, m.Len(), "length should be 3")
	}
	{
		// check all items functionality
		a := m.Items()
		assert.Equal(t, a, map[KeyType]int{
			reflect.TypeOf(Value1{}):     67,
			reflect.TypeOf(Value3{}):     99,
			reflect.TypeOf(ValueStr("")): 69,
		})
	}
	{
		// check reflect type assertions
		k := reflect.TypeFor[Value1]()
		assert.Panics(t, func() { m.Has(k) }, "using Has with reflect.Type should panic")
		assert.Panics(t, func() { m.Get(k) }, "using Get with reflect.Type should panic")
		assert.Panics(t, func() { m.Put(k, 0) }, "using Put with reflect.Type should panic")
		assert.Panics(t, func() { m.Del(k) }, "using Del with reflect.Type should panic")
	}
}
