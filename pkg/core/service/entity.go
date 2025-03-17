package service

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/core/event"
	"ijkcode.tech/volumixer/pkg/widget"
	corev1 "ijkcode.tech/volumixer/proto/core/v1"
	widgetv1 "ijkcode.tech/volumixer/proto/widget/v1"
	"log/slog"
)

type EntityServiceHandler struct {
	Entities *entity.Context
	Log      *slog.Logger
}

func (h EntityServiceHandler) EntityList(ctx context.Context, req *connect.Request[corev1.EntityListRequest]) (*connect.Response[corev1.EntityListResponse], error) {
	h.Log.InfoContext(ctx, "got entity list request")
	entities := make([]*corev1.EntityInfo, 0)
	for ent := range h.Entities.All() {
		entities = append(entities, &corev1.EntityInfo{
			EntityId:   ent.ID().String(),
			Components: h.componentsToAny(ent.Components()),
		})
	}
	return connect.NewResponse(&corev1.EntityListResponse{
		Entities: entities,
	}), nil
}

func (h EntityServiceHandler) EntityById(ctx context.Context, req *connect.Request[corev1.EntityByIdRequest]) (*connect.Response[corev1.EntityByIdResponse], error) {
	h.Log.InfoContext(ctx, "got entity info request")
	id, err := uuid.Parse(req.Msg.EntityId)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}
	ent, ok := h.Entities.Get(id)
	if !ok {
		return nil, errors.New("entity not found")
	}
	return connect.NewResponse(&corev1.EntityByIdResponse{
		Entity: &corev1.EntityInfo{
			EntityId:   ent.ID().String(),
			Components: h.componentsToAny(ent.Components()),
		},
	}), nil
}

func (h EntityServiceHandler) EventStream(ctx context.Context, req *connect.Request[corev1.EventStreamRequest], res *connect.ServerStream[corev1.EventStreamResponse]) error {
	h.Log.InfoContext(ctx, "got event stream request")
	streamEvents := make(chan *corev1.EventStreamResponse, 10) //TODO is buffering needed ???

	bus := h.Entities.EventBus()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.EntityAddedEvent) {
		streamEvents <- &corev1.EventStreamResponse{
			Event: &corev1.EventStreamResponse_EntityAdded{
				EntityAdded: &corev1.EntityAddedEvent{
					EntityId:   evt.Entity.ID().String(),
					Components: h.componentsToAny(evt.Entity.Components()),
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.EntityRemovedEvent) {
		streamEvents <- &corev1.EventStreamResponse{
			Event: &corev1.EventStreamResponse_EntityRemoved{
				EntityRemoved: &corev1.EntityRemovedEvent{
					EntityId: evt.Entity.ID().String(),
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.ComponentUpdatedEvent) {
		component := h.componentToAny(evt.Component)
		if component == nil {
			return
		}
		streamEvents <- &corev1.EventStreamResponse{
			Event: &corev1.EventStreamResponse_ComponentUpdated{
				ComponentUpdated: &corev1.ComponentUpdatedEvent{
					EntityId:  evt.Entity.ID().String(),
					Component: component,
				},
			},
		}
	})()
	defer event.SubscribeFunc(bus, func(ctx context.Context, evt entity.ComponentRemovedEvent) {
		component := h.componentToAny(evt.Component)
		if component == nil {
			return
		}
		streamEvents <- &corev1.EventStreamResponse{
			Event: &corev1.EventStreamResponse_ComponentRemoved{
				ComponentRemoved: &corev1.ComponentRemovedEvent{
					EntityId:  evt.Entity.ID().String(),
					Component: component,
				},
			},
		}
	})()

	if req.Msg.SimulateState {
		for ent := range h.Entities.All() {
			evt := &corev1.EventStreamResponse{
				Event: &corev1.EventStreamResponse_EntityAdded{
					EntityAdded: &corev1.EntityAddedEvent{
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
	for _, component := range components {
		anyComponent := h.componentToAny(component)
		if anyComponent != nil {
			result = append(result, anyComponent)
		}
	}
	return result
}

func (h EntityServiceHandler) componentToAny(component any) *anypb.Any {
	protoComponent := h.componentToProto(component)
	if protoComponent == nil {
		return nil
	}
	anyComponent, err := anypb.New(protoComponent)
	if err != nil {
		h.Log.Warn("error marshaling component", "error", err)
		return nil
	}
	return anyComponent
}

func (h EntityServiceHandler) componentToProto(component any) proto.Message {
	switch component := component.(type) {
	case widget.InfoComponent:
		return &widgetv1.InfoComponent{
			Name: component.Name,
		}
	case widget.VolumeComponent:
		return &widgetv1.VolumeComponent{
			Level: component.Level,
			Muted: component.Muted,
		}
	default:
		//TODO implement dynamic component to proto conversion
		return nil
	}
}
