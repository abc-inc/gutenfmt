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

package gfmt_test

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func Test_Write_Types(t *testing.T) {
	type data struct {
		kind reflect.Kind
		in   interface{}
		out  string
	}
	type unit struct {
		b *strings.Builder
		w GenericWriter
	}

	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2

	tests := []data{
		{reflect.Bool, true, "true"},
		{reflect.Int, int(-4), "-4"},
		{reflect.Int8, int8(-8), "-8"},
		{reflect.Int16, int16(-16), "-16"},
		{reflect.Int32, int32(-32), "-32"},
		{reflect.Int64, int64(-64), "-64"},
		{reflect.Uint, uint(4), "4"},
		{reflect.Uint8, uint8(8), "8"},
		{reflect.Uint16, uint16(16), "16"},
		{reflect.Uint32, uint32(32), "32"},
		{reflect.Uint64, uint64(64), "64"},
		{reflect.Uintptr, uintptr(128), "128"},
		{reflect.Float32, math.MaxFloat32, fmt.Sprint(math.MaxFloat32)},
		{reflect.Float64, math.MaxFloat64, fmt.Sprint(math.MaxFloat64)},
		// {reflect.Complex64, complex64(-2.71i), "(0-2.71i)"},   // not supported by JSON
		// {reflect.Complex128, complex128(-3.14i), "(0-3.14i)"}, // not supported by JSON
		// {reflect.Array, [3]string{"a", "b", "c"}, "[a b c]"},
		// {reflect.Chan, true, "true"},                          // does not make sense
		// {reflect.Func, Test_Write_Types, "Test_Write_Types"},  // not supported by JSON
		// {reflect.Interface, true, "true"},
		// {reflect.Map, m, "{}"},
		// {reflect.Ptr, true, "true"},
		// {reflect.Slice, true, "true"},
		{reflect.String, "str", "str"},
		// {reflect.Struct, true, "true"},
		// {reflect.UnsafePointer, true, "true"},
	}

	for _, tt := range tests {
		sPrettyJSON := &strings.Builder{}
		sJSON := &strings.Builder{}
		// sTab := &strings.Builder{}
		sText := &strings.Builder{}

		prettyJSON := NewPrettyJSON(sPrettyJSON)
		json := NewJSON(sJSON)
		// tab := NewTab(sTab)
		text := NewText(sText)

		us := []unit{
			{sPrettyJSON, prettyJSON},
			{sJSON, json},
			// {sTab, tab},
			{sText, text},
		}

		for _, u := range us {
			t.Run(tt.kind.String()+"_"+reflect.TypeOf(u.w).Elem().Name(), func(t *testing.T) {
				want := tt.out

				postProc := func(s string) string { return s }
				if _, ok := u.w.(*Text); ok {
					postProc = func(s string) string { return strings.Trim(strings.ReplaceAll(s, ",", "\n"), "[]") }
				} else if _, ok := u.w.(*Tab); ok {
					postProc = func(s string) string { return strings.ReplaceAll(s, ",", "\n") }
				}

				_, err := u.w.Write(tt.in)
				NoError(t, err)
				Equal(t, want, u.b.String())

				if f, ok := u.w.(*JSON); ok && f.Style != "" {
					// pretty JSON is too hard to verify, so we skip further tests
					return
				}

				if _, ok := u.w.(*JSON); ok && tt.in == want {
					// JSON quotes strings, so the expected output needs to be quoted
					want = "\"" + want + "\""
				}

				// array
				u.b.Reset()
				_, err = u.w.Write([2]interface{}{tt.in, tt.in})
				NoError(t, err)
				Equal(t, postProc(fmt.Sprintf("[%s,%s]", want, want)), u.b.String())

				// slice
				u.b.Reset()
				_, err = u.w.Write([]interface{}{tt.in, tt.in})
				NoError(t, err)
				Equal(t, postProc(fmt.Sprintf("[%s,%s]", want, want)), u.b.String())

				// map
				if _, ok := u.w.(*JSON); !ok {
					// JSON does not support arbitrary maps

					u.b.Reset()
					_, err = u.w.Write(map[interface{}]interface{}{tt.in: tt.in})
					NoError(t, err)
					Equal(t, postProc(fmt.Sprintf("%s:%s", want, want)), u.b.String())
				}
			})
		}
	}
}
