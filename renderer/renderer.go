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

// Package renderer provides various implementations to convert values to strings.
package renderer

import (
	"errors"
	"reflect"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

// ErrUnsupported is the error resulting if a Renderer does not support the type.
var ErrUnsupported = errors.New("unsupported type")

// Renderer converts the given parameter to its string representation.
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

// Noop always returns an empty string and no error.
func Noop(_ interface{}) (s string, err error) {
	return
}

// NoopRenderer returns a simple Renderer that always returns an empty string and no error.
func NoopRenderer() Renderer { return RendererFunc(Noop) }

// CompRenderer combines multiple Renderers, each handling a different type.
type CompRenderer struct {
	byType map[string]Renderer
}

// NewComp creates and initializes a new CompRenderer.
func NewComp() *CompRenderer {
	return &CompRenderer{make(map[string]Renderer)}
}

// Renderer converts the given parameter to its string representation.
// If none of the registered Renderers can handle the given value, an error is returned.
func (cr CompRenderer) Render(i interface{}) (string, error) {
	if r, ok := cr.byType[meta.TypeName(reflect.TypeOf(i))]; ok {
		return r.Render(i)
	}
	return "", ErrUnsupported
}

// SetRenderer registers the Renderer for the given type.
// If a Renderer already exists for the type, it is replaced.
func (cr *CompRenderer) SetRenderer(n string, r Renderer) {
	cr.byType[n] = r
}

// SetRendererFunc registers the RendererFunc for the given type.
// If a Renderer already exists for the type, it is replaced.
func (cr *CompRenderer) SetRendererFunc(n string, r RendererFunc) {
	cr.byType[n] = r
}
