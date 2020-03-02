package gogress

import (
	"math/rand"
	"text/template"

	"github.com/fatih/color"
)

type Decorator func(frame *FrameContext, gridCols, cols int) string

var Decorators = template.FuncMap{
	"bar":       bar,
	"prefix":    prefix,
	"counter":   1,
	"timeSpent": 1,
	"speed":     1,
	"percent":   1,
	"timeLeft":  1,
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
	"rnd":        rnd,
}

func AddDecorator(name string, decorator *Decorator) {
	Decorators[name] = decorator
}

func RemoveDecorator(name string) {
	delete(Decorators, name)
}

func rndcolor(s string) string {
	c := rand.Intn(int(color.FgWhite-color.FgBlack)) + int(color.FgBlack)
	return color.New(color.Attribute(c)).Sprint(s)
}

func rnd(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	return args[rand.Intn(len(args))]
}
