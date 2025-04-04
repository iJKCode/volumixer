package entity

import (
	"github.com/ijkcode/volumixer/pkg/util/valset"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextRelations(t *testing.T) {
	c := NewContext(nil)
	s := c.storage.Load()

	c1 := c.SubContext()
	c2 := c.SubContext()
	c11 := c1.SubContext()
	c3 := c.SubContext()

	// check collectChildren
	assert.Equal(t, valset.From(c1, c2, c11, c3), c.collectChildren(nil))
	assert.Equal(t, valset.From(c11), c1.collectChildren(nil))
	assert.Empty(t, c2.collectChildren(nil))
	assert.Empty(t, c11.collectChildren(nil))
	assert.Empty(t, c3.collectChildren(nil))

	// check context parents
	assert.Nil(t, c.parent)
	assert.Equal(t, c, c1.parent)
	assert.Equal(t, c, c2.parent)
	assert.Equal(t, c1, c11.parent)
	assert.Equal(t, c, c3.parent)
	// check context children
	assert.Equal(t, valset.From(c1, c2, c3), c.children)
	assert.Equal(t, valset.From(c11), c1.children)
	assert.Empty(t, c2.children)
	assert.Empty(t, c11.children)
	assert.Empty(t, c3.children)
	// check storage references
	assert.Equal(t, s, c1.storage.Load())
	assert.Equal(t, s, c2.storage.Load())
	assert.Equal(t, s, c11.storage.Load())
	assert.Equal(t, s, c3.storage.Load())

	// close leaf context
	assert.NoError(t, c2.Close())
	assert.Nil(t, c2.parent)
	assert.Empty(t, c2.children)
	assert.NotEqual(t, s, c2.storage.Load())
	// check other contexts
	assert.Equal(t, s, c1.storage.Load())
	assert.Equal(t, s, c11.storage.Load())
	assert.Equal(t, s, c3.storage.Load())
	assert.Equal(t, c, c1.parent)
	assert.Equal(t, c1, c11.parent)
	assert.Equal(t, c, c3.parent)
	assert.Equal(t, valset.From(c1, c3), c.children)
	assert.Equal(t, valset.From(c11), c1.children)
	assert.Empty(t, c3.children)

	// close internal context
	assert.NoError(t, c1.Close())
	assert.Nil(t, c1.parent)
	assert.Empty(t, c1.children)
	assert.NotEqual(t, s, c1.storage.Load())
	assert.Nil(t, c11.parent)
	assert.Empty(t, c11.children)
	assert.NotEqual(t, s, c11.storage.Load())
	assert.False(t, c.children.Has(c1))

	// close root context
	assert.NoError(t, c.Close())
	assert.Nil(t, c.storage.Load())
	assert.Nil(t, c3.storage.Load())
	assert.Nil(t, c.parent)
	assert.Nil(t, c3.parent)
	assert.Empty(t, c.children)
	assert.Empty(t, c3.children)
}

func TestContextEntities(t *testing.T) {
	c := NewContext(nil)

	c1 := c.SubContext()
	c2 := c.SubContext()
	c11 := c1.SubContext()
	c3 := c.SubContext()

	// add some entities
	ce1 := NewEntity()
	assert.NoError(t, c.Add(ce1))
	c1e1 := NewEntity()
	assert.NoError(t, c1.Add(c1e1))
	c1e2 := NewEntity()
	assert.NoError(t, c1.Add(c1e2))
	c11e1 := NewEntity()
	assert.NoError(t, c11.Add(c11e1))
	c2e1 := NewEntity()
	assert.NoError(t, c2.Add(c2e1))

	// check direct entity ownership
	assert.False(t, c1.Owns(ce1))
	assert.True(t, c1.Owns(c1e1))
	assert.True(t, c1.Owns(c1e2))
	assert.False(t, c1.Owns(c11e1))
	assert.False(t, c1.Owns(c2e1))
	assert.True(t, c11.Owns(c11e1))
	assert.False(t, c11.Owns(c1e1))
	assert.False(t, c11.Owns(c2e1))

	// check nested entity ownership
	assert.False(t, c1.OwnsDeep(ce1))
	assert.True(t, c1.OwnsDeep(c1e1))
	assert.True(t, c1.OwnsDeep(c1e2))
	assert.True(t, c1.OwnsDeep(c11e1))
	assert.False(t, c1.OwnsDeep(c2e1))
	assert.True(t, c11.OwnsDeep(c11e1))
	assert.False(t, c11.OwnsDeep(c1e1))
	assert.False(t, c11.OwnsDeep(c2e1))

	// close empty context
	assert.NoError(t, c3.Close())
	assert.True(t, c.Has(ce1))
	assert.True(t, c.Has(c1e1))
	assert.True(t, c.Has(c1e2))
	assert.True(t, c.Has(c11e1))
	assert.True(t, c.Has(c2e1))

	// attempt close after close
	assert.NoError(t, c3.Close())

	// close inner context
	assert.NoError(t, c1.Close())
	assert.True(t, c.Has(ce1))
	assert.False(t, c.Has(c1e1))
	assert.False(t, c.Has(c1e2))
	assert.False(t, c.Has(c11e1))
	assert.True(t, c.Has(c2e1))

	// close root context
	assert.NoError(t, c.Close())
	assert.False(t, c.Has(ce1))
	assert.False(t, c.Has(c1e1))
	assert.False(t, c.Has(c1e2))
	assert.False(t, c.Has(c11e1))
	assert.False(t, c.Has(c2e1))
}

func TestNamedEntities(t *testing.T) {

	c := NewContext(nil)

	c1 := c.SubContext()
	c2 := c.SubContext()
	c11 := c1.SubContext()

	// add some entities
	ce1 := NewEntity()
	assert.NoError(t, c.Add(ce1))
	c1e1 := NewEntity(WithName("c1e1"))
	assert.NoError(t, c1.Add(c1e1))
	c1e2 := NewEntity(WithName("c1e2"))
	assert.NoError(t, c1.Add(c1e2))
	c1e3 := NewEntity()
	assert.NoError(t, c1.Add(c1e2))
	c11e1 := NewEntity(WithName("c11e1"))
	assert.NoError(t, c11.Add(c11e1))
	c2e1 := NewEntity()
	assert.NoError(t, c2.Add(c2e1))

	// check name properties
	assert.True(t, c1e1.HasName())
	assert.Equal(t, "c1e1", c1e1.Name())
	assert.False(t, c1e3.HasName())

	// check context has named
	assert.False(t, c1.HasNamed(ce1.Name()))
	assert.True(t, c1.HasNamed(c1e1.Name()))
	assert.True(t, c1.HasNamed(c1e2.Name()))
	assert.False(t, c1.HasNamed(c1e3.Name()))
	assert.False(t, c1.HasNamed(c11e1.Name()))
	assert.False(t, c1.HasNamed(c2e1.Name()))
	assert.True(t, c11.HasNamed(c11e1.Name()))
	assert.False(t, c11.HasNamed(c1e1.Name()))
	assert.False(t, c11.HasNamed(c2e1.Name()))

	// check context get named
	ent, ok := c1.GetNamed(c1e1.Name())
	assert.True(t, ok)
	assert.Equal(t, c1e1, ent)
	ent, ok = c1.GetNamed(c1e2.Name())
	assert.True(t, ok)
	assert.Equal(t, c1e2, ent)
	ent, ok = c11.GetNamed(c11e1.Name())
	assert.True(t, ok)
	assert.Equal(t, c11e1, ent)

	// close inner context
	assert.NoError(t, c1.Close())
	assert.False(t, c1.HasNamed(c1e1.Name()))
	assert.False(t, c1.HasNamed(c1e2.Name()))
	assert.False(t, c1.HasNamed(c1e3.Name()))
	assert.False(t, c11.HasNamed(c11e1.Name()))
}
