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
	"text/tabwriter"

	"github.com/abc-inc/gutenfmt/internal/meta"
	"github.com/abc-inc/gutenfmt/renderer"
)

type Tab struct {
	w        *countingWriter
	Renderer *renderer.CompRenderer
}

func NewTab(w io.Writer) *Tab {
	return &Tab{wrapCountingWriter(w), renderer.NewComp()}
}

func (f Tab) Write(i interface{}) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := f.Renderer.Render(i); err == nil {
		return f.w.WriteString(s)
	} else if err != renderer.ErrUnsupported {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return f.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !meta.IsContainerType(typ.Kind()) {
		return fmt.Fprint(f.w, i)
	}

	f.w.cnt = 0
	tw := tabwriter.NewWriter(f.w, 1, 4, 3, ' ', 0)

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		_, err := f.writeSlice(tw, reflect.ValueOf(i))
		return int(f.w.cnt), err
	case reflect.Map:
		_, err := f.writeMap(tw, i)
		return int(f.w.cnt), err
	default:
		_, err := f.writeStruct(tw, i)
		return int(f.w.cnt), err
	}
}

func (f Tab) writeSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	if v.Type().Elem().Kind() == reflect.Struct {
		n, err := f.writeStructSlice(tw, v)
		return n, err
	} else if v.Type().Elem().Kind() == reflect.Map {
		n, err := f.writeMapSlice(tw, v)
		return n, err
	}

	if v.Len() == 0 {
		return 0, nil
	}

	n, err := f.w.WriteString(meta.StrVal(v.Index(0).Interface()))
	if err != nil {
		return n, err
	}

	cnt := n
	for idx := 1; idx < v.Len(); idx++ {
		n, err = f.w.WriteString("\n" + meta.StrVal(v.Index(idx).Interface()))
		cnt += n
		if err != nil {
			return cnt, err
		}
	}
	return cnt, tw.Flush()
}

func (f Tab) writeMap(tw *tabwriter.Writer, i interface{}) (int, error) {
	r := renderer.FromMap("\t", "\t\n")
	return renderer.RenderTab(tw, r, i)
}

func (f Tab) writeMapSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	r := renderer.FromMapSlice("\t", "\t\n", v.Type())
	return renderer.RenderTab(tw, r, v.Interface())
}

func (f Tab) writeStruct(tw *tabwriter.Writer, i interface{}) (int, error) {
	r := renderer.FromStruct("\t", "\t\n", reflect.TypeOf(i))
	return renderer.RenderTab(tw, r, i)
}

func (f Tab) writeStructSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	r := renderer.FromStructSlice("\t", "\t\n", v.Type())
	return renderer.RenderTab(tw, r, v.Interface())
}
