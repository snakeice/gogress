package gogress

import "github.com/snakeice/gogress/format"

type FrameContext struct {
	bar     *Progress
	current int64
	max     int64
}

func (f *FrameContext) format() *format.ProgressFormat {
	return &f.bar.Format
}
