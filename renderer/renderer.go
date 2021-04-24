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

package renderer

import (
	"errors"
	"reflect"
	"strings"
	"text/template"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

var ErrUnsupported = errors.New("unsupported type")

//Renderer converts the given parameter to a string representation.
type Renderer interface {
	Render(i interface{}) (string, error)
}

// RendererFunc is an adapter to allow the use of ordinary functions as Renderers.
// If f is a function with the appropriate signature,
// RendererFunc(f) is a Renderer that calls f.
type RendererFunc func(i interface{}) (string, error)

func (r RendererFunc) Render(i interface{}) (string, error) {
	return r(i)
}

//Noop always returns an empty string and no error.
func Noop(i interface{}) (s string, err error) {
	return
}

// NoopRenderer returns a simple Renderer
// that always returns an empty string and no error.
func NoopRenderer() Renderer { return RendererFunc(Noop) }

type CompRenderer struct {
	byType map[string]Renderer
}

func NewComp() *CompRenderer {
	return &CompRenderer{make(map[string]Renderer)}
}

func (cr CompRenderer) Render(i interface{}) (string, error) {
	r, ok := cr.byType[meta.TypeName(reflect.TypeOf(i))]
	if !ok {
		return "", ErrUnsupported
	}
	return r.Render(i)
}

func (cr *CompRenderer) SetRenderer(n string, r Renderer) {
	cr.byType[n] = r
}

func (cr *CompRenderer) SetRendererFunc(n string, r RendererFunc) {
	cr.byType[n] = r
}

func FromTemplate(tmpl template.Template) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		if err := tmpl.Execute(b, i); err != nil {
			return "", err
		}
		return b.String(), nil
	})
}
