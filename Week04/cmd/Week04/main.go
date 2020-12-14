package main

import (
	"Week04/pkg/lifecycle"
	http2 "Week04/pkg/server/http"
	"context"
	"log"
)

func main() {

	srv := http2.NewServer(func(options *http2.Options) {
		options.Addr = "127.0.0.1:80"
	})
	app := lifecycle.NewApp()
	app.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			return srv.Start()
		},
		OnStop: func(ctx context.Context) error {
			return srv.Stop()
		},
	})

	if err := app.Run(context.Background()); err != nil {
		log.Printf("app failed: error[%v]\n", err)
		return
	}
}
