package gogress

import "github.com/snakeice/gogress/format"

type FrameContext struct {
	bar           *Progress
	Current       int64
	Max           int64
	Width         int
	elementNo     int
	usedWidth     int
	frameNo       int64
	MessagePrefix string
	IsFinish      bool
	recalc        []Decorator
}

func (f *FrameContext) Format() *format.ProgressFormat {
	return &f.bar.Format
}

func (f *FrameContext) FrameNo() int64 {
	return f.frameNo
}

func NewFrame(bar *Progress, current, max int64, width int, frameNo int64) *FrameContext {
	return &FrameContext{
		bar:           bar,
		Current:       current,
		Max:           max,
		Width:         width,
		elementNo:     0,
		usedWidth:     0,
		frameNo:       frameNo,
		MessagePrefix: bar.messagePrefix,
		IsFinish:      bar.isFinish,
		recalc:        []Decorator{},
	}
}

func (f *FrameContext) Copy() *FrameContext {
	return &FrameContext{
		bar:           f.bar,
		Current:       f.Current,
		Max:           f.Max,
		Width:         f.Width,
		elementNo:     0,
		frameNo:       f.FrameNo(),
		usedWidth:     0,
		MessagePrefix: f.MessagePrefix,
		IsFinish:      f.IsFinish,
	}
}
