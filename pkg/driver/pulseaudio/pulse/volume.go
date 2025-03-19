package pulse

import (
	"github.com/jfreymuth/pulse/proto"
)

func VolumeMax(channels proto.ChannelVolumes) proto.Volume {
	var out proto.Volume
	for _, v := range channels {
		if v > out {
			out = v
		}
	}
	return out
}
