package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ijkcode/volumixer-api/gen/go/command/v1/commandv1connect"
	"github.com/ijkcode/volumixer-api/gen/go/entity/v1/entityv1connect"
	"github.com/ijkcode/volumixer/pkg/core/command"
	"github.com/ijkcode/volumixer/pkg/core/entity"
	"github.com/ijkcode/volumixer/pkg/core/event"
	"github.com/ijkcode/volumixer/pkg/core/server"
	"github.com/ijkcode/volumixer/pkg/core/service"
	"github.com/ijkcode/volumixer/pkg/driver/pulseaudio"
	"github.com/ijkcode/volumixer/pkg/widget"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	// register components
	_ "github.com/ijkcode/volumixer/pkg/widget"
)

func main() {
	logLevel := slog.LevelInfo

	slog.SetLogLoggerLevel(logLevel)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	bus := event.NewBus()
	repo := entity.NewContext(bus)

	srv, err := server.NewServer("0.0.0.0:5000")
	if err != nil {
		log.Error("error creating server", "error", err)
	}

	srv.ReflectNames().Add(entityv1connect.EntityServiceName)
	srv.ServeMux().Handle(entityv1connect.NewEntityServiceHandler(&service.EntityServiceHandler{
		Log:      log.With("handler", "EntityService"),
		Entities: repo,
	}))

	srv.ReflectNames().Add(commandv1connect.VolumeServiceName)
	srv.ServeMux().Handle(commandv1connect.NewVolumeServiceHandler(&widget.VolumeServiceHandler{
		Log:      log.With("handler", "VolumeService"),
		Entities: repo,
	}))

	proxyUrl, err := url.Parse("http://localhost:5173")
	if err != nil {
		log.Error("error parsing proxy url", "error", err)
	}
	proxyClient := httputil.NewSingleHostReverseProxy(proxyUrl)
	srv.ServeMux().HandleFunc("/app/", func(res http.ResponseWriter, req *http.Request) {
		proxyClient.ServeHTTP(res, req)
	})

	event.SubscribeAll(bus, event.Func(func(ctx context.Context, event any) {
		log.Info("entity event", "type", fmt.Sprintf("%T", event), "data", event)
	}))

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		defer stop()
		log.Info("starting event router")
		bus.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		defer stop()
		log.Info("starting pulseaudio driver")
		err := pulseaudio.NewConnection(log, repo, "").Run(ctx)
		if err != nil {
			log.Error("running pulseaudio connection", "error", err)
		}
	}()
	go func() {
		defer wg.Done()
		defer stop()
		log.Info("starting http server", "endpoint", srv.Endpoint())
		err := srv.Serve()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error running http server", "error", err)
		}
	}()

	//time.Sleep(2 * time.Second)
	//scenarioLowerVolume(log, repo)
	//stop()

	<-ctx.Done()
	log.Info("shutting down...")

	err = srv.Stop()
	if err != nil {
		log.Error("error stopping http server", "error", err)
	}

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
