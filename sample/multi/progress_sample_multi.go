package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/snakeice/gogress"
)

const (
	TOTAL = 100
	BARS  = 3
)

func processBar(bar *gogress.Progress) {
	for bar.GetMax() > bar.GetCurrent() {
		bar.Inc()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}

	bar.Finish()
}

func main() {
	pool := gogress.NewPool()
	pool.Start()
	for i := 0; i < BARS; i++ {
		bar := pool.NewBar(TOTAL)
		bar.Prefix("Processing " + strconv.Itoa(i))
		go processBar(bar)
	}

	bar := pool.NewBar(TOTAL)
	bar.Add(50).Prefix("Other processing")
	processBar(bar)
	pool.RemoveBar(bar)

	pool.Wait()

}
