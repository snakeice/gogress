# gogress


Simple terminal progress bar with Go.

Features:

- Customizable progress bar
- Multi-progress bar
- Multi-threaded multi-progress bar
- IO wrapper

Based on [cheggaaa/pb](https://github.com/cheggaaa/pb)

# Sample usage - Simple progress bar

```go
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
```

[![asciicast](https://asciinema.org/a/wqZKNwxiQErdrVG4fDlFdQqLZ.svg)](https://asciinema.org/a/wqZKNwxiQErdrVG4fDlFdQqLZ)

# Sample decorator

```go
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
```

[![asciicast](https://asciinema.org/a/BjDqCQYQyQGZJeD3D61absJfq.svg)](https://asciinema.org/a/BjDqCQYQyQGZJeD3D61absJfq)

# Sample usage - Multi-progress bar

```go
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
```

[![asciicast](https://asciinema.org/a/REaVd9hpZ6iEZN6LYz0uqEG4B.svg)](https://asciinema.org/a/REaVd9hpZ6iEZN6LYz0uqEG4B)
