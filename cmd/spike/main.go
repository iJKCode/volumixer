package main

import (
	"context"
	server2 "ijkcode.tech/volumixer/pkg/core/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ijkcode.tech/volumixer/pkg/server"
	corepb "ijkcode.tech/volumixer/proto/core"
	"ijkcode.tech/volumixer/proto/core/corepbconnect"
)

func main() {
	var err error

	srv, err := server2.NewServer(":8080")
	if err != nil {
		log.Fatalf("error: creating server: %v", err)
	}

	srv.Mux().HandleFunc("/query", queryHandler)

	srv.Mux().Handle(corepbconnect.NewCoreServiceHandler(&server.ProtoServiceHandler{}))
	srv.RegisterConnectServiceName(corepbconnect.CoreServiceName)
	err = corepb.RegisterCoreServiceHandler(srv.GatewayArgs())
	if err != nil {
		log.Printf("error: registering core gateway: %v", err)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		log.Printf("info: shutting down server")
		err := srv.Stop()
		if err != nil {
			log.Printf("error: stopping server: %v", err)
		}
	}()

	log.Printf("info: starting server")
	err = srv.Serve()
	if err != nil {
		log.Printf("error: running server: %v", err)
	}
}

func queryHandler(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("hello world"))
	if err != nil {
		log.Printf("error: writing response: %v", err)
	}
}
