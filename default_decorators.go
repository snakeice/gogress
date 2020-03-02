package gogress

import (
	"math"
	"strings"

	"github.com/snakeice/gogress/format"
)

func bar(frame *FrameContext, gridCols, cols int) string {
	bar := " "
	barSize := cols
	if barSize <= 0 {
		bar = ""
	} else if frame.max > 0 {
		decPercent := float64(frame.current) / float64(frame.max)
		doneSize := int(math.Ceil(decPercent * float64(barSize)))
		emptySize := barSize - doneSize

		if emptySize < 0 {
			emptySize = 0
		}

		if doneSize > barSize {
			doneSize = barSize
		}
		bar += frame.format().BoxStart

		cursorLen := format.EscapeAwareRuneCountInString(frame.format().Completed)
		if emptySize <= 0 {
			bar += strings.Repeat(frame.format().Completed, doneSize/cursorLen)
		} else if doneSize > 0 {
			cursorEndLen := format.EscapeAwareRuneCountInString(frame.format().Completed)
			cursorRepetitions := (doneSize - cursorEndLen) / cursorLen
			bar += strings.Repeat(frame.format().Completed, cursorRepetitions)
			bar += frame.format().Current
		}

		emptyLen := format.EscapeAwareRuneCountInString(frame.format().Empty)
		bar += strings.Repeat(frame.format().Empty, emptySize/emptyLen)
		bar += frame.format().BoxEnd
	} else {
		pos := barSize - int(frame.current)%int(barSize)
		bar += frame.format().BoxStart
		if pos-1 > 0 {
			bar += strings.Repeat(frame.format().Empty, pos-1)
		}
		bar += frame.format().Current
		if barSize-pos-1 > 0 {
			bar += strings.Repeat(frame.format().Empty, barSize-pos-1)
		}
		bar += frame.format().BoxEnd
	}
	return bar
}

func prefix(frame *FrameContext, gridCols, cols int) string {
	msg := frame.bar.messagePrefix

	if len(msg) > gridCols {
		msg = msg[:gridCols-3] + "..."

	} else if len(msg) < gridCols {
		msg += strings.Repeat(" ", gridCols-len(msg))
	}

	return msg
}
