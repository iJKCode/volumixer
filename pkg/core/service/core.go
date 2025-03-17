package service

import (
	"connectrpc.com/connect"
	"context"
	corev1 "ijkcode.tech/volumixer/proto/core/v1"
	"log/slog"
	"runtime"
)

type CoreServiceHandler struct {
	Log *slog.Logger
}

func (c CoreServiceHandler) Health(ctx context.Context, _ *connect.Request[corev1.HealthRequest]) (*connect.Response[corev1.HealthResponse], error) {
	c.Log.InfoContext(ctx, "got health request")
	return connect.NewResponse(&corev1.HealthResponse{}), nil
}

func (c CoreServiceHandler) ServerInfo(ctx context.Context, _ *connect.Request[corev1.ServerInfoRequest]) (*connect.Response[corev1.ServerInfoResponse], error) {
	c.Log.InfoContext(ctx, "got server info request")
	return connect.NewResponse(&corev1.ServerInfoResponse{
		ServerVersion:  "0.0.0-dev",
		ServerPlatform: runtime.GOOS,
	}), nil
}
