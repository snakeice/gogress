package gogress

import (
	"bytes"
	"strings"
	"sync/atomic"
	"text/template"
)

const DefaultTemplate = `{{prefix . 2}} {{bar . 5}} {{percent . 1}} {{counter . 1}} {{speed . 1}} {{timeSpent . 1}} {{timeLeft . 1}}`

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
	if tp.lastContext != nil && frame.Width != tp.lastContext.Width {
		atomic.AddInt64(&tp.frameNo, 1)
		tp.lastContext.Width = frame.Width
		tp.parseContext(frame)
	} else {
		atomic.AddInt64(&tp.frameNo, 1)
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
	tp.lastContext = frame.Copy()
	var tpl bytes.Buffer
	if err := tp.template.Execute(&tpl, frame); err != nil {
		tp.lastFrame = err.Error()
	} else {
		tp.lastFrame = tpl.String()

		toRecalc := len(frame.recalc)
		if toRecalc == 0 {
			return
		}
		staticWidth := len(tp.lastFrame) - (toRecalc * adElPlaceholderLen)

		if frame.Width-staticWidth <= 0 {
			tp.lastFrame = strings.ReplaceAll(tp.lastFrame, adElPlaceholder, "")
		} else {
			max := (frame.Width - staticWidth) / toRecalc
			for _, element := range frame.recalc {
				tp.lastFrame = strings.Replace(tp.lastFrame, adElPlaceholder, element(frame, max), 1)
			}
		}
	}
}
