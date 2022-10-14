package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/bingoohuang/pb"
)

func main() {
	pool := &pb.Pool{}
	first := pb.Full.New(1000).Set("prefix", "First ").SetMaxWidth(100)
	second := pb.Full.New(1000).Set("prefix", "Second").SetMaxWidth(100)
	third := pb.Full.New(1000).Set("prefix", "Third ").SetMaxWidth(100)
	pool.Start(first, second, third)
	wg := new(sync.WaitGroup)
	for _, bar := range []*pb.ProgressBar{first, second, third} {
		wg.Add(1)
		go func(cb *pb.ProgressBar) {
			for n := 0; n < 1000; n++ {
				cb.Increment()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			}
			cb.Finish()
			wg.Done()
		}(bar)
	}
	wg.Wait()
}
