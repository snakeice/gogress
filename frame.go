package gogress

import "github.com/snakeice/gogress/format"

type FrameContext struct {
	Bar       *Progress
	Current   int64
	Max       int64
	Width     int
	elementNo int
	usedWidth int
	frameNo   int64
}

func (f *FrameContext) format() *format.ProgressFormat {
	return &f.Bar.Format
}

func (f *FrameContext) SpinString() string {
	return f.Bar.Format.SpinString
}

func (f *FrameContext) FrameNo() int64 {
	return f.frameNo
}

func NewFrame(bar *Progress, current, max int64, width int, frameNo int64) *FrameContext {
	return &FrameContext{
		Bar:       bar,
		Current:   current,
		Max:       max,
		Width:     width,
		elementNo: 0,
		usedWidth: 0,
		frameNo:   frameNo,
	}
}
