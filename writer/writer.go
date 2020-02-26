package writer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var cuuAndEd = fmt.Sprintf("%c[%%dA%[1]c[J", 27)
var notTTY = errors.New("Not a terminal")

type Writer struct {
	out        io.Writer
	buffer     bytes.Buffer
	lineCount  int
	fd         uintptr
	isTerminal bool
}

func New(out io.Writer) *Writer {
	var writer = &Writer{}
	writer.out = out
	if f, ok := out.(*os.File); ok {
		writer.fd = f.Fd()
		writer.isTerminal = terminal.IsTerminal(int(writer.fd))
	}
	return writer
}

func (w *Writer) GetWidth() (int, error) {
	if w.isTerminal {
		width, _, err := terminal.GetSize(int(w.fd))
		return width, err
	}
	return 0, notTTY
}

func (w *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	return w.buffer.ReadFrom(r)
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *Writer) WriteString(s string) (n int, err error) {
	return w.buffer.WriteString(s)
}

func (w *Writer) Flush(lineCount int) (err error) {
	if w.lineCount > 0 {
		w.clearLines()
	}
	w.lineCount = lineCount
	_, err = w.buffer.WriteTo(w.out)
	w.buffer.Reset()
	return err
}
