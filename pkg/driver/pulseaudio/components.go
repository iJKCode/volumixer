package pulseaudio

import (
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"ijkcode.tech/volumixer/pkg/widget"
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
	var level float64
	if len(channels) > 0 {
		for _, channel := range channels {
			level += channel.Norm()
		}
		level /= float64(len(channels))
	}
	level = math.Round(level/VolumePrecision) * VolumePrecision
	return widget.VolumeComponent{
		Level: float32(level),
		Muted: muted,
	}
}
