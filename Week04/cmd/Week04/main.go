package main

import (
	"Week04/api/order"
	"Week04/api/product"
	"Week04/internal/Week04/data"
	"Week04/internal/Week04/service"
	"Week04/internal/Week04/wires"
	"Week04/pkg/lifecycle"
	grpctransport "Week04/pkg/transport/grpc"
	httptransport "Week04/pkg/transport/http"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	mux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	err := order.RegisterOrderServiceHandlerFromEndpoint(context.Background(), mux, "127.0.0.1:8081", options)
	if err != nil {
		return
	}

	grpcsrv := grpctransport.NewServer(grpctransport.WithAddress("127.0.0.1:8081"))
	httpsrv := httptransport.NewServer(
		httptransport.WithAddress("127.0.0.1:8080"),
		httptransport.WithMux(mux),
	)

	app := lifecycle.NewApp()
	app.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			db, err := gorm.Open("mysql", "root:13643566666@tcp(127.0.0.1:3306)/wblog?parseTime=true")
			if err != nil {
				panic(err)
			}
			productService := wires.InitProductService(db)
			orderService := &service.OrderService{
				Dao: &data.OrderData{
					DB: db,
				},
			}
			product.RegisterProductServiceServer(grpcsrv.Server, productService)
			order.RegisterOrderServiceServer(grpcsrv.Server, orderService)
			return grpcsrv.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return grpcsrv.Stop(ctx)
		},
	})
	app.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			time.Sleep(time.Second)
			return httpsrv.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return httpsrv.Stop(ctx)
		},
	})

	if err := app.Run(context.Background()); err != nil {
		log.Printf("app failed: error[%v]\n", err)
		return
	}
}
