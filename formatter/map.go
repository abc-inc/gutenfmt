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

package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/render"
)

// FromMap creates a Formatter that outputs all map entries in unspecified order.
func FromMap(sep, delim string) Formatter {
	return Func(func(i any) (string, error) {
		return FromMapKeys(sep, delim, reflect.ValueOf(i).MapKeys()...).Format(i)
	})
}

// FromMapKeys creates a Formatter that outputs entries for the given keys.
// Unlike FromMap, map entries are formatted in the specified order.
// If a key is given multiple times, it will be rendered multiple times.
func FromMapKeys(sep, delim string, ks ...reflect.Value) Formatter {
	return Func(func(i any) (string, error) {
		m := reflect.ValueOf(i)
		b := &strings.Builder{}
		for _, mk := range ks {
			mv := m.MapIndex(mk)
			s := ""
			if mv.IsValid() {
				s = render.ToString(mv.Interface())
			}
			n := render.ToString(mk.Interface())
			if _, err := fmt.Fprintf(b, "%s%s%s%s", n, sep, s, delim); err != nil {
				return "", err
			}
		}
		return b.String(), nil
	})
}

// FromMapSlice creates a Formatter that formats a slice of maps.
func FromMapSlice(sep, delim string) Formatter {
	contains := func(es []string, s string) bool {
		for _, e := range es {
			if e == s {
				return true
			}
		}
		return false
	}

	return Func(func(i any) (string, error) {
		v := reflect.ValueOf(i)
		if v.Len() == 0 {
			return "", nil
		}

		b := &strings.Builder{}

		var ks []string
		e := reflect.ValueOf(v.Index(0).Interface())
		for idx, k := range e.MapKeys() {
			if idx > 0 {
				b.WriteString(sep)
			}
			n := render.ToString(k.Interface())
			if !contains(ks, n) {
				ks = append(ks, n)
				b.WriteString(n)
			}
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			m := reflect.ValueOf(v.Index(i).Interface())
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := m.MapIndex(reflect.ValueOf(k))
				if v.IsValid() {
					b.WriteString(render.ToString(v))
				}
			}
		}

		return b.String(), nil
	})
}

// FromMapSliceKeys creates a Formatter that outputs a slice of maps.
func FromMapSliceKeys(sep, delim string, ks ...reflect.Value) Formatter {
	if len(ks) == 0 {
		return NoopFormatter()
	}

	return Func(func(i any) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		b.WriteString(render.ToString(ks[0].Interface()))
		for idx := 1; idx < len(ks); idx++ {
			b.WriteString(sep)
			b.WriteString(render.ToString(ks[idx].Interface()))
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := v.Index(i).MapIndex(k)
				if v.IsValid() {
					b.WriteString(render.ToString(v))
				}
			}
		}

		return b.String(), nil
	})
}
