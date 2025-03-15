package entity

import (
	"github.com/google/uuid"
	"ijkcode.tech/volumixer/pkg/util/typemap"
	"ijkcode.tech/volumixer/pkg/util/typeset"
	"sync"
	"sync/atomic"
)

type ID = uuid.UUID
type ComponentType = typemap.KeyType
type CommandType = typeset.ValType

type Entity struct {
	id         ID
	ctx        atomic.Pointer[Context]
	mut        sync.RWMutex
	components typeset.TypeSet
	handlers   typemap.TypeMap[Handler[any]]
}

func NewEntity(options ...func(*Entity)) *Entity {
	entity := &Entity{}
	for _, option := range options {
		option(entity)
	}
	if entity.id == uuid.Nil {
		entity.id = uuid.New()
	}
	return entity
}

func WithID(id ID) func(*Entity) {
	return func(entity *Entity) {
		entity.id = id
	}
}

func WithComponents(components []any) func(*Entity) {
	return func(entity *Entity) {
		entity.SetComponents(components...)
	}
}

func WithHandlers(handlers typemap.TypeMap[Handler[any]]) func(*Entity) {
	return func(entity *Entity) {
		for cmd, handler := range handlers.Items() {
			entity.SetHandlerType(cmd, handler)
		}
	}
}

func (e *Entity) ID() ID {
	return e.id
}

func (e *Entity) IsActive() bool {
	return e.ctx.Load() != nil
}

func (e *Entity) publishEvent(value any) {
	c := e.ctx.Load()
	if c == nil {
		return
	}
	c.publishEvent(value)
}
