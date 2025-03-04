// Copyright 2021 The gutenfmt authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gfmt

import (
	"fmt"
	"io"
	"reflect"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/styles"
)

// Writer is the interface that wraps the generic Write method.
type Writer interface {
	// Write writes i to the underlying output stream.
	//
	// It returns the number of bytes written and any error encountered that
	// caused the write to stop early.
	// Write must not modify the given parameter, even temporarily.
	Write(i any) (int, error)
}

// IOWriter wraps an io.Writer and implements the gfmt.Writer interface.
type IOWriter struct {
	writer io.Writer
}

var _ Writer = (*IOWriter)(nil)

// WrapIOWriter wraps the given io.Writer as a gfmt.Writer.
func WrapIOWriter(wrapper io.Writer) Writer {
	return IOWriter{wrapper}
}

// Write uses the default formats for its operand and writes to the internal io.Writer.
func (w IOWriter) Write(i any) (int, error) {
	return fmt.Fprint(w.writer, i)
}

// countingWriter counts the number of bytes written to the underlying Writer.
//
// Note that direct writes to the wrapped Writer do not increase the counter.
type countingWriter struct {
	writer io.Writer
	cnt    int
}

var _ io.Writer = (*countingWriter)(nil)
var _ io.StringWriter = (*countingWriter)(nil)

// wrapCountingWriter encapsulates the given Writer.
// If w is already a countingWriter, it is returned instead.
func wrapCountingWriter(w io.Writer) *countingWriter {
	if cw, ok := w.(*countingWriter); ok {
		return cw
	}
	return &countingWriter{w, 0}
}

// Write writes len(p) bytes from p to the underlying Writer.
// It returns the number of bytes written and any error encountered that caused
// the write to stop early.
// Additionally, it increases the total number of bytes written by n.
func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.writer.Write(p)
	cw.cnt += n
	return n, err
}

// WriteString writes s to the underlying Writer.
// It returns the number of bytes written and any error encountered that caused
// the write to stop early.
// Additionally, it increases the total number of bytes written by n.
func (cw *countingWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(cw.writer, s)
	cw.cnt += n
	return n, err
}

// highlight formats text in the given format and writes it to the Writer.
// If an unknown style is specified, a fallback style is used instead.
func highlight(w io.Writer, l chroma.Lexer, text string, s *chroma.Style) error {
	l = chroma.Coalesce(l)
	it, err := l.Tokenise(nil, text)
	if err != nil {
		return err
	}

	return formatters.TTY8.Format(w, s, it)
}

// Opt allows to customize the behavior of a Writer.
type Opt[W Writer] func(writer *W)

// WithPretty enables pretty-printing for the given Writer.
func WithPretty[W Writer]() Opt[W] {
	return func(w *W) {
		switch reflect.TypeOf(w) {
		case reflect.TypeOf(&JSON{}):
			any(w).(*JSON).Indent = "  "
		case reflect.TypeOf(&YAML{}):
			any(w).(*YAML).Indent = 2
		}
	}
}

// WithStyle sets the syntax highlighting style for the given Writer.
func WithStyle[W Writer](s *chroma.Style) Opt[W] {
	if s == styles.Fallback {
		s = nil
	}
	return func(w *W) {
		switch reflect.TypeOf(w) {
		case reflect.TypeOf(&JSON{}):
			any(w).(*JSON).Style = s
		case reflect.TypeOf(&YAML{}):
			any(w).(*YAML).Style = s
		}
	}
}
