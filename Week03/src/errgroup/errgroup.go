package errgroup

import (
	"context"
	"fmt"
	"sync"
)

type Group struct {
	wg       sync.WaitGroup
	rest     chan func(ctx context.Context, cancel context.CancelFunc) error
	fs       []func(ctx context.Context, cancel context.CancelFunc) error
	err      error
	errOnce  sync.Once
	poolOnce sync.Once
	ctx      context.Context
	cancel   func()
}

func WithContext(ctx context.Context) *Group {
	return newGroup(ctx, nil)
}

func WithCancel(ctx context.Context) *Group {
	ctx, cancelFunc := context.WithCancel(ctx)
	return newGroup(ctx, cancelFunc)
}

func newGroup(ctx context.Context, cancel func()) *Group {
	return &Group{
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		poolOnce: sync.Once{},
		cancel:   cancel,
	}
}

func (g *Group) WithPool(n int) *Group {
	// n个任务池
	g.rest = make(chan func(ctx context.Context, cancel context.CancelFunc) error, n)
	go func() {
		g.poolOnce.Do(func() {
			// 初始化协程池
			// rest空时会阻塞，每有一个任务就启动一个协程执行
			for f := range g.rest {
				go g.do(f)
			}
		})
	}()
	return g
}

func (g *Group) do(f func(ctx context.Context, cancel context.CancelFunc) error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("task error = %s\n", r)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f(g.ctx, g.cancel)
}

func (g *Group) Go(f func(ctx context.Context, cancel context.CancelFunc) error) {
	g.wg.Add(1)
	// 是否以协程池执行任务
	if g.rest != nil {
		select {
		// 如果rest没有满，使用协程池执行任务
		case g.rest <- f:
		// 如果rest满了，进入任务缓存队列
		default:
			g.fs = append(g.fs, f)
		}
		return
	}
	// 开启单个协程执行单个任务
	go g.do(f)
}

func (g *Group) Wait() error {
	defer func() {
		if g.rest != nil {
			close(g.rest)
		}
		if g.cancel != nil {
			g.cancel()
		}
	}()
	if g.rest != nil {
		for _, f := range g.fs {
			g.rest <- f
		}
	}
	g.wg.Wait()
	return g.err
}
