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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/meta"
	"github.com/abc-inc/gutenfmt/renderer"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers/j"
	"github.com/alecthomas/chroma/styles"
	"github.com/mattn/go-isatty"
)

type JSON struct {
	w        io.Writer
	Renderer *renderer.CompRenderer
	Indent   string
	Style    string
}

func NewAutoJSON(w io.Writer) *JSON {
	if w == os.Stdout && isatty.IsTerminal(os.Stdout.Fd()) {
		return NewPrettyJSON(w)
	}
	return NewJSON(w)
}

func NewJSON(w io.Writer) *JSON {
	return &JSON{w, renderer.NewComp(), "", ""}
}

func NewPrettyJSON(w io.Writer) *JSON {
	return &JSON{w, renderer.NewComp(), "  ", "native"}
}

func (f JSON) Write(i interface{}) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := f.Renderer.Render(i); err == nil {
		return io.WriteString(f.w, s)
	} else if err != renderer.ErrUnsupported {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return f.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !meta.IsContainerType(typ.Kind()) {
		return fmt.Fprint(f.w, meta.StrVal(i))
	}

	b := &strings.Builder{}
	m := json.NewEncoder(b)
	m.SetEscapeHTML(false)
	m.SetIndent("", f.Indent)
	if err := m.Encode(i); err != nil {
		return 0, err
	}

	s := strings.TrimSuffix(b.String(), "\n")
	if f.Style == "" {
		return io.WriteString(f.w, s)
	}

	w := wrapCountingWriter(f.w)
	if err := f.highlight(f.w, s, f.Style); err != nil {
		return 0, err
	}
	return int(w.cnt), nil
}

func (f JSON) highlight(w io.Writer, source, style string) error {
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	l := chroma.Coalesce(j.JSON)
	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return formatters.TTY8.Format(w, s, it)
}
