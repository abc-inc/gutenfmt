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
	"os"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/formatter"
	"github.com/abc-inc/gutenfmt/internal/json"
	"github.com/abc-inc/gutenfmt/internal/render"
	"github.com/alecthomas/chroma/lexers/j"
	"github.com/mattn/go-isatty"
)

// JSON is a generic Writer that formats arbitrary values as JSON.
type JSON struct {
	w         io.Writer
	Formatter *formatter.CompFormatter
	Indent    string
	Style     string
}

// NewAutoJSON creates and initializes new JSON Writer with or without formatting.
// The JSON encoder uses indentation and ANSII escape sequences are used for coloring,
// if the underlying Writer is stdout on an interactive terminal.
func NewAutoJSON(w io.Writer) *JSON {
	if w == os.Stdout && isatty.IsTerminal(os.Stdout.Fd()) {
		return NewPrettyJSON(w)
	}
	return NewJSON(w)
}

// NewJSON creates a new JSON Writer.
func NewJSON(w io.Writer) *JSON {
	return &JSON{w, formatter.NewComp(), "", ""}
}

// NewPrettyJSON creates a new JSON Writer with indentation and coloring.
func NewPrettyJSON(w io.Writer) *JSON {
	return &JSON{w, formatter.NewComp(), "  ", "native"}
}

// Write writes the JSON representation of the given value to the underlying Writer.
func (w JSON) Write(i any) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := w.Formatter.Format(i); err == nil {
		return io.WriteString(w.w, s)
	} else if !errors.Is(err, formatter.ErrUnsupported) {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return w.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !isContainerType(typ.Kind()) {
		return fmt.Fprint(w.w, render.ToString(i))
	}

	b := &strings.Builder{}
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.SetIndent("", w.Indent)
	if err := e.Encode(i); err != nil {
		return 0, err
	}

	s := strings.TrimSuffix(b.String(), "\n")
	if w.Style == "" {
		return io.WriteString(w.w, s)
	}

	cw := wrapCountingWriter(w.w)
	if err := highlight(cw, j.JSON, s, w.Style); err != nil {
		return 0, err
	}
	return int(cw.cnt), nil
}
