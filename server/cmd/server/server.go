package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/burbokop/architecture-lab3/server/channels"
)

type HttpPortNumber int

// ChatApiServer configures necessary handlers and starts listening on a configured port.
type ChatApiServer struct {
	Port HttpPortNumber

	ListVmsHandler        channels.HttpHandlerFunc
	DiscConnectionHandler channels.HttpHandlerFunc

	server *http.Server
}

// Start will set all handlers and start listening.
// If this methods succeeds, it does not return until server is shut down.
// Returned error will never be nil.
func (s *ChatApiServer) Start() error {
	if s.ListVmsHandler == nil {
		return fmt.Errorf("HTTP ListVmsHandler is not defined - cannot start")
	}
	if s.DiscConnectionHandler == nil {
		return fmt.Errorf("HTTP DiscConnectionHandler is not defined - cannot start")
	}
	if s.Port == 0 {
		return fmt.Errorf("port is not defined")
	}

	handler := new(http.ServeMux)
	handler.HandleFunc("/vm_list", s.ListVmsHandler)
	handler.HandleFunc("/connect_disc", s.DiscConnectionHandler)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: handler,
	}

	return s.server.ListenAndServe()
}

// Stops will shut down previously started HTTP server.
func (s *ChatApiServer) Stop() error {
	if s.server == nil {
		return fmt.Errorf("server was not started")
	}
	return s.server.Shutdown(context.Background())
}
