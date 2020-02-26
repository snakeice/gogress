package gogress

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/snakeice/gogress/format"
)

type printerContext struct {
	bar          *Progress
	width        int
	max          int64
	current      int64
	countersBox  string
	timeSpentBox string
	speedBox     string
	barBox       string
	percentBox   string
	timeLeftBox  string
	end          string
	out          string
	mu           sync.Mutex
}

func newPrintContex(bar *Progress, max, current int64) *printerContext {
	return &printerContext{
		bar:     bar,
		width:   bar.GetWidth(),
		max:     max,
		current: current,
	}
}

func (pc *printerContext) Update(max, current int64) *printerContext {
	pc.width = pc.bar.GetWidth()
	pc.max = max
	pc.current = current
	return pc

}

func (pc *printerContext) String() string {
	return pc.out + pc.end
}

func (pc *printerContext) Feed() string {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.feedCounters()
	pc.feedBar()
	pc.feedEnd()

	pc.bar.last = pc.String()
	return pc.String()
}

func (pc *printerContext) feedCounters() {
	percent := float64(pc.current) / (float64(pc.max) / 100.0)

	pc.percentBox = fmt.Sprintf(" %.02f%%", percent)

	currentFromStart := pc.current - pc.bar.startValue
	fromStart := time.Since(pc.bar.startTime)
	lastChangeTime := pc.bar.changeTime
	fromChange := lastChangeTime.Sub(pc.bar.startTime)

	var left time.Duration
	select {
	case <-pc.bar.finish:
		left = (fromStart / time.Second) * time.Second
		pc.timeLeftBox = fmt.Sprintf(" %s", left.String())
	default:
		if pc.bar.ShowElapsedTime {
			pc.timeSpentBox = fmt.Sprintf(" %s", (fromStart/time.Second)*time.Second)
		}

		if pc.bar.ShowTimeLeft && currentFromStart > 0 {
			perEntry := fromChange / time.Duration(currentFromStart)
			if pc.max > 0 {
				left = time.Duration(pc.max-pc.current) * perEntry
				left -= time.Since(lastChangeTime)
				left = (left / time.Second) * time.Second
			}
			if left > 0 {
				timeLeft := format.Format(int64(left)).To(format.U_DURATION).String()
				pc.timeLeftBox = fmt.Sprintf(" %s", timeLeft)
			}
		}
	}

	if len(pc.timeLeftBox) < pc.bar.TimeBoxWidth {
		pc.timeLeftBox = fmt.Sprintf("%s%s", strings.Repeat(" ", pc.bar.TimeBoxWidth-len(pc.timeLeftBox)), pc.timeLeftBox)
	}

	if pc.bar.ShowSpeed && currentFromStart > 0 {
		fromStart := time.Since(pc.bar.startTime)
		speed := float64(currentFromStart) / (float64(fromStart) / float64(time.Second))
		pc.speedBox = " " + format.Format(int64(speed)).To(pc.bar.Units).Width(pc.bar.UnitsWidth).PerSec().String()
	}

	if pc.bar.ShowCounters {
		current := format.Format(pc.current).To(pc.bar.Units).Width(pc.bar.UnitsWidth)
		if pc.max > 0 {
			totalS := format.Format(pc.max).To(pc.bar.Units).Width(pc.bar.UnitsWidth)
			pc.countersBox = fmt.Sprintf(" %s/%s", current, totalS)
		} else {
			pc.countersBox = fmt.Sprintf(" %s/?", current)
		}
	}
}

func (pc *printerContext) getBarSize() int {
	barWidth := format.EscapeAwareRuneCountInString(
		pc.bar.Format.BoxStart +
			pc.bar.Format.BoxEnd +
			pc.percentBox +
			pc.timeSpentBox +
			pc.timeLeftBox +
			pc.speedBox +
			pc.countersBox +
			pc.bar.messagePrefix)

	return pc.width - barWidth
}

func (pc *printerContext) format() *format.ProgressFormat {
	return pc.bar.Format
}

func (pc *printerContext) feedBar() {
	bar := " "
	barSize := pc.getBarSize() - 1
	if barSize <= 0 {
		bar = ""
	} else if pc.max > 0 {
		decPercent := float64(pc.current) / float64(pc.max)
		doneSize := int(math.Ceil(decPercent * float64(barSize)))
		emptySize := barSize - doneSize

		if emptySize < 0 {
			emptySize = 0
		}

		if doneSize > barSize {
			doneSize = barSize
		}
		bar += pc.format().BoxStart

		cursorLen := format.EscapeAwareRuneCountInString(pc.format().Completed)
		if emptySize <= 0 {
			bar += strings.Repeat(pc.format().Completed, doneSize/cursorLen)
		} else if doneSize > 0 {
			cursorEndLen := format.EscapeAwareRuneCountInString(pc.format().Completed)
			cursorRepetitions := (doneSize - cursorEndLen) / cursorLen
			bar += strings.Repeat(pc.format().Completed, cursorRepetitions)
			bar += pc.format().Current
		}

		emptyLen := format.EscapeAwareRuneCountInString(pc.format().Empty)
		bar += strings.Repeat(pc.format().Empty, emptySize/emptyLen)
		bar += pc.format().BoxEnd
	} else {
		pos := barSize - int(pc.current)%int(barSize)
		bar += pc.format().BoxStart
		if pos-1 > 0 {
			bar += strings.Repeat(pc.format().Empty, pos-1)
		}
		bar += pc.format().Current
		if barSize-pos-1 > 0 {
			bar += strings.Repeat(pc.format().Empty, barSize-pos-1)
		}
		bar += pc.format().BoxEnd
	}
	pc.barBox = bar
}

func (pc *printerContext) feedEnd() {
	pc.out = pc.bar.messagePrefix +
		pc.countersBox +
		pc.barBox +
		pc.percentBox +
		pc.speedBox +
		pc.timeSpentBox +
		pc.timeLeftBox

	if cl := format.EscapeAwareRuneCountInString(pc.out); cl < pc.width {
		pc.end = strings.Repeat(" ", pc.width-cl)
	}

}
