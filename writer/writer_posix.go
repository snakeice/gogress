// +build !windows

package writer

import (
	"fmt"
)

func (w *Writer) clearLines() {
	fmt.Fprintf(w.out, cuuAndEd, w.lineCount)
}
