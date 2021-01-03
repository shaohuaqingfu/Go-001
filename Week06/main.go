package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	// 请求状态数量
	Success int64
	Fail    int64
	// 上次请求时更新的时间
	last time.Time
}

type SlideWindow struct {
	buckets []*Bucket
	width   time.Duration
	// bucket容量
	size int
	// 上次请求的记录bucket的下标
	idx int
	mu  sync.Mutex
}

func OpenSlideWindow(w time.Duration, s int) *SlideWindow {
	var buckets []*Bucket
	now := time.Now()
	for i := 0; i < s; i++ {
		buckets = append(buckets, &Bucket{
			last:    now,
			Success: 0,
			Fail:    0,
		})
	}

	return &SlideWindow{
		buckets: buckets,
		size:    s,
		width:   w,
		idx:     0,
	}
}

func (sw *SlideWindow) IncrSuccess() {
	bucket := sw.getWindow()
	atomic.AddInt64(&bucket.Success, 1)
}

func (sw *SlideWindow) getWindow() *Bucket {
	t := time.Now()
	// 获取上个请求的bucket
	sw.mu.Lock()
	defer sw.mu.Unlock()
	bucket := sw.buckets[sw.idx]

	if t.Add(-sw.width).After(bucket.last) {
		// 如果时间间隔超过了指定的窗口宽度，则记录到下一个窗口中
		sw.idx = (sw.idx + 1) % sw.size
		sw.buckets[sw.idx] = &Bucket{
			Success: 0,
			Fail:    0,
			last:    t,
		}
		bucket = sw.buckets[sw.idx]
	}
	return bucket
}

func (sw *SlideWindow) GetSuccess() int64 {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	var res int64
	for _, bucket := range sw.buckets {
		res += bucket.Success
	}
	return res
}

func (sw *SlideWindow) GetFail() int64 {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	var res int64
	for _, bucket := range sw.buckets {
		res += bucket.Fail
	}
	return res
}

func (sw *SlideWindow) IncrFail() {
	bucket := sw.getWindow()
	atomic.AddInt64(&bucket.Fail, 1)
}

func main() {
	group := sync.WaitGroup{}
	sw := OpenSlideWindow(1*time.Second, 10)
	var num int64 = 0
	for i := 0; i < 10; i++ {
		group.Add(1)
		go func() {
			for j := 0; j < 10; j++ {
				n := rand.Intn(3)
				if n == 1 {
					sw.IncrFail()
				} else {
					sw.IncrSuccess()
				}
				printBucket(&num, sw)
				time.Sleep(time.Duration(n) * time.Second)
			}
			group.Done()
		}()
	}
	group.Wait()
	time.Sleep(time.Second * 2)
	fmt.Println("success = [", sw.GetSuccess(), "]    failed = [", sw.GetFail(), "]")
}

func printBucket(num *int64, sw *SlideWindow) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	fmt.Printf("[num:%d] ", *num)
	*num = *num + 1
	for _, bucket := range sw.buckets {
		fmt.Printf("[s:%d,f:%d,t:%d] ", bucket.Success, bucket.Fail, bucket.last.Second())
	}
	fmt.Println()
}
