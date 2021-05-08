/**
 * Copyright 2021 The gutenfmt authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gfmt

import (
	"io"
)

// GenericWriter is the interface that wraps the generic Write method.
type GenericWriter interface {
	// Write writes i to the underlying output stream.
	//
	// It returns the number of bytes written and any error encountered that
	// caused the write to stop early.
	// Write must not modify the given parameter, even temporarily.
	Write(i interface{}) (int, error)
}

// countingWriter counts the number of bytes written to the underlying writer.
//
// Note that direct writes to the wrapped Writer do not increase the counter.
type countingWriter struct {
	w   io.Writer
	cnt uint64
}

// wrapCountingWriter encapsulates the given Writer.
//
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
	n, err := cw.w.Write(p)
	cw.cnt += uint64(n)
	return n, err
}

// WriteString writes s to the underlying Writer.
// It returns the number of bytes written and any error encountered that caused
// the write to stop early.
// Additionally, it increases the total number of bytes written by n.
func (cw *countingWriter) WriteString(s string) (n int, err error) {
	n, err = io.WriteString(cw.w, s)
	cw.cnt += uint64(n)
	return n, err
}
