package entity

import (
	"errors"
	"fmt"
	"ijkcode.tech/volumixer/pkg/core/event"
	"ijkcode.tech/volumixer/pkg/util/valset"
	"sync"
	"sync/atomic"
)

var ErrContextClosed = errors.New("registry context is closed")

type Context struct {
	storage  atomic.Pointer[sharedStorage]
	parent   *Context
	children valset.Set[*Context]
}

func NewContext(events *event.Bus) *Context {
	ctx := &Context{
		parent:   nil,
		children: valset.Make[*Context](),
	}
	ctx.storage.Store(&sharedStorage{
		mut:      sync.RWMutex{},
		events:   events,
		entities: make(map[ID]*Entity),
	})
	return ctx
}

func (c *Context) SubContext() *Context {
	s, unlock := c.getStorageWrite()
	defer unlock()
	if s == nil {
		return nil
	}

	child := &Context{
		parent:   c,
		children: valset.Make[*Context](),
	}
	child.storage.Store(s)
	c.children.Put(child)

	return child
}

func (c *Context) Close() error {
	// no-op if context already closed
	s, unlock := c.getStorageWrite()
	defer unlock()
	if s == nil {
		return nil
	}

	// get self and all child contexts
	contexts := valset.Make[*Context]()
	contexts.Put(c)
	c.collectChildren(contexts)

	// collect entities to remove
	remove := valset.Make[*Entity]()
	for _, ent := range s.entities {
		if contexts.Has(ent.ctx.Load()) {
			remove.Put(ent)
		}
	}

	// remove entities from storage
	var errs []error
	for ent := range remove {
		err := s.remove(ent)
		if err != nil {
			err := fmt.Errorf("removing entity %s", ent.ID())
			errs = append(errs, err)
		}
	}

	// unlink context relations
	if c.parent != nil {
		c.parent.children.Del(c)
	}
	for ctx := range contexts {
		ctx.storage.Swap(nil)
		ctx.parent = nil
		ctx.children = nil
	}

	if errs != nil {
		return errors.Join(errs...)
	}
	return nil
}

func (c *Context) IsActive() bool {
	return c.storage.Load() != nil
}

func (c *Context) EventBus() *event.Bus {
	s := c.storage.Load()
	if s == nil {
		return nil
	}
	return s.events
}

func (c *Context) publishEvent(value any) {
	s := c.storage.Load()
	if s == nil {
		return
	}
	event.Publish(s.events, value)
}

func (c *Context) Get(id ID) (*Entity, bool) {
	s, unlock := c.getStorageRead()
	defer unlock()
	if s == nil {
		return nil, false
	}

	ent, ok := s.entities[id]
	return ent, ok
}

func (c *Context) Add(ent *Entity) error {
	s, unlock := c.getStorageWrite()
	defer unlock()
	if s == nil {
		return ErrContextClosed
	}

	return s.add(ent, c)
}

func (c *Context) Remove(ent *Entity) error {
	s, unlock := c.getStorageWrite()
	defer unlock()
	if s == nil {
		return ErrContextClosed
	}

	return s.remove(ent)
}

func (c *Context) RemoveId(id ID) error {
	s, unlock := c.getStorageWrite()
	defer unlock()
	if s == nil {
		return ErrContextClosed
	}

	ent, ok := s.entities[id]
	if !ok {
		return nil
	}

	return s.remove(ent)
}

func (c *Context) Has(ent *Entity) bool {
	e, ok := c.Get(ent.ID())
	return ok && e == ent
}

func (c *Context) HasId(id ID) bool {
	_, ok := c.Get(id)
	return ok
}

func (c *Context) Owns(ent *Entity) bool {
	s, unlock := c.getStorageRead()
	defer unlock()
	if s == nil {
		return false
	}

	contexts := c.collectChildren(nil)
	contexts.Put(c)

	return contexts.Has(ent.ctx.Load())
}

func (c *Context) collectChildren(s valset.Set[*Context]) valset.Set[*Context] {
	if s == nil {
		s = valset.Make[*Context]()
	}
	for child := range c.children {
		s.Put(child)
		child.collectChildren(s)
	}
	return s
}
