package sync

import (
	"Week09/errors"
	"context"
)

func Go(ctx context.Context, f func(ctx context.Context) error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errors.PrintError(err.(error))
			}
		}()
		err := f(ctx)
		if err != nil {
			panic(err)
		}
	}()
}
