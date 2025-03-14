package entity

import (
	"errors"
	"sync"
)

var ErrEntityDuplicateID = errors.New("duplicate entity id")

type sharedStorage struct {
	mut      sync.RWMutex
	entities map[ID]*Entity
}

func (s *sharedStorage) add(ent *Entity, ctx *Context) error {
	id := ent.ID()
	if e, ok := s.entities[id]; ok && (e != ent || e.ctx.Load() != ctx) {
		return ErrEntityDuplicateID
	}

	s.entities[id] = ent
	ent.ctx.Store(ctx)

	//TODO trigger events removed entity

	return nil
}

func (s *sharedStorage) remove(ent *Entity) error {
	id := ent.ID()
	e, ok := s.entities[id]
	if !ok || e != ent {
		// trying to remove a different or nonexistent entity
		return nil
	}

	delete(s.entities, id)
	ent.ctx.Store(nil)

	//TODO trigger events for added entity

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
