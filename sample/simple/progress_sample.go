package main

import (
	"fmt"
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
	for i := 1; i <= TOTAL; i++ {
		bar.Inc()
		bar.Prefix(fmt.Sprintf("Solving problem %d", i))
		time.Sleep(time.Second / 120)
	}
	bar.FinishPrint("All Solved")
}
