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
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/formatter"
	"github.com/abc-inc/gutenfmt/internal/json"
	"github.com/abc-inc/gutenfmt/internal/render"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// JSON is a generic Writer that formats arbitrary values as JSON.
type JSON struct {
	writer    io.Writer
	Formatter *formatter.CompFormatter
	Indent    string
	Style     *chroma.Style
}

// NewJSON creates a new JSON Writer.
func NewJSON(w io.Writer, opts ...Opt[JSON]) *JSON {
	gw := &JSON{writer: w, Formatter: formatter.NewComp()}
	for _, opt := range opts {
		opt(gw)
	}
	return gw
}

// Write writes the JSON representation of the given value to the underlying Writer.
func (w JSON) Write(i any) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := w.Formatter.Format(i); err == nil {
		return io.WriteString(w.writer, s)
	} else if !errors.Is(err, formatter.ErrUnsupported) {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return w.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !isContainerType(typ.Kind()) {
		return fmt.Fprint(w.writer, render.ToString(i))
	}

	b := &strings.Builder{}
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.SetIndent("", w.Indent)
	if err := e.Encode(i); err != nil {
		return 0, err
	}

	s := strings.TrimSuffix(b.String(), "\n")
	if w.Style == nil || w.Style.Name == "noop" {
		// Take a shortcut if no syntax highlighting should be applied.
		// This is about 50x faster than having chroma to tokenize the JSON.
		return io.WriteString(w.writer, s)
	}

	cw := wrapCountingWriter(w.writer)
	if err := highlight(w.writer, lexers.Get("json"), s, w.Style); err != nil {
		return 0, err
	}
	return cw.cnt, nil
}
