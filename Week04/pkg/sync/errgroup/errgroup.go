package errgroup

import (
	"context"
	"fmt"
	"sync"
)

type Group struct {
	wg       sync.WaitGroup
	rest     chan func(ctx context.Context) error
	fs       []func(ctx context.Context) error
	err      error
	errOnce  sync.Once
	poolOnce sync.Once
	ctx      context.Context
	cancel   func()
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	return newGroup(ctx, nil)
}

func WithCancel(ctx context.Context) (*Group, context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	return newGroup(ctx, cancelFunc)
}

func newGroup(ctx context.Context, cancel func()) (*Group, context.Context) {
	return &Group{
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		poolOnce: sync.Once{},
		cancel:   cancel,
	}, ctx
}

func (g *Group) WithPool(n int) *Group {
	// n个任务池
	g.rest = make(chan func(ctx context.Context) error, n)
	go func() {
		g.poolOnce.Do(func() {
			// 初始化协程池
			// rest空时会阻塞，每有一个任务就启动一个协程执行
			for f := range g.rest {
				go g.do(g.ctx, f)
			}
		})
	}()
	return g
}

func (g *Group) do(ctx context.Context, f func(ctx context.Context) error) {
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
	err = f(ctx)
}

func (g *Group) Go(f func(ctx context.Context) error) {
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
	go g.do(g.ctx, f)
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
