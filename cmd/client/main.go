package main

import (
	"connectrpc.com/connect"
	"context"
	"google.golang.org/protobuf/types/known/anypb"
	"ijkcode.tech/volumixer/pkg/core/component"
	corev1 "ijkcode.tech/volumixer/proto/core/v1"
	"ijkcode.tech/volumixer/proto/core/v1/corev1connect"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	// register components
	_ "ijkcode.tech/volumixer/pkg/widget"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)

	meta := component.ListMetadata()
	slog.Info("registered components", "count", len(meta))

	client := corev1connect.NewEntityServiceClient(
		http.DefaultClient,
		"http://localhost:5000",
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	res, err := client.EventStream(ctx, connect.NewRequest(&corev1.EventStreamRequest{
		SimulateState: true,
	}))
	if err != nil {
		slog.Error("failed to connect event stream", "error", err.Error())
	} else {
		for res.Receive() {
			slog.Debug("got event", "event", res.Msg())
			switch evt := res.Msg().Event.(type) {
			case *corev1.EventStreamResponse_EntityAdded:
				slog.Info("entity added", "entity", evt.EntityAdded.EntityId)
				for _, msg := range evt.EntityAdded.Components {
					cmp := unmarshalComponent(msg)
					slog.Info("component updated", "entity", evt.EntityAdded.EntityId, "component", cmp)
				}
			case *corev1.EventStreamResponse_EntityRemoved:
				slog.Info("entity removed", "entity", evt.EntityRemoved.EntityId)
			case *corev1.EventStreamResponse_ComponentUpdated:
				cmp := unmarshalComponent(evt.ComponentUpdated.Component)
				slog.Info("component updated", "entity", evt.ComponentUpdated.EntityId, "component", cmp)
			case *corev1.EventStreamResponse_ComponentRemoved:
				cmp := unmarshalComponent(evt.ComponentRemoved.Component)
				slog.Info("component removed", "entity", evt.ComponentRemoved.EntityId, "component", cmp)
			default:
				slog.Info("invalid event type", "event", evt)
			}
		}
	}
}

func unmarshalComponent(msg *anypb.Any) any {
	cmp, _, err := component.UnmarshalProtoAny(msg)
	if err != nil {
		slog.Warn("failed to unmarshal component", "error", err.Error())
		return nil
	}
	return cmp
}
