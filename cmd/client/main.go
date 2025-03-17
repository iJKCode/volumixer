package main

import (
	"connectrpc.com/connect"
	"context"
	corev1 "ijkcode.tech/volumixer/proto/core/v1"
	"ijkcode.tech/volumixer/proto/core/v1/corev1connect"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	client := corev1connect.NewEntityServiceClient(
		http.DefaultClient,
		"http://localhost:5000",
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	res, err := client.EventStream(ctx, connect.NewRequest(&corev1.EventStreamRequest{
		SimulateState: false,
	}))
	if err != nil {
		slog.Error("failed to connect event stream", "error", err.Error())
	} else {
		for res.Receive() {
			slog.Info("got event", "event", res.Msg())
		}
	}
}
