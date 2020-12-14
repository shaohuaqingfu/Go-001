package sync

import (
	"context"
	"errors"
	"log"
)

func Go(ctx context.Context, f func(ctx context.Context) error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[系统错误] 123 error = %+v\n", err)
			}
		}()
		err := f(ctx)
		if err != nil {
			log.Printf("[系统错误] 456 error = %s\n", err.Error())
		}
		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				log.Printf("[系统错误] 789 error = %+v\n", ctx.Err())
			}
			return
		}
	}()
}
