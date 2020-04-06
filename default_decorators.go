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
	cols -= len(frame.Format().BoxStart) + len(frame.Format().BoxEnd)
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
		bar += frame.Format().BoxStart

		cursorLen := format.EscapeAwareRuneCountInString(frame.Format().Completed)
		if emptySize <= 0 {
			bar += strings.Repeat(frame.Format().Completed, doneSize/cursorLen)
		} else if doneSize > 0 {
			cursorEndLen := format.EscapeAwareRuneCountInString(frame.Format().Completed)
			cursorRepetitions := (doneSize - cursorEndLen) / cursorLen
			bar += strings.Repeat(frame.Format().Completed, cursorRepetitions)
			bar += frame.Format().Current
		}

		emptyLen := format.EscapeAwareRuneCountInString(frame.Format().Empty)
		bar += strings.Repeat(frame.Format().Empty, emptySize/emptyLen)
		bar += frame.Format().BoxEnd
	} else {
		pos := cols - int(frame.Current)%int(cols)
		bar += frame.Format().BoxStart
		if pos-1 > 0 {
			bar += strings.Repeat(frame.Format().Empty, pos-1)
		}
		bar += frame.Format().Current
		if cols-pos-1 > 0 {
			bar += strings.Repeat(frame.Format().Empty, cols-pos-1)
		}
		bar += frame.Format().BoxEnd
	}
	return bar
}

func prefix(frame *FrameContext, cols int) string {
	msg := frame.MessagePrefix

	if cols <= 3 {
		msg = strings.Repeat(".", cols)

	} else if len(msg) > cols {
		msg = msg[:cols-3] + "..."

	}
	// else if len(msg) < cols {
	// 	//msg += strings.Repeat(" ", cols-len(msg))
	// }

	return msg
}

func counter(frame *FrameContext, cols int) string {
	var counterBox string
	current := format.Format(frame.Current).To(frame.bar.Units).Width(frame.bar.UnitsWidth)
	if frame.Max > 0 {
		totalS := format.Format(frame.Max).To(frame.bar.Units).Width(frame.bar.UnitsWidth)
		counterBox = fmt.Sprintf("%s/%s", current, totalS)
	} else {
		counterBox = fmt.Sprintf("%s/?", current)
	}

	return counterBox
}

func timeSpent(frame *FrameContext, cols int) string {

	var left time.Duration
	var timeSpentBox string
	var fromStart time.Duration
	if frame.IsFinish {
		fromStart = frame.bar.finishedTime.Sub(frame.bar.startTime)
		left = (fromStart / time.Second) * time.Second
		timeSpentBox = left.String()
	} else {
		fromStart = time.Since(frame.bar.startTime)
		timeSpentBox = ((fromStart / time.Second) * time.Second).String()
	}

	return timeSpentBox
}

func speed(frame *FrameContext, cols int) string {
	var fromStart time.Duration
	if frame.IsFinish {
		fromStart = frame.bar.finishedTime.Sub(frame.bar.startTime)
	} else {
		fromStart = time.Since(frame.bar.startTime)
	}
	currentFromStart := frame.Current - frame.bar.startValue

	var speedBox string
	speed := float64(currentFromStart) / (float64(fromStart) / float64(time.Second))
	speedBox = format.Format(int64(speed)).To(frame.bar.Units).Width(frame.bar.UnitsWidth).PerSec().String()

	return speedBox
}

func percent(frame *FrameContext, cols int) string {
	percent := float64(frame.Current) / (float64(frame.Max) / 100.0)
	return fmt.Sprintf("%.02f%%", percent)
}

func timeLeft(frame *FrameContext, cols int) string {
	currentFromStart := frame.Current - frame.bar.startValue
	lastChangeTime := frame.bar.changeTime
	fromChange := lastChangeTime.Sub(frame.bar.startTime)

	var left time.Duration
	var timeLeftBox string
	select {
	case <-frame.bar.finish:
	default:
		if currentFromStart > 0 {
			perEntry := fromChange / time.Duration(currentFromStart)
			if frame.Max > 0 {
				left = time.Duration(frame.Max-frame.Current) * perEntry
				left -= time.Since(lastChangeTime)
				left = (left / time.Second) * time.Second
			}
			if left > 0 {
				timeLeftBox = format.Format(int64(left)).To(format.U_DURATION).String()
			}
		}
	}
	return timeLeftBox
}

func getSpinString(frame *FrameContext) string {
	return frame.Format().SpinString
}

func spin(frame *FrameContext, cols int) string {
	char := getSpinString(frame)[frame.frameNo%int64(len(getSpinString(frame)))]
	return string(char)
}

func frameNo(frame *FrameContext, cols int) string {
	return strconv.FormatInt(frame.FrameNo(), 10)
}
