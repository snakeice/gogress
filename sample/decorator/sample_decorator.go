package main

import (
	"time"

	"github.com/snakeice/gogress"
	"github.com/snakeice/gogress/format"
)

const TOTAL = 500

func main() {
	barFormat := format.ProgressFormat{
		BoxStart:   "|",
		BoxEnd:     "|",
		Empty:      "_",
		Current:    ">",
		Completed:  "-",
		SpinString: "\\|/-",
	}

	template := `{{prefix . 2 | green }} {{spin . 1 | rndcolor }}  {{percent . 1 | cyan }} {{counter . 1 | red }} {{speed . 1 | blue}}`

	bar := gogress.New(TOTAL)

	if err := bar.SetTemplate(template); err != nil {
		panic(err)
	}

	bar.Format = barFormat
	bar.Prefix("Processing")

	bar.Start()
	bar.Prefix("Downloading life")
	for i := 1; i <= TOTAL; i++ {
		bar.Inc()
		time.Sleep(time.Second / 120)
	}
	bar.FinishPrint("All Solved")
}
