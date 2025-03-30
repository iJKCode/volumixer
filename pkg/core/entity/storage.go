package entity

import (
	"errors"
	"github.com/ijkcode/volumixer/pkg/core/event"
	"sync"
)

var ErrEntityDuplicateID = errors.New("duplicate entity id")

type sharedStorage struct {
	mut      sync.RWMutex
	events   *event.Bus
	entities map[ID]*Entity
}

func (s *sharedStorage) add(ent *Entity, ctx *Context) error {
	id := ent.ID()
	e, ok := s.entities[id]
	if ok {
		if e == ent && e.ctx.Load() == ctx {
			// entity is already bound under the specified context
			return nil
		}
		return ErrEntityDuplicateID
	}
	if ent.HasName() {
		if e, ok := ctx.named[ent.name]; ok && (e != ent) {
			return ErrEntityDuplicateName
		}
	}

	s.entities[id] = ent
	ent.ctx.Store(ctx)
	if ent.HasName() {
		ctx.named[ent.name] = ent
	}

	event.Publish(s.events, EntityAddedEvent{
		Entity: ent,
	})

	return nil
}

func (s *sharedStorage) remove(ent *Entity) error {
	id := ent.ID()
	e, ok := s.entities[id]
	if !ok || e != ent {
		// trying to remove a different or nonexistent entity
		return nil
	}

	if ent.HasName() {
		ctx := ent.ctx.Load()
		if ctx != nil {
			delete(ctx.named, ent.name)
		}
	}

	delete(s.entities, id)
	ent.ctx.Store(nil)

	event.Publish(s.events, EntityRemovedEvent{
		Entity: ent,
	})

	return nil
}

func (c *Context) getStorageRead() (*sharedStorage, func()) {
	s := c.storage.Load()
	if s == nil {
		return nil, func() {}
	}

	s.mut.RLock()
	return s, func() {
		s.mut.RUnlock()
	}
}

func (c *Context) getStorageWrite() (*sharedStorage, func()) {
	s := c.storage.Load()
	if s == nil {
		return nil, func() {}
	}

	s.mut.Lock()
	return s, func() {
		s.mut.Unlock()
	}
}
