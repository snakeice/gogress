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
	BARS  = 6
)

func main() {
	pool := gogress.NewPool()
	pool.RefreshRate = time.Second / 30

	newBar := func() *gogress.Progress {
		bar := pool.NewBar(TOTAL)
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
			other.SetMax(TOTAL)
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			max := 100 * time.Millisecond
			other.Prefix(fmt.Sprintf("Other %d", li))
			wg.Add(1)
			go func() {
				for other.GetCurrent() < TOTAL {
					time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
					other.Add(rng.Intn(3))
				}
				pool.RemoveBar(other)
				wg.Done()
			}()

			for bar.GetCurrent() < TOTAL {
				bar.Prefix(fmt.Sprintf(name, rng.Intn(3)))
				time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
				bar.Add(rng.Intn(3))
			}
		}()
	}

	wg.Wait()
	pool.FinishAll()
}
