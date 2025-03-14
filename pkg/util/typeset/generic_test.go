package typeset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericFunctions(t *testing.T) {
	m := Make()
	p := &m

	// add some values
	Put(p, Value1{42})
	Put(p, Value2{"hello"})
	Put(p, ValueStr("world"))

	{
		// check has functionality
		assert.True(t, Has[Value1](p), "should contain Value1")
		assert.True(t, Has[Value2](p), "should contain Value2")
		assert.False(t, Has[Value3](p), "should not contain Value3")
		assert.True(t, Has[ValueStr](p), "should contain ValueStr")
	}
	{
		// check get functionality
		v1, ok := Get[Value1](p)
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, Value1{42}, v1)
		v2, ok := Get[Value2](p)
		assert.True(t, ok, "Value2 should exist")
		assert.Equal(t, Value2{"hello"}, v2)
		_, ok = Get[Value3](p)
		assert.False(t, ok, "Value3 should not exist")
		v4, ok := Get[ValueStr](p)
		assert.True(t, ok, "ValueStr should exist")
		assert.Equal(t, ValueStr("world"), v4)
	}
	{
		// check replace functionality
		Put[Value1](p, Value1{99})
		v, ok := Get[Value1](p)
		assert.True(t, ok, "Value1 should exist")
		assert.Equal(t, Value1{99}, v)
	}
	{
		// check delete functionality
		Del[Value1](p)
		assert.False(t, Has[Value1](p), "Value1 should be removed")
		assert.True(t, Has[Value2](p), "Value2 should still exist")
	}
	{
		// re-add deleted value
		Put[Value1](p, Value1{67})
		v, ok := Get[Value1](p)
		assert.True(t, ok, "Value1 should exist again")
		assert.Equal(t, Value1{67}, v)
	}
}
