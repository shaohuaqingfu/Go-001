package main

import (
	"Week03/src/errgroup"
	"Week03/src/sync"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// 1.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。
func main() {
	homework()
}

func homework() {

	exitChan := make(chan bool, 1)
	signalChan := make(chan os.Signal, 1)
	done := make(chan error, 2)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()
	group := errgroup.WithCancel(ctx).WithPool(10)

	server1 := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	server2 := &http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: nil,
	}
	sync.Go(ctx, func(ctx context.Context) error {
		for s := range signalChan {
			switch s {
			case syscall.SIGINT:
				fmt.Printf("receive ctrl^c")
				done <- errors.New("receive ctrl^c")
				break
			case syscall.SIGTERM:
				fmt.Printf("receive ctrl^\\")
				//cancel()
				break
			case syscall.SIGQUIT:
				fmt.Printf("receive ctrl^\\")
				//cancel()
				break
			}
		}
		return nil
	})
	sync.Go(ctx, func(ctx context.Context) error {
		for i := 0; i < cap(done); i++ {
			if i == 0 {
				<-done
				close(exitChan)
			} else {
				<-done
			}
		}
		return nil
	})
	group.Go(func(ctx context.Context, cancel context.CancelFunc) error {
		sync.Go(ctx, func(ctx context.Context) error {
			<-exitChan
			err := server1.Shutdown(ctx)
			if err == nil || errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		})
		err := server1.ListenAndServe()
		if err != nil {
			done <- err
		}
		fmt.Printf("server1 监听结束\n")
		return err
	})
	group.Go(func(ctx context.Context, cancel context.CancelFunc) error {
		sync.Go(ctx, func(ctx context.Context) error {
			<-exitChan
			err := server2.Shutdown(ctx)
			if err == nil || errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		})
		err := server2.ListenAndServe()
		if err != nil {
			done <- err
		}
		fmt.Printf("server2 监听结束\n")
		return err
	})
	err := group.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("[系统错误] error = %+v\n", errors.Wrap(err, "server error"))
	}
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("系统正常退出")
	}
}
