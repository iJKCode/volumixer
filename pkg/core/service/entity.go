package service

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	entityv1 "github.com/ijkcode/volumixer-api/gen/go/entity/v1"
	"github.com/ijkcode/volumixer/pkg/core/component"
	"github.com/ijkcode/volumixer/pkg/core/entity"
	"github.com/ijkcode/volumixer/pkg/core/event"
	"google.golang.org/protobuf/types/known/anypb"
	"log/slog"
)

type EntityServiceHandler struct {
	Entities *entity.Context
	Log      *slog.Logger
}

func (h EntityServiceHandler) EntityList(ctx context.Context, req *connect.Request[entityv1.EntityListRequest]) (*connect.Response[entityv1.EntityListResponse], error) {
	h.Log.InfoContext(ctx, "got entity list request")
	entities := make([]*entityv1.EntityInfo, 0)
	for ent := range h.Entities.All() {
		entities = append(entities, &entityv1.EntityInfo{
			EntityId:   ent.ID().String(),
			Components: h.componentsToAny(ent.Components()),
		})
	}
	return connect.NewResponse(&entityv1.EntityListResponse{
		Entities: entities,
	}), nil
}

func (h EntityServiceHandler) EntityById(ctx context.Context, req *connect.Request[entityv1.EntityByIdRequest]) (*connect.Response[entityv1.EntityByIdResponse], error) {
	h.Log.InfoContext(ctx, "got entity info request")
	id, err := uuid.Parse(req.Msg.EntityId)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}
	ent, ok := h.Entities.Get(id)
	if !ok {
		return nil, errors.New("entity not found")
	}
	return connect.NewResponse(&entityv1.EntityByIdResponse{
		Entity: &entityv1.EntityInfo{
			EntityId:   ent.ID().String(),
			Components: h.componentsToAny(ent.Components()),
		},
	}), nil
}

func (h EntityServiceHandler) EventStream(ctx context.Context, req *connect.Request[entityv1.EventStreamRequest], res *connect.ServerStream[entityv1.EventStreamResponse]) error {
	h.Log.InfoContext(ctx, "got event stream request")
	streamEvents := make(chan *entityv1.EventStreamResponse, 10) //TODO is buffering needed ???

	bus := h.Entities.EventBus()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.EntityAddedEvent) {
		streamEvents <- &entityv1.EventStreamResponse{
			Event: &entityv1.EventStreamResponse_EntityAdded{
				EntityAdded: &entityv1.EntityAddedEvent{
					EntityId:   evt.Entity.ID().String(),
					Components: h.componentsToAny(evt.Entity.Components()),
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.EntityRemovedEvent) {
		streamEvents <- &entityv1.EventStreamResponse{
			Event: &entityv1.EventStreamResponse_EntityRemoved{
				EntityRemoved: &entityv1.EntityRemovedEvent{
					EntityId: evt.Entity.ID().String(),
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.ComponentUpdatedEvent) {
		cmp := h.componentToAny(evt.Component)
		if cmp == nil {
			return
		}
		streamEvents <- &entityv1.EventStreamResponse{
			Event: &entityv1.EventStreamResponse_ComponentUpdated{
				ComponentUpdated: &entityv1.ComponentUpdatedEvent{
					EntityId:  evt.Entity.ID().String(),
					Component: cmp,
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.ComponentRemovedEvent) {
		cmp := h.componentToAny(evt.Component)
		if cmp == nil {
			return
		}
		streamEvents <- &entityv1.EventStreamResponse{
			Event: &entityv1.EventStreamResponse_ComponentRemoved{
				ComponentRemoved: &entityv1.ComponentRemovedEvent{
					EntityId:  evt.Entity.ID().String(),
					Component: cmp,
				},
			},
		}
	})()

	if req.Msg.SimulateState {
		for ent := range h.Entities.All() {
			evt := &entityv1.EventStreamResponse{
				Event: &entityv1.EventStreamResponse_EntityAdded{
					EntityAdded: &entityv1.EntityAddedEvent{
						EntityId:   ent.ID().String(),
						Components: h.componentsToAny(ent.Components()),
					},
				},
			}
			err := res.Send(evt)
			if err != nil {
				h.Log.WarnContext(ctx, "error sending event to stream", "error", err)
				return nil
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case evt, ok := <-streamEvents:
			if !ok {
				return nil
			}
			if evt == nil {
				continue
			}
			err := res.Send(evt)
			if err != nil {
				h.Log.WarnContext(ctx, "error sending event to stream", "error", err)
				return nil
			}
		}
	}
}

func (h EntityServiceHandler) componentsToAny(components []any) []*anypb.Any {
	var result []*anypb.Any
	for _, cmp := range components {
		msg := h.componentToAny(cmp)
		if msg != nil {
			result = append(result, msg)
		}
	}
	return result
}

func (h EntityServiceHandler) componentToAny(cmp any) *anypb.Any {
	msg, meta, err := component.MarshalProtoAny(cmp)
	if err != nil {
		if errors.Is(err, component.ErrUnknownComponent) {
			return nil
		}
		h.Log.Warn("failed to marshal component", "error", err, "type", fmt.Sprintf("%T", cmp))
		return nil
	}
	if !meta.ShouldReplicate() {
		return nil
	}
	return msg
}
