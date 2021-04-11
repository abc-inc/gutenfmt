package gfmt

import (
	"io"
	"reflect"
)

type InterfaceWriter interface {
	Write(i interface{}) (int, error)
}

type countingWriter struct {
	w   io.Writer
	cnt uint64
}

func newCountingWriter(w io.Writer) *countingWriter {
	if cw, ok := w.(*countingWriter); ok {
		return cw
	}
	return &countingWriter{w, 0}
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.cnt += uint64(n)
	return n, err
}

func (cw *countingWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(cw.w, s)
	cw.cnt += uint64(n)
	return n, err
}

func isCompositeType(typ reflect.Type) bool {
	k := typ.Kind()
	return k == reflect.Struct || k == reflect.Slice || k == reflect.Map || k == reflect.Array
}
