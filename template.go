package gogress

import (
	"bytes"
	"sync/atomic"
	"text/template"
)

const DefaultTemplate = `{{spin . 1}}{{prefix . 2}}{{bar . 3}}{{percent . 1}}{{counter . 1}}{{speed . 1}}{{timeSpent . 1}}{{timeLeft . 1}}{{frameNo . 1}}`

/*
	"bar":       bar,
	"prefix":    prefix,
	"counter":   counter,
	"timeSpent": timeSpent,
	"speed":     speed,
	"percent":   percent,
	"timeLeft":  timeLeft,
*/

/*
 +-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----+
 | C1  | C2  | C3  | C4  | C5  | C6  | C7  | C8  | C9  | C10 | C11 | C12 |
 +-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----+

	Like bootstrap grid system
*/
type TemplateParser struct {
	template    *template.Template
	lastContext *FrameContext
	frameNo     int64
	width       int
	lastFrame   string
}

func NewTemplateParserEmpty() *TemplateParser {
	tp := &TemplateParser{
		template: template.New(""),
	}
	tp.template.Funcs(Colors)
	tp.template.Funcs(Decorators)
	_ = tp.UpdateTemplate(DefaultTemplate)
	return tp
}

func NewTemplateParser(templateString string) (*TemplateParser, error) {
	tp := NewTemplateParserEmpty()
	err := tp.UpdateTemplate(templateString)
	return tp, err
}

func (tp *TemplateParser) Last() string {
	return tp.lastFrame
}

func (tp *TemplateParser) UpdateFrame(frame *FrameContext) {
	if frame.Bar.isFinish && frame.Width != tp.width {
		if tp.lastContext == nil {
			tp.lastContext.Width = tp.width
			tp.parseContext(tp.lastContext)
		}

	} else {
		tp.parseContext(frame)
	}
}

func (tp *TemplateParser) UpdateTemplate(templateString string) error {
	if _, err := tp.template.Parse(templateString); err != nil {
		return err
	}
	return nil
}

func (tp *TemplateParser) parseContext(frame *FrameContext) {
	atomic.AddInt64(&tp.frameNo, 1)
	tp.lastContext = frame
	var tpl bytes.Buffer
	if err := tp.template.Execute(&tpl, frame); err != nil {
		tp.lastFrame = err.Error()
	} else {
		tp.lastFrame = tpl.String()
	}
}
