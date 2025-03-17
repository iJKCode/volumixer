package entity

type EntityAddedEvent struct {
	Entity *Entity
}

type EntityRemovedEvent struct {
	Entity *Entity
}

type HandlersUpdatedEvent struct {
	Entity *Entity
}

type ComponentUpdatedEvent struct {
	Entity    *Entity
	Component any
}

type ComponentRemovedEvent struct {
	Entity    *Entity
	Component any
}
