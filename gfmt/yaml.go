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
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/formatter"
	"github.com/abc-inc/gutenfmt/internal/render"
	"github.com/alecthomas/chroma/lexers/y"
	"gopkg.in/yaml.v3"
)

// YAML is a generic Writer that formats arbitrary values as YSON.
type YAML struct {
	w         io.Writer
	Formatter *formatter.CompFormatter
	Indent    int
	Style     string
}

// NewYAML creates a new YAML Writer.
func NewYAML(w io.Writer) *YAML {
	return &YAML{w, formatter.NewComp(), 2, ""}
}

// NewPrettyYAML creates a new YAML Writer with indentation and coloring.
func NewPrettyYAML(w io.Writer) *YAML {
	return &YAML{w, formatter.NewComp(), 2, "native"}
}

// Write writes the YAML representation of the given value to the underlying Writer.
func (w YAML) Write(i interface{}) (int, error) {
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
	e := yaml.NewEncoder(b)
	e.SetIndent(w.Indent)
	if err := e.Encode(i); err != nil {
		return 0, err
	}

	s := strings.TrimSuffix(b.String(), "\n")
	if w.Style == "" {
		return io.WriteString(w.w, s)
	}

	cw := wrapCountingWriter(w.w)
	if err := highlight(cw, y.YAML, s, w.Style); err != nil {
		return 0, err
	}
	return int(cw.cnt), nil
}
