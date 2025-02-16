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
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/render"
	"github.com/abc-inc/gutenfmt/meta"
)

// FromStruct creates a new Formatter for a struct type.
func FromStruct(sep, delim string, typ reflect.Type) Formatter {
	fs := meta.Resolve(typ)
	if len(fs) == 0 {
		return NoopFormatter()
	}

	return Func(func(i any) (string, error) {
		v := reflect.Indirect(reflect.ValueOf(i))
		b := &strings.Builder{}
		for _, f := range fs {
			b.WriteString(f.Name)
			b.WriteString(sep)
			if !v.FieldByName(f.Field).IsZero() {
				b.WriteString(render.ToString(v.FieldByName(f.Field).Interface()))
			}
			b.WriteString(delim)
		}
		return b.String()[:b.Len()-len(delim)], nil
	})
}

// FromStructSlice creates a new Formatter for a struct slice.
// The fields are determined by the slice's element type.
func FromStructSlice(sep, delim string, typ reflect.Type) Formatter {
	fs := meta.Resolve(typ.Elem())
	if len(fs) == 0 {
		return NoopFormatter()
	}

	return Func(func(i any) (string, error) {
		b := &strings.Builder{}
		for _, f := range fs {
			b.WriteString(sep)
			b.WriteString(f.Name)
		}
		b.WriteString(delim)

		v := reflect.ValueOf(i)
		for idx := 0; idx < v.Len(); idx++ {
			e := reflect.Indirect(v.Index(idx))
			b.WriteString(render.ToString(e.Field(0).Interface()))
			for pIdx := 1; pIdx < len(fs); pIdx++ {
				b.WriteString(sep)
				b.WriteString(render.ToString(e.Field(pIdx).Interface()))
			}
			b.WriteString(delim)
		}
		return b.String()[len(sep) : b.Len()-len(delim)], nil
	})
}
