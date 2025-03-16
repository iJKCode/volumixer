package main

import (
	"context"
	"fmt"
	"ijkcode.tech/volumixer/pkg/core/command"
	"ijkcode.tech/volumixer/pkg/core/entity"
	"ijkcode.tech/volumixer/pkg/core/event"
	"ijkcode.tech/volumixer/pkg/driver/pulseaudio"
	"ijkcode.tech/volumixer/pkg/widget"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	logLevel := slog.LevelInfo

	slog.SetLogLoggerLevel(logLevel)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	bus := event.NewBus()
	repo := entity.NewContext(bus)

	event.SubscribeAll(bus, event.Func(func(ctx context.Context, event any) {
		log.Info("entity event", "type", fmt.Sprintf("%T", event), "data", event)
	}))

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer stop()
		bus.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		defer stop()
		err := pulseaudio.NewConnection(log, repo, "").Run(ctx)
		if err != nil {
			log.Error("running pulseaudio connection", "error", err)
		}
	}()

	//time.Sleep(2 * time.Second)
	//scenarioLowerVolume(log, repo)
	//stop()

	<-ctx.Done()
	log.Info("shutting down...")
	wg.Wait()
	log.Info("shutdown complete")
}

func scenarioMuteUnmute(log *slog.Logger, entities *entity.Context) {
	for ent := range entities.All() {
		log.Info("muting", "entity", ent)
		err := command.DispatchEntity(ent, widget.VolumeMuteCommand{
			Mute: true,
		})
		if err != nil {
			log.Error("dispatch mute command", "entity", ent, "error", err)
		}

		time.Sleep(1 * time.Second)

		log.Info("unmuting", "entity", ent)
		err = command.DispatchEntity(ent, widget.VolumeMuteCommand{
			Mute: false,
		})
		if err != nil {
			log.Error("dispatch unmute command", "entity", ent, "error", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func scenarioLowerVolume(log *slog.Logger, entities *entity.Context) {
	for ent := range entities.All() {
		vol, ok := entity.GetComponent[widget.VolumeComponent](ent)
		if !ok {
			log.Error("volume component not found", "component", ent)
		}

		levelLower := vol.Level / 2.0
		levelRestore := vol.Level

		log.Info("lowering volume", "entity", ent, "level", levelLower)
		err := command.DispatchEntity(ent, widget.VolumeChangeCommand{
			Level: levelLower,
		})
		if err != nil {
			log.Error("dispatch volume lower command", "entity", ent, "error", err)
		}

		time.Sleep(1 * time.Second)

		log.Info("restoring", "entity", ent, "level", levelRestore)
		err = command.DispatchEntity(ent, widget.VolumeChangeCommand{
			Level: levelRestore,
		})
		if err != nil {
			log.Error("dispatch volume restore command", "entity", ent, "error", err)
		}

		time.Sleep(1 * time.Second)
	}
}
