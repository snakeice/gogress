package main

import (
	"time"

	"github.com/snakeice/gogress"
)

const TOTAL = 500

func main() {
	bar := gogress.New(TOTAL)
	bar.Start()
	bar.Prefix("Downloading life")
	for i := 1; i <= TOTAL; i++ {
		bar.Inc()
		time.Sleep(time.Second / 120)
	}
	bar.FinishPrint("All Solved")
}
