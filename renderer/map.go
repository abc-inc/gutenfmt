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
	"fmt"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

// FromMap creates a Renderer that outputs all map entries in unspecified order.
func FromMap(sep, delim string) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		return FromMapKeys(sep, delim, reflect.ValueOf(i).MapKeys()...).Render(i)
	})
}

// FromMapKeys creates a Renderer that outputs entries for the given keys.
// Unlike FromMap, map entries are rendered in the specified order.
// If key is given multiple times, it will be rendered multiple times.
// If a key is not present
func FromMapKeys(sep, delim string, ks ...reflect.Value) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		m := reflect.ValueOf(i)
		b := &strings.Builder{}
		for _, mk := range ks {
			mv := m.MapIndex(mk)
			str := ""
			if mv.IsValid() {
				str = meta.ToString(mv.Interface())
			}
			n := meta.ToString(mk.Interface())
			if _, err := fmt.Fprintf(b, "%s%s%s%s", n, sep, str, delim); err != nil {
				return "", err
			}
		}
		return b.String(), nil
	})
}

// FromMapSlice creates a Renderer that renders a slice of maps.
func FromMapSlice(sep, delim string, _ reflect.Type) Renderer {
	contains := func(es []string, s string) bool {
		for _, e := range es {
			if e == s {
				return true
			}
		}
		return false
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		if v.Len() == 0 {
			return "", nil
		}

		b := &strings.Builder{}

		var ks []string
		e := v.Index(0)
		for idx, k := range e.MapKeys() {
			if idx > 0 {
				b.WriteString(sep)
			}
			n := meta.ToString(k.Interface())
			if !contains(ks, n) {
				ks = append(ks, n)
				b.WriteString(n)
			}
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			m := v.Index(i)
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := m.MapIndex(reflect.ValueOf(k))
				if v.IsValid() {
					b.WriteString(meta.ToString(v))
				}
			}
		}

		return b.String(), nil
	})
}

// FromMapSliceKeys creates a renderer that outputs a slice of maps.
func FromMapSliceKeys(sep, delim string, ks ...reflect.Value) Renderer {
	if len(ks) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		b.WriteString(meta.ToString(ks[0].Interface()))
		for idx := 1; idx < len(ks); idx++ {
			b.WriteString(sep)
			b.WriteString(meta.ToString(ks[idx].Interface()))
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := v.Index(i).MapIndex(k)
				if v.IsValid() {
					b.WriteString(meta.ToString(v))
				}
			}
		}

		return b.String(), nil
	})
}
