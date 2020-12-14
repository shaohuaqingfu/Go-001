package tracker

import (
	"Week02/src/sync"
	"context"
	"log"
	"time"
)

type Tracker struct {
	grPool *sync.GRPool
}

func NewTracker() *Tracker {
	return &Tracker{
		grPool: sync.NewGRPool(100),
	}
}

func (t *Tracker) Event(data string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := t.grPool.Go(ctx, func() {
		// DoSomething
	})
	if err != nil {
		log.Printf("error = %s", err.Error())
		return
	}
}

func (t *Tracker) Shutdown() {
	t.grPool.Wait()
}
