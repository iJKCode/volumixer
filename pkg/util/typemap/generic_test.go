package typemap

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGenericFunctions(t *testing.T) {
	m := Make[int]()
	p := &m

	// add some values
	Put[Value1](p, 42)
	Put[Value3](p, 99)
	Put[ValueStr](p, 69)

	{
		// check has functionality
		assert.True(t, Has[Value1](p), "should contain Value1")
		assert.False(t, Has[Value2](p), "should not contain Value2")
		assert.True(t, Has[Value3](p), "should contain Value3")
		assert.True(t, Has[ValueStr](p), "should contain ValueStr")
		assert.False(t, Has[string](p), "should not contain string")
	}
	{
		// check get functionality
		v, ok := Get[Value1](p)
		assert.True(t, ok, "should contain Value1")
		assert.Equal(t, 42, v)
		v, ok = Get[Value2](p)
		assert.False(t, ok, "should not contain Value2")
		v, ok = Get[Value3](p)
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
		v, ok := Get[Value3](p)
		assert.True(t, ok, "should contain Value3")
		assert.Equal(t, 99, v)
		v, ok = m.Get(Value3{"test"})
		assert.True(t, ok, "should contain Value3")
		assert.Equal(t, 99, v)
	}
	{
		// check replace functionality
		Put[Value1](p, 101)
		v, ok := Get[Value1](p)
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, 101, v)
	}
	{
		// check delete functionality
		Del[Value1](p)
		assert.False(t, Has[Value1](p), "Value1 should be removed")
		assert.True(t, Has[Value3](p), "Value3 should still exist")
	}
	{
		// re-add deleted value
		Put[Value1](p, 67)
		v, ok := Get[Value1](p)
		assert.True(t, ok, "Value1 should exist again")
		assert.Equal(t, 67, v)
	}
	{
		// check all items functionality
		a := m.Items()
		assert.Equal(t, a, map[KeyType]int{
			reflect.TypeFor[Value1]():   67,
			reflect.TypeFor[Value3]():   99,
			reflect.TypeFor[ValueStr](): 69,
		})
	}
}
