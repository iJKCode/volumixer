package pulseaudio

import (
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/widget"
)

var _ entityProcessor = (*sourceEntityProcessor)(nil)
var _ commandProcessor = (*sourceCommandProcessor)(nil)

type sourceEntityProcessor struct {
	c *Connection
}

func (s sourceEntityProcessor) updateAll() {

	items, err := s.c.client.GetSourceInfoList()
	if err != nil {
		s.c.log.Error("Failed to fetch source info list: ", "error", err)
	}
	for _, item := range items {
		s.processInfo(item)
	}
}

func (s sourceEntityProcessor) updateId(index uint32) {
	info, err := s.c.client.GetSourceInfoByIndex(index)
	if err != nil {
		s.c.log.Error("error while fetching source info", "index", index, "error", err)
		return
	}
	s.processInfo(info)
}

func (s sourceEntityProcessor) removeId(index uint32) {
	name := s.entityName(index)
	s.c.removeEntity(name)
}

func (s sourceEntityProcessor) entityName(index uint32) string {
	return fmt.Sprintf("pulseaudio-source-%d", index)
}

func (s sourceEntityProcessor) processInfo(info *proto.GetSourceInfoReply) {
	index := info.SourceIndex
	s.c.log.Debug("processing source info", "index", index, "info", info)
	s.c.updateEntity(
		s.entityName(index),
		sourceCommandProcessor{s.c, index},
		processStreamComponent(widget.StreamTypeCapture),
		processInfoComponent(info.SourceName, info.Properties),
		processVolumeComponent(info.ChannelVolumes, info.Mute),
	)
}

type sourceCommandProcessor struct {
	c     *Connection
	index uint32
}

func (s sourceCommandProcessor) registerCommands(ent *entity.Entity) {
	entity.SetHandler(ent, s.processVolumeChangeCommand)
	entity.SetHandler(ent, s.processVolumeMuteCommand)
}

func (s sourceCommandProcessor) processVolumeChangeCommand(ent *entity.Entity, cmd widget.VolumeChangeCommand) error {
	level := proto.NormVolume(float64(cmd.Level))
	return s.c.client.Command(&proto.SetSourceVolume{
		SourceIndex:    s.index,
		ChannelVolumes: []proto.Volume{level},
	})
}
func (s sourceCommandProcessor) processVolumeMuteCommand(ent *entity.Entity, cmd widget.VolumeMuteCommand) error {
	return s.c.client.Command(&proto.SetSourceMute{
		SourceIndex: s.index,
		Mute:        cmd.Mute,
	})
}
