package gogress

import (
	"io"
)

type Reader struct {
	io.Reader
	bar *Progress
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Add(n)
	return
}

func (r *Reader) Close() (err error) {
	r.bar.Finish()
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}

type Writer struct {
	io.Writer
	bar *Progress
}

func (r *Writer) Write(p []byte) (n int, err error) {
	n, err = r.Writer.Write(p)
	r.bar.Add(n)
	return
}

func (r *Writer) Close() (err error) {
	r.bar.Finish()
	if closer, ok := r.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return
}
