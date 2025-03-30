package pulseaudio

import (
	"fmt"
	"github.com/ijkcode/volumixer/pkg/core/entity"
	"github.com/ijkcode/volumixer/pkg/widget"
	"github.com/jfreymuth/pulse/proto"
)

var _ entityProcessor = (*sinkEntityProcessor)(nil)
var _ commandProcessor = (*sinkCommandProcessor)(nil)

type sinkEntityProcessor struct {
	c *Connection
}

func (s sinkEntityProcessor) updateAll() {

	items, err := s.c.client.GetSinkInfoList()
	if err != nil {
		s.c.log.Error("Failed to fetch sink info list: ", "error", err)
	}
	for _, item := range items {
		s.processInfo(item)
	}
}

func (s sinkEntityProcessor) updateId(index uint32) {
	info, err := s.c.client.GetSinkInfoByIndex(index)
	if err != nil {
		s.c.log.Error("error while fetching sink info", "index", index, "error", err)
		return
	}
	s.processInfo(info)
}

func (s sinkEntityProcessor) removeId(index uint32) {
	name := s.entityName(index)
	s.c.removeEntity(name)
}

func (s sinkEntityProcessor) entityName(index uint32) string {
	return fmt.Sprintf("pulseaudio-sink-%d", index)
}

func (s sinkEntityProcessor) processInfo(info *proto.GetSinkInfoReply) {
	index := info.SinkIndex
	s.c.log.Debug("processing sink info", "index", index, "info", info)
	s.c.updateEntity(
		s.entityName(index),
		sinkCommandProcessor{s.c, index},
		processStreamComponent(widget.StreamTypePlayback),
		processInfoComponent(info.SinkName, info.Properties),
		processVolumeComponent(info.ChannelVolumes, info.Mute),
	)
}

type sinkCommandProcessor struct {
	c     *Connection
	index uint32
}

func (s sinkCommandProcessor) registerCommands(ent *entity.Entity) {
	entity.SetHandler(ent, s.processVolumeChangeCommand)
	entity.SetHandler(ent, s.processVolumeMuteCommand)
}

func (s sinkCommandProcessor) processVolumeChangeCommand(ent *entity.Entity, cmd widget.VolumeChangeCommand) error {
	level := proto.NormVolume(float64(cmd.Level))
	return s.c.client.Command(&proto.SetSinkVolume{
		SinkIndex:      s.index,
		ChannelVolumes: []proto.Volume{level},
	})
}
func (s sinkCommandProcessor) processVolumeMuteCommand(ent *entity.Entity, cmd widget.VolumeMuteCommand) error {
	return s.c.client.Command(&proto.SetSinkMute{
		SinkIndex: s.index,
		Mute:      cmd.Mute,
	})
}
