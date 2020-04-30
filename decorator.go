package gogress

import (
	"math/rand"
	"text/template"

	"github.com/fatih/color"
)

type Decorator func(frame *FrameContext, cols int) string

const (
	adElPlaceholder    = "%_ad_el_%"
	adElPlaceholderLen = len(adElPlaceholder)
)

var Decorators = template.FuncMap{
	"bar":       adaptativeDecorator(bar),
	"prefix":    wrapDecorator(prefix),
	"counter":   wrapDecorator(counter),
	"timeSpent": wrapDecorator(timeSpent),
	"speed":     wrapDecorator(speed),
	"percent":   wrapDecorator(percent),
	"timeLeft":  wrapDecorator(timeLeft),
	"spin":      wrapDecorator(spin),
	"frameNo":   wrapDecorator(frameNo),
}

var Colors = template.FuncMap{
	"black":      color.New(color.FgBlack).SprintFunc(),
	"red":        color.New(color.FgRed).SprintFunc(),
	"green":      color.New(color.FgGreen).SprintFunc(),
	"yellow":     color.New(color.FgYellow).SprintFunc(),
	"blue":       color.New(color.FgBlue).SprintFunc(),
	"magenta":    color.New(color.FgMagenta).SprintFunc(),
	"cyan":       color.New(color.FgCyan).SprintFunc(),
	"white":      color.New(color.FgWhite).SprintFunc(),
	"resetcolor": color.New(color.Reset).SprintFunc(),
	"rndcolor":   rndcolor,
}

func AddDecorator(name string, decorator Decorator) {
	Decorators[name] = wrapDecorator(decorator)
}

func getColWidth(total int) int {
	return total / 12
}

func wrapDecorator(decorator Decorator) Decorator {
	return Decorator(func(frame *FrameContext, colsGrid int) string {
		frame.Width -= 1 // line break
		frame.elementNo += 1
		cols := getColWidth(frame.Width) * colsGrid
		frame.usedWidth += cols
		response := decorator(frame, cols)
		if len(response) >= cols {
			return response[:cols]
		} else {
			return response
		}
	})
}

func adaptativeDecorator(decorator Decorator) Decorator {
	return Decorator(func(frame *FrameContext, cols int) string {
		frame.recalc = append(frame.recalc, decorator)
		return adElPlaceholder
	})
}

func RemoveDecorator(name string) {
	delete(Decorators, name)
}

func rndcolor(s string) string {
	c := rand.Intn(int(color.FgWhite-color.FgBlack)) + int(color.FgBlack) - 1
	return color.New(color.Attribute(c)).Sprint(s)
}
