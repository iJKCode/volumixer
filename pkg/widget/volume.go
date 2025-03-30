package widget

import (
	widgetv1 "github.com/ijkcode/volumixer-api/gen/go/widget/v1"
	"github.com/ijkcode/volumixer/pkg/core/component"
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
