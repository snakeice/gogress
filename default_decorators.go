package gogress

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/snakeice/gogress/format"
)

func bar(frame *FrameContext, cols int) string {
	var bar string
	cols -= len(frame.Bar.Format.BoxStart) + len(frame.Bar.Format.BoxEnd)
	if cols <= 0 {
		bar = ""
	} else if frame.Max > 0 {
		decPercent := float64(frame.Current) / float64(frame.Max)
		doneSize := int(math.Ceil(decPercent * float64(cols)))
		emptySize := cols - doneSize

		if emptySize < 0 {
			emptySize = 0
		}

		if doneSize > cols {

			doneSize = cols
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
		pos := cols - int(frame.Current)%int(cols)
		bar += frame.format().BoxStart
		if pos-1 > 0 {
			bar += strings.Repeat(frame.format().Empty, pos-1)
		}
		bar += frame.format().Current
		if cols-pos-1 > 0 {
			bar += strings.Repeat(frame.format().Empty, cols-pos-1)
		}
		bar += frame.format().BoxEnd
	}
	return bar
}

func prefix(frame *FrameContext, cols int) string {
	msg := frame.Bar.messagePrefix

	if cols <= 3 {
		msg = strings.Repeat(".", cols)

	} else if len(msg) > cols {
		msg = msg[:cols-3] + "..."

	} else if len(msg) < cols {
		msg += strings.Repeat(" ", cols-len(msg))
	}

	return msg
}

func counter(frame *FrameContext, cols int) string {
	var counterBox string
	current := format.Format(frame.Current).To(frame.Bar.Units).Width(frame.Bar.UnitsWidth)
	if frame.Max > 0 {
		totalS := format.Format(frame.Max).To(frame.Bar.Units).Width(frame.Bar.UnitsWidth)
		counterBox = fmt.Sprintf("%s/%s", current, totalS)
	} else {
		counterBox = fmt.Sprintf("%s/?", current)
	}

	return counterBox
}

func timeSpent(frame *FrameContext, cols int) string {
	fromStart := time.Since(frame.Bar.startTime)

	var left time.Duration
	var timeSpentBox string
	select {
	case <-frame.Bar.finish:
		left = (fromStart / time.Second) * time.Second
		timeSpentBox = left.String()
	default:
		timeSpentBox = fmt.Sprintf("%s", (fromStart/time.Second)*time.Second)
	}

	return timeSpentBox
}

func speed(frame *FrameContext, cols int) string {
	fromStart := time.Since(frame.Bar.startTime)
	currentFromStart := frame.Current - frame.Bar.startValue

	var speedBox string
	speed := float64(currentFromStart) / (float64(fromStart) / float64(time.Second))
	speedBox = format.Format(int64(speed)).To(frame.Bar.Units).Width(frame.Bar.UnitsWidth).PerSec().String()

	return speedBox
}

func percent(frame *FrameContext, cols int) string {
	percent := float64(frame.Current) / (float64(frame.Max) / 100.0)
	return fmt.Sprintf("%.02f%%", percent)
}

func timeLeft(frame *FrameContext, cols int) string {
	currentFromStart := frame.Current - frame.Bar.startValue
	lastChangeTime := frame.Bar.changeTime
	fromChange := lastChangeTime.Sub(frame.Bar.startTime)

	var left time.Duration
	var timeLeftBox string
	select {
	case <-frame.Bar.finish:
	default:
		if currentFromStart > 0 {
			perEntry := fromChange / time.Duration(currentFromStart)
			if frame.Max > 0 {
				left = time.Duration(frame.Max-frame.Current) * perEntry
				left -= time.Since(lastChangeTime)
				left = (left / time.Second) * time.Second
			}
			if left > 0 {
				timeLeft := format.Format(int64(left)).To(format.U_DURATION).String()
				timeLeftBox = fmt.Sprintf("%s", timeLeft)
			}
		}
	}
	return timeLeftBox
}

func getSpinString(frame *FrameContext) string {
	return frame.Bar.Format.SpinString
}

func spin(frame *FrameContext, cols int) string {
	char := getSpinString(frame)[frame.frameNo%int64(len(getSpinString(frame)))]
	return string(char)
}

func frameNo(frame *FrameContext, cols int) string {
	return strconv.FormatInt(frame.FrameNo(), 10)
}
