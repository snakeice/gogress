package main

import (
	"time"

	"github.com/snakeice/gogress"
)

const TOTAL = 500

func main() {
	bar := gogress.New(TOTAL)
	bar.ShowElapsedTime = true
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()
	bar.Prefix("Downloading life")
	for i := 1; i <= TOTAL; i++ {
		bar.Inc()
		time.Sleep(time.Second / 120)
	}
	bar.FinishPrint("All Solved")
}
