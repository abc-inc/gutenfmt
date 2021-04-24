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
	"fmt"
	"io"
	"reflect"

	"github.com/abc-inc/gutenfmt/internal/meta"
	"github.com/abc-inc/gutenfmt/renderer"
)

type Text struct {
	Sep      string
	Delim    string
	w        io.Writer
	Renderer *renderer.CompRenderer
}

func NewText(w io.Writer) *Text {
	return &Text{":", "\n", w, renderer.NewComp()}
}

func (f Text) Write(i interface{}) (int, error) {
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
	} else if str, ok := i.(fmt.Stringer); ok {
		return fmt.Fprint(f.w, str.String())
	} else if !meta.IsContainerType(typ.Kind()) {
		return fmt.Fprint(f.w, i)
	}

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		return f.writeSlice(reflect.ValueOf(i))
	case reflect.Map:
		return f.writeMap(reflect.ValueOf(i).MapRange())
	default:
		return 0, renderer.ErrUnsupported
	}
}

func (f Text) writeSlice(v reflect.Value) (int, error) {
	if v.Len() == 0 {
		return 0, nil
	}

	n, err := f.Write(v.Index(0).Interface())
	if err != nil {
		return n, err
	}

	cnt := n
	for idx := 1; idx < v.Len(); idx++ {
		n, err := io.WriteString(f.w, f.Delim)
		cnt += n
		if err != nil {
			return cnt, err
		}

		n, err = f.Write(v.Index(idx).Interface())
		cnt += n
		if err != nil {
			return cnt, err
		}
	}
	return cnt, nil
}

func (f Text) writeMap(iter *reflect.MapIter) (int, error) {
	if !iter.Next() {
		return 0, nil
	}

	n, err := f.writeKeyVal(iter.Key().Interface(), iter.Value().Interface())
	if err != nil {
		return n, err
	}
	cnt := n

	for iter.Next() {
		n, err := io.WriteString(f.w, f.Delim)
		cnt += n
		if err != nil {
			return cnt, err
		}

		n, err = f.writeKeyVal(iter.Key().Interface(), iter.Value().Interface())
		cnt += n
		if err != nil {
			return cnt, err
		}
	}

	return cnt, nil
}

func (f Text) writeKeyVal(k, v interface{}) (int, error) {
	cnt := 0

	n, err := f.Write(k)
	cnt += n
	if err != nil {
		return cnt, err
	}

	n, err = io.WriteString(f.w, f.Sep)
	cnt += n
	if err != nil {
		return cnt, err
	}

	n, err = f.Write(v)
	cnt += n
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}
