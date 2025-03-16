package pulseaudio

import (
	"context"
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/driver/pulseaudio/pulse"
	"log/slog"
)

const EventChanSize = 20

type Connection struct {
	log      *slog.Logger
	uri      string
	client   *pulse.Client
	entities *entity.Context
}

func NewConnection(log *slog.Logger, entities *entity.Context, uri string) *Connection {
	log = log.With("driver", "pulseaudio", "server", uri)
	return &Connection{
		log:      log,
		uri:      uri,
		client:   nil,
		entities: entities,
	}
}

func (c *Connection) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := pulse.NewClient()
	if err != nil {
		return err
	}

	defer func() {
		err := c.client.Close()
		if err != nil {
			c.log.Warn("error while closing connection", "error", err)
		}
	}()

	c.client = client
	c.log.Info("pulseaudio client info", "version", c.client.GetClientVersion())

	info, err := c.client.GetServerInfo()
	if err != nil {
		return fmt.Errorf("failed to fetch server info: %w", err)
	}

	c.log.Info("pulseaudio server info",
		"version", info.PackageVersion,
		"package", info.PackageName,
		"username", info.Username,
		"hostname", info.Hostname,
	)

	err = c.client.SetEventSubscription(proto.SubscriptionMaskAll)
	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	eventCh := make(chan any, EventChanSize)
	c.client.SetEventCallback(func(event any) {
		select {
		case eventCh <- event:
		default:
			c.log.Warn("failed to push event to queue", "event", event)
		}
	})

	c.log.Info("pulseaudio connection initialised")

	sinkEntityProcessor{c}.updateAll()
	sourceEntityProcessor{c}.updateAll()
	sinkInputEntityProcessor{c}.updateAll()
	sourceOutputEntityProcessor{c}.updateAll()

	c.processEvents(ctx, eventCh)

	c.log.Info("pulseaudio connection closed")

	return nil
}

type entityProcessor interface {
	updateAll()
	updateId(uint32)
	removeId(uint32)
	entityName(uint32) string
}

type commandProcessor interface {
	registerCommands(ent *entity.Entity)
}

func (c *Connection) processEvents(ctx context.Context, eventCh chan any) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-eventCh:
			if !ok {
				return
			}

			c.log.Debug("handling pulseaudio message", "event", event)

			switch event := event.(type) {

			case *proto.ConnectionClosed:
				c.log.Warn("received connection closed event", "event", event)
				return

			case *proto.SubscribeEvent:
				c.switchSubscriberEventFacility(event)
			}
		}
	}
}

func (c *Connection) switchSubscriberEventFacility(event *proto.SubscribeEvent) {
	switch event.Event & proto.EventFacilityMask {
	case proto.EventSink:
		c.switchSubscribeEventType(event, sinkEntityProcessor{c})
	case proto.EventSource:
		c.switchSubscribeEventType(event, sourceEntityProcessor{c})
	case proto.EventSinkSinkInput:
		c.switchSubscribeEventType(event, sinkInputEntityProcessor{c})
	case proto.EventSinkSourceOutput:
		c.switchSubscribeEventType(event, sourceOutputEntityProcessor{c})
	default:
		c.log.Debug("unknown event", "event", event)
	}
}

func (c *Connection) switchSubscribeEventType(event *proto.SubscribeEvent, processor entityProcessor) {
	switch event.Event & proto.EventTypeMask {
	case proto.EventNew, proto.EventChange:
		processor.updateId(event.Index)
	case proto.EventRemove:
		processor.removeId(event.Index)
	default:
	}
}

func (c *Connection) updateEntity(name string, cmd commandProcessor, components ...any) {
	ent, exists := c.entities.GetNamed(name)
	if exists {
		c.log.Debug("updating entity", "name", name, "components", components)
		ent.SetComponents(components...)
	} else {
		c.log.Debug("creating entity", "name", name, "components", components)
		ent = entity.NewEntity(entity.WithName(name))
		ent.SetComponents(processDriverComponent(c))
		ent.SetComponents(components...)
		if cmd != nil {
			cmd.registerCommands(ent)
		}
		err := c.entities.Add(ent)
		if err != nil {
			c.log.Error("error while adding entity", "name", name, "error", err)
			return
		}
	}
}

func (c *Connection) removeEntity(name string) {
	ent, exists := c.entities.GetNamed(name)
	if exists {
		c.log.Debug("removing entity", "name", name)
		err := c.entities.Remove(ent)
		if err != nil {
			c.log.Error("error while removing entity", "name", name, "error", err)
		}
	}
}
