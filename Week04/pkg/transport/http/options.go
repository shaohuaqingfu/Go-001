package httptransport

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// ServerOption is gRPC server option.
type ServerOption func(o *serverOptions)

type serverOptions struct {
	Address string
	Mux     *runtime.ServeMux
}

// WithAddress is bind address option.
func WithAddress(a string) ServerOption {
	return func(o *serverOptions) {
		o.Address = a
	}
}

func WithMux(mux *runtime.ServeMux) ServerOption {
	return func(o *serverOptions) {
		o.Mux = mux
	}
}
