package pulseaudio

import (
	"fmt"
	"github.com/ijkcode/volumixer/pkg/driver/pulseaudio/pulse"
	"github.com/ijkcode/volumixer/pkg/widget"
	"github.com/jfreymuth/pulse/proto"
	"math"
)

const VolumePrecision = 0.001

const PropApplicationName = "application.name"

func processDriverComponent(conn *Connection) widget.DriverComponent {
	instance := conn.uri
	if instance == "" {
		instance = "local"
	}
	return widget.DriverComponent{
		Driver:   "pulseaudio",
		Instance: instance,
	}
}

func processStreamComponent(typ widget.StreamType) widget.StreamComponent {
	return widget.StreamComponent{
		Type: typ,
	}
}

func processInfoComponent(name string, info proto.PropList) widget.InfoComponent {
	appname, ok := info[PropApplicationName]
	if ok {
		name = fmt.Sprintf("%s - %s", appname.String(), name)
	}
	return widget.InfoComponent{
		Name: name,
	}
}

func processVolumeComponent(channels proto.ChannelVolumes, muted bool) widget.VolumeComponent {
	level := pulse.VolumeMax(channels).Norm()
	level = math.Round(level/VolumePrecision) * VolumePrecision
	return widget.VolumeComponent{
		Level: float32(level),
		Muted: muted,
	}
}
