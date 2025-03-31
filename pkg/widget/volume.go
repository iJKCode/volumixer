package widget

import (
	"connectrpc.com/connect"
	"context"
	"github.com/google/uuid"
	commandv1 "github.com/ijkcode/volumixer-api/gen/go/command/v1"
	"github.com/ijkcode/volumixer-api/gen/go/command/v1/commandv1connect"
	widgetv1 "github.com/ijkcode/volumixer-api/gen/go/widget/v1"
	"github.com/ijkcode/volumixer/pkg/core/command"
	"github.com/ijkcode/volumixer/pkg/core/component"
	"github.com/ijkcode/volumixer/pkg/core/entity"
	"log/slog"
)

type VolumeComponent struct {
	Level float32
	Muted bool
}

type VolumeChangeCommand struct {
	Level float32
}

type VolumeMuteCommand struct {
	Mute bool
}

func init() {
	component.RegisterFatal(VolumeComponentMarshaller{})
}

type VolumeComponentMarshaller struct{}

var _ component.Marshaler[VolumeComponent, *widgetv1.VolumeComponent] = (*VolumeComponentMarshaller)(nil)

func (m VolumeComponentMarshaller) MarshalProto(cmp VolumeComponent) (*widgetv1.VolumeComponent, error) {
	return &widgetv1.VolumeComponent{
		Level: cmp.Level,
		Muted: cmp.Muted,
	}, nil
}

func (m VolumeComponentMarshaller) UnmarshalProto(msg *widgetv1.VolumeComponent) (VolumeComponent, error) {
	return VolumeComponent{
		Level: msg.Level,
		Muted: msg.Muted,
	}, nil
}

var _ commandv1connect.VolumeServiceHandler = (*VolumeServiceHandler)(nil)

type VolumeServiceHandler struct {
	Log      *slog.Logger
	Entities *entity.Context
}

func (v VolumeServiceHandler) SetVolumeLevel(ctx context.Context, req *connect.Request[commandv1.SetVolumeLevelRequest]) (*connect.Response[commandv1.SetVolumeLevelResponse], error) {
	v.Log.InfoContext(ctx, "got set volume level command", "entity", req.Msg.EntityId, "level", req.Msg.Level)
	id, err := uuid.Parse(req.Msg.EntityId)
	if err != nil {
		return nil, err
	}
	err = command.DispatchContext(v.Entities, id, VolumeChangeCommand{
		Level: req.Msg.Level,
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&commandv1.SetVolumeLevelResponse{}), nil
}

func (v VolumeServiceHandler) SetVolumeMute(ctx context.Context, req *connect.Request[commandv1.SetVolumeMuteRequest]) (*connect.Response[commandv1.SetVolumeMuteResponse], error) {
	v.Log.InfoContext(ctx, "got set volume mute command", "entity", req.Msg.EntityId, "level", req.Msg.Mute)
	id, err := uuid.Parse(req.Msg.EntityId)
	if err != nil {
		return nil, err
	}
	err = command.DispatchContext(v.Entities, id, VolumeMuteCommand{
		Mute: req.Msg.Mute,
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&commandv1.SetVolumeMuteResponse{}), nil
}
