package pulseaudio

import (
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/widget"
)

var _ entityProcessor = (*sourceOutputEntityProcessor)(nil)
var _ commandProcessor = (*sourceOutputCommandProcessor)(nil)

type sourceOutputEntityProcessor struct {
	c *Connection
}

func (s sourceOutputEntityProcessor) updateAll() {
	items, err := s.c.client.GetSourceOutputInfoList()
	if err != nil {
		s.c.log.Error("Failed to fetch source output info list: ", "error", err)
	}
	for _, item := range items {
		s.updateInfo(item)
	}
}

func (s sourceOutputEntityProcessor) updateId(index uint32) {
	info, err := s.c.client.GetSourceOutputInfoByIndex(index)
	if err != nil {
		s.c.log.Error("error while fetching source output info", "index", index, "error", err)
		return
	}
	s.updateInfo(info)
}

func (s sourceOutputEntityProcessor) removeId(index uint32) {
	name := s.entityName(index)
	s.c.removeEntity(name)
}

func (s sourceOutputEntityProcessor) entityName(index uint32) string {
	return fmt.Sprintf("pulseaudio-source-output-%d", index)
}

func (s sourceOutputEntityProcessor) updateInfo(info *proto.GetSourceOutputInfoReply) {
	index := info.SourceOutpuIndex
	s.c.log.Debug("processing source output info", "index", index, "info", info)
	s.c.updateEntity(
		s.entityName(index),
		sourceOutputCommandProcessor{s.c, index},
		processStreamComponent(widget.StreamTypeInput),
		processInfoComponent(info.MediaName, info.Properties),
		processVolumeComponent(info.ChannelVolumes, info.Muted),
	)
}

type sourceOutputCommandProcessor struct {
	c     *Connection
	index uint32
}

func (s sourceOutputCommandProcessor) registerCommands(ent *entity.Entity) {
	entity.SetHandler(ent, s.processVolumeChangeCommand)
	entity.SetHandler(ent, s.processVolumeMuteCommand)
}

func (s sourceOutputCommandProcessor) processVolumeChangeCommand(ent *entity.Entity, cmd widget.VolumeChangeCommand) error {
	level := proto.NormVolume(float64(cmd.Level))
	return s.c.client.Command(&proto.SetSourceOutputVolume{
		SourceOutputIndex: s.index,
		ChannelVolumes:    []proto.Volume{level},
	})
}
func (s sourceOutputCommandProcessor) processVolumeMuteCommand(ent *entity.Entity, cmd widget.VolumeMuteCommand) error {
	return s.c.client.Command(&proto.SetSourceOutputMute{
		SourceOutputIndex: s.index,
		Mute:              cmd.Mute,
	})
}
