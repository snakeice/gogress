package main

import (
	"time"

	"github.com/snakeice/gogress"
	"github.com/snakeice/gogress/format"
)

const TOTAL = 100

func main() {
	barFormat := format.ProgressFormat{
		BoxStart:  "|",
		BoxEnd:    "|",
		Empty:     "_",
		Current:   ">",
		Completed: "-",
		Spin:      []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
	}

	template := `{{prefix . 2 | green }} {{spin . 2 | yellow }} {{bar . -1 | rndcolor }} {{percent . 1 | cyan }} {{counter . 1 | red }} {{speed . 1 | blue}}`

	bar := gogress.New(TOTAL)

	if err := bar.SetTemplate(template); err != nil {
		panic(err)
	}

	bar.Format = barFormat
	bar.Prefix("Downloading life")

	bar.Start()
	for i := 1; i <= TOTAL; i++ {
		bar.Inc()
		time.Sleep(time.Second / 60)
	}
	bar.FinishPrint("All Solved")
}
