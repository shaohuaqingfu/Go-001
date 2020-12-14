package main

import (
	"Week04/api/order"
	"Week04/internal/Week04/service"
	"Week04/pkg/lifecycle"
	grpctransport "Week04/pkg/transport/grpc"
	"context"
	"log"
)

func main() {

	server := grpctransport.NewServer(grpctransport.WithAddress("127.0.0.1:8081"))

	order.RegisterOrderServiceServer(server.Server, &service.OrderService{})

	app := lifecycle.NewApp()
	app.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(context.Background())
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(context.Background())
		},
	})

	if err := app.Run(context.Background()); err != nil {
		log.Printf("app failed: error[%v]\n", err)
		return
	}
}
