package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/snakeice/gogress"
)

const (
	TOTAL = 100
	BARS  = 5
)

func main() {
	pool := gogress.NewPool()

	newBar := func() *gogress.Progress {
		bar := pool.NewBar(TOTAL)
		bar.ShowElapsedTime = true
		bar.ShowTimeLeft = true
		bar.ShowCounters = true
		return bar
	}

	var wg sync.WaitGroup
	wg.Add(BARS)
	pool.Start()

	for i := 0; i < BARS; i++ {
		bar := newBar()

		name := fmt.Sprintf("Task %d - %%d:", i)
		bar.Prefix(name)
		var li = i
		go func() {
			other := newBar()
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			max := 100 * time.Millisecond
			other.Prefix(fmt.Sprintf("Other %d", li))
			for bar.GetCurrent() < TOTAL {
				bar.Prefix(fmt.Sprintf(name, rng.Intn(3)))
				time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
				bar.Add(rng.Intn(3))
			}
			wg.Add(1)
			for other.GetCurrent() < TOTAL {
				time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
				other.Add(rng.Intn(3))
			}
			wg.Done()

		}()
	}

	wg.Wait()
	pool.FinishAll()
}
