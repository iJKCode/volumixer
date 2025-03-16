package pulseaudio

import (
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/widget"
)

var _ entityProcessor = &sinkInputEntityProcessor{}
var _ commandProcessor = &sinkInputCommandProcessor{}

type sinkInputEntityProcessor struct {
	c *Connection
}

func (s sinkInputEntityProcessor) updateAll() {
	items, err := s.c.client.GetSinkInputInfoList()
	if err != nil {
		s.c.log.Error("Failed to fetch sink input info list: ", "error", err)
	}
	for _, item := range items {
		s.updateInfo(item)
	}
}

func (s sinkInputEntityProcessor) updateId(index uint32) {
	info, err := s.c.client.GetSinkInputInfoByIndex(index)
	if err != nil {
		s.c.log.Error("error while fetching sink input info", "index", index, "error", err)
		return
	}
	s.updateInfo(info)
}

func (s sinkInputEntityProcessor) removeId(index uint32) {
	name := s.entityName(index)
	s.c.removeEntity(name)
}

func (s sinkInputEntityProcessor) entityName(index uint32) string {
	return fmt.Sprintf("pulseaudio-sink-input-%d", index)
}

func (s sinkInputEntityProcessor) updateInfo(info *proto.GetSinkInputInfoReply) {
	index := info.SinkInputIndex
	s.c.log.Debug("processing sink input info", "index", index, "info", info)
	s.c.updateEntity(
		s.entityName(index),
		sinkInputCommandProcessor{s.c, index},
		processStreamComponent(widget.StreamTypeOutput),
		processInfoComponent(info.MediaName, info.Properties),
		processVolumeComponent(info.ChannelVolumes, info.Muted),
	)
}

type sinkInputCommandProcessor struct {
	c     *Connection
	index uint32
}

func (s sinkInputCommandProcessor) registerCommands(ent *entity.Entity) {
	entity.SetHandler(ent, s.processVolumeChangeCommand)
	entity.SetHandler(ent, s.processVolumeMuteCommand)
}

func (s sinkInputCommandProcessor) processVolumeChangeCommand(ent *entity.Entity, cmd widget.VolumeChangeCommand) error {
	level := proto.NormVolume(float64(cmd.Level))
	return s.c.client.Command(&proto.SetSinkInputVolume{
		SinkInputIndex: s.index,
		ChannelVolumes: []proto.Volume{level},
	})
}
func (s sinkInputCommandProcessor) processVolumeMuteCommand(ent *entity.Entity, cmd widget.VolumeMuteCommand) error {
	return s.c.client.Command(&proto.SetSinkInputMute{
		SinkInputIndex: s.index,
		Mute:           cmd.Mute,
	})
}
