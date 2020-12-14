package grpctransport

import (
	"context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	opts serverOptions
}

func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	srv := grpc.NewServer()
	return &Server{
		srv,
		options,
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}
	return s.Serve(listener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	return nil
}
