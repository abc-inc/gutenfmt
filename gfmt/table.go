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

package gfmt

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"github.com/abc-inc/gutenfmt/formatter"
	"github.com/abc-inc/gutenfmt/internal/render"
)

// Tab is a generic Writer that formats arbitrary values as ASCII table.
type Tab struct {
	cw        *countingWriter
	Formatter *formatter.CompFormatter
}

// NewTab creates a new table Writer.
func NewTab(w io.Writer) *Tab {
	return &Tab{wrapCountingWriter(w), formatter.NewComp()}
}

// Write formats the given value as a table and writes it to the underlying Writer.
func (w Tab) Write(i any) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := w.Formatter.Format(i); err == nil {
		return w.cw.WriteString(s)
	} else if !errors.Is(err, formatter.ErrUnsupported) {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return w.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !isContainerType(typ.Kind()) {
		return fmt.Fprint(w.cw, i)
	}

	w.cw.cnt = 0
	tw := tabwriter.NewWriter(w.cw, 4, 4, 1, ' ', 0)

	switch typ.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Array:
		_, err := w.writeSlice(tw, reflect.ValueOf(i))
		return w.cw.cnt, err
	case reflect.Map:
		_, err := w.writeMap(tw, i)
		return w.cw.cnt, err
	default:
		_, err := w.writeStruct(tw, i)
		return w.cw.cnt, err
	}
}

// writeSlice formats a slice of any type to a string.
func (w Tab) writeSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	if v.Type().Elem().Kind() == reflect.Struct ||
		(v.Type().Elem().Kind() == reflect.Ptr && v.Type().Elem().Elem().Kind() == reflect.Struct) {

		return w.writeStructSlice(tw, v)
	}
	if v.Len() == 0 {
		return 0, nil
	}
	if reflect.TypeOf(reflect.Indirect(v.Index(0)).Interface()).Kind() == reflect.Map {
		return w.writeMapSlice(tw, v)
	}

	cnt, err := w.cw.WriteString(render.ToString(v.Index(0).Interface()))
	if err != nil {
		return cnt, err
	}

	for idx := 1; idx < v.Len(); idx++ {
		n, err := w.cw.WriteString("\n" + render.ToString(v.Index(idx).Interface()))
		cnt += n
		if err != nil {
			return cnt, err
		}
	}
	return cnt, tw.Flush()
}

// writeMap formats a map to a tabular string representation.
func (w Tab) writeMap(tw *tabwriter.Writer, i any) (int, error) {
	f := formatter.FromMap("\t", "\t\n")
	return formatter.FormatTab(tw, f, i)
}

// writeMapSlice formats a map slice to a tabular string representation.
func (w Tab) writeMapSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	f := formatter.FromMapSlice("\t", "\t\n")
	return formatter.FormatTab(tw, f, v.Interface())
}

// writeStruct formats a struct to a tabular string representation.
func (w Tab) writeStruct(tw *tabwriter.Writer, i any) (int, error) {
	f := formatter.FromStruct("\t", "\t\n", reflect.TypeOf(i))
	return formatter.FormatTab(tw, f, i)
}

// writeStructSlice formats a struct slice to a tabular string representation.
func (w Tab) writeStructSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	f := formatter.FromStructSlice("\t", "\t\n", v.Type())
	return formatter.FormatTab(tw, f, v.Interface())
}
