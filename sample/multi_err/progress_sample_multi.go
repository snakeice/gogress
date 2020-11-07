package main

import (
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
		bar.Set(56)
		return bar
	}

	pool.Start()
	defer pool.FinishAll()

	for i := 0; i < BARS; i++ {
		_ = newBar()
	}

}
