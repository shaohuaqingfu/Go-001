package http

import (
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

type Handler http.Handler

type Options struct {
	Addr    string
	Handler Handler
}

type Option func(*Options)

func NewServer(opts ...Option) *Server {
	options := Options{
		Addr:    "127.0.0.1",
		Handler: nil,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &Server{
		server: &http.Server{
			Addr:    options.Addr,
			Handler: options.Handler,
		},
	}
}

func (s *Server) Start() error {
	server := s.server
	return server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) Stop() error {
	return s.server.Close()
}
