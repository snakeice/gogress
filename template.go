package gogress

import (
	"text/template"
)

/*
 +-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----+
 | C1  | C2  | C3  | C4  | C5  | C6  | C7  | C8  | C9  | C10 | C11 | C12 |
 +-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----+

	Like bootstrap grid system
*/

// type RowPart struct {
// 	decorator Decorator
// 	cols      int
// 	aux       interface{}
// }

type Template struct {
}

func NewTemplate(templateString string) (*template.Template, error) {
	// var result []RowPart
	t := template.New("")
	t.Funcs(Colors)
	t.Funcs(Decorators)
	if _, err := t.Parse(templateString); err != nil {
		return nil, err
	}

	return t, nil
}
