package sync

import (
	"context"
	"errors"
	"log"
	"sync"
)

type GRPool struct {
	group    *sync.WaitGroup
	n        int
	resource chan bool
}

func NewGRPool(n int) *GRPool {
	return &GRPool{
		group:    &sync.WaitGroup{},
		n:        n,
		resource: make(chan bool, n),
	}
}

func (gp *GRPool) Go(ctx context.Context, fn func()) error {
	if gp.resource == nil {
		return errors.New("resource has closed")
	}
	// 控制一组goroutine的执行结束
	gp.group.Add(1)
	// 控制剩余资源数
	gp.resource <- true
	// 控制goroutine结束
	done := make(chan bool, 1)
	defer close(done)

	go func() {
		defer func() {
			<-gp.resource
			gp.group.Done()
			done <- true
		}()
		defer func() {
			if err := recover(); err != nil {
				log.Fatalf("[系统错误] error = %s", err)
			}
		}()

		fn()
	}()

	select {
	// 正常结束
	case <-done:
		return nil
	// 超时
	case <-ctx.Done():
		return errors.New("超时")
	}
}

func (gp *GRPool) Wait() {
	gp.group.Wait()
}

func (gp *GRPool) Close() {
	close(gp.resource)
}
