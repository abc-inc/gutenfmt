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

	"github.com/abc-inc/gutenfmt/formatter"
)

// Text is a generic Writer that formats arbitrary values as plain text.
type Text struct {
	w         io.Writer
	Formatter *formatter.CompFormatter
	Sep       string
	Delim     string
}

// NewText creates a new text Writer.
func NewText(w io.Writer) *Text {
	return &Text{w, formatter.NewComp(), ":", "\n"}
}

// Write writes the text representation of the given value to the underlying Writer.
func (w Text) Write(i interface{}) (int, error) {
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
	} else if str, ok := i.(fmt.Stringer); ok {
		return fmt.Fprint(w.w, str.String())
	} else if !isContainerType(typ.Kind()) {
		return fmt.Fprint(w.w, i)
	}

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		return w.writeSlice(reflect.ValueOf(i))
	case reflect.Map:
		return w.writeMap(reflect.ValueOf(i).MapRange())
	default:
		return 0, formatter.ErrUnsupported
	}
}

// writeSlice writes the text representation of the given slice to the underlying Writer.
func (w Text) writeSlice(v reflect.Value) (int, error) {
	if v.Len() == 0 {
		return 0, nil
	}

	n, err := w.Write(v.Index(0).Interface())
	if err != nil {
		return n, err
	}

	cnt := n
	for idx := 1; idx < v.Len(); idx++ {
		n, err := io.WriteString(w.w, w.Delim)
		cnt += n
		if err != nil {
			return cnt, err
		}

		n, err = w.Write(v.Index(idx).Interface())
		cnt += n
		if err != nil {
			return cnt, err
		}
	}
	return cnt, nil
}

// writeMap writes the text representation of the given map to the underlying Writer.
func (w Text) writeMap(iter *reflect.MapIter) (int, error) {
	if !iter.Next() {
		return 0, nil
	}

	n, err := w.writeKeyVal(iter.Key().Interface(), iter.Value().Interface())
	if err != nil {
		return n, err
	}
	cnt := n

	for iter.Next() {
		n, err := io.WriteString(w.w, w.Delim)
		cnt += n
		if err != nil {
			return cnt, err
		}

		n, err = w.writeKeyVal(iter.Key().Interface(), iter.Value().Interface())
		cnt += n
		if err != nil {
			return cnt, err
		}
	}

	return cnt, nil
}

// writeKeyVal writes the text representation of a map entry to the underlying Writer.
func (w Text) writeKeyVal(k, v interface{}) (int, error) {
	cnt := 0

	n, err := w.Write(k)
	cnt += n
	if err != nil {
		return cnt, err
	}

	n, err = io.WriteString(w.w, w.Sep)
	cnt += n
	if err != nil {
		return cnt, err
	}

	n, err = w.Write(v)
	cnt += n
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}
