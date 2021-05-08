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
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

// FromStruct creates a new Renderer for a struct type.
//
// It looks for the "json" key in the fields' tag and uses the name, if defined.
// As a special case, if the field tag is "-", the field is omitted.
// Options like "omitempty" are ignored.
func FromStruct(sep, delim string, typ reflect.Type) Renderer {
	fns, pns := jsonMetadata(typ)
	if len(pns) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.Indirect(reflect.ValueOf(i))
		b := &strings.Builder{}
		for idx := range fns {
			b.WriteString(pns[idx])
			b.WriteString(sep)
			b.WriteString(meta.ToString(v.FieldByName(fns[idx]).Interface()))
			b.WriteString(delim)
		}
		return b.String()[:b.Len()-len(delim)], nil
	})
}

// FromStructSlice creates a new Renderer for a struct slice.
//
// The fields are determined by the slice's element type.
// Like FromStruct, it looks for the "json" key in the struct field's tag.
func FromStructSlice(sep, delim string, typ reflect.Type) Renderer {
	_, pns := jsonMetadata(typ.Elem())
	if len(pns) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		b.WriteString(strings.Join(pns, sep))
		b.WriteString(delim)

		v := reflect.ValueOf(i)
		for idx := 0; idx < v.Len(); idx++ {
			e := reflect.Indirect(v.Index(idx))
			b.WriteString(meta.ToString(e.Field(0).Interface()))
			for pIdx := 1; pIdx < len(pns); pIdx++ {
				b.WriteString(sep)
				b.WriteString(meta.ToString(e.Field(pIdx).Interface()))
			}
			b.WriteString(delim)
		}
		return b.String()[:b.Len()-len(delim)], nil
	})
}
