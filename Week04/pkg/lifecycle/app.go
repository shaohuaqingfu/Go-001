package lifecycle

import (
	"Week04/pkg/sync/errgroup"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	opts   Options
	hooks  []Hook
	cancel func()
}

type Option func(o *Options)

type Options struct {
	StartTimeout time.Duration
	StopTimeout  time.Duration
	Signals      []os.Signal
	SignalFn     func(*App, os.Signal)
}

type Hook struct {
	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

func NewApp(opts ...Option) *App {

	options := Options{
		StartTimeout: time.Second * 30,
		StopTimeout:  time.Second * 30,
		Signals: []os.Signal{
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGINT,
		},
		SignalFn: func(app *App, signal os.Signal) {
			switch signal {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				app.Stop()
			default:
			}
		},
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &App{
		opts: options,
	}
}

func (app *App) Append(hook Hook) {
	app.hooks = append(app.hooks, hook)
}

func (app *App) Run(ctx context.Context) error {
	// 应用取消
	ctx, app.cancel = context.WithCancel(ctx)
	// 协程组
	group, groupCtx := errgroup.WithCancel(ctx)

	// 对所有的回调初始化
	for _, hook := range app.hooks {
		// 设置局部变量, 防止goroutine启动后, 调用for range中的hook的地址被改变
		hook := hook
		// 启动停止协程
		if hook.OnStop != nil {
			group.Go(func(ctx context.Context) error {
				// 监听group done
				<-groupCtx.Done()
				// 停止回调超时
				stopCtx, cancel := context.WithTimeout(ctx, app.opts.StopTimeout)
				defer cancel()
				return hook.OnStop(stopCtx)
			})
		}
		// 启动初始化协程
		if hook.OnStart != nil {
			group.Go(func(ctx context.Context) error {
				// 初始化回调超时
				startCtx, cancel := context.WithTimeout(ctx, app.opts.StartTimeout)
				defer cancel()
				return hook.OnStart(startCtx)
			})
		}
	}

	if len(app.opts.Signals) == 0 {
		return group.Wait()
	}
	signalChan := make(chan os.Signal, len(app.opts.Signals))
	// 监听linux signal
	signal.Notify(signalChan, app.opts.Signals...)

	group.Go(func(_ context.Context) error {
		for {
			select {
			// 监听协程组的done，如果有一个goroutine出现error，则返回error
			case <-groupCtx.Done():
				return groupCtx.Err()
			// 监听linux signal，有指定信号会触发应用结束
			case sig := <-signalChan:
				if app.opts.SignalFn != nil {
					app.opts.SignalFn(app, sig)
				}
			}
		}
	})
	return group.Wait()
}

func (app *App) Stop() {
	if app.cancel != nil {
		app.cancel()
	}
}
