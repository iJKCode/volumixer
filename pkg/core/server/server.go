package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	endpoint     string
	httpMux      *http.ServeMux
	httpServer   *http.Server
	reflectNames *ReflectNames
}

func NewServer(endpoint string) (*Server, error) {

	// create http mux
	mux := http.NewServeMux()

	// construct server
	srv := &Server{
		endpoint: endpoint,
		httpMux:  mux,
		httpServer: &http.Server{
			Handler: h2c.NewHandler(mux, &http2.Server{}),
			Addr:    endpoint,
		},
		reflectNames: &ReflectNames{},
	}

	// register connect reflector
	reflector := grpcreflect.NewReflector(srv.reflectNames)
	path1, handler1 := grpcreflect.NewHandlerV1(reflector)
	srv.httpMux.Handle(path1, handler1)
	path2, handler2 := grpcreflect.NewHandlerV1Alpha(reflector)
	srv.httpMux.Handle(path2, handler2)

	return srv, nil
}

func (s *Server) Endpoint() string {
	return s.endpoint
}

func (s *Server) ServeMux() *http.ServeMux {
	return s.httpMux
}

func (s *Server) ReflectNames() *ReflectNames {
	return s.reflectNames
}

func (s *Server) Serve() error {
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Stop() error {
	deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(deadline)
}
