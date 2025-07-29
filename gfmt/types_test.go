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

package gfmt_test

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/stretchr/testify/require"
)

func Test_Write_Types(t *testing.T) {
	type data struct {
		kind reflect.Kind
		in   any
		out  string
	}
	type unit struct {
		b *strings.Builder
		w gfmt.Writer
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
		{reflect.Float32, float32(math.MaxFloat32), strconv.FormatFloat(math.MaxFloat32, 'f', -1, 32)},
		{reflect.Float64, math.MaxFloat64, strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64)},
		// {reflect.Complex64, complex64(-2.71i), "(0-2.71i)"},   // not supported by JSON
		// {reflect.Complex128, complex128(-3.14i), "(0-3.14i)"}, // not supported by JSON
		// {reflect.Chan, ch, "chan<- bool"},                     // not supported by JSON
		// {reflect.Func, Test_Write_Types, "Test_Write_Types"},  // not supported by JSON
		{reflect.String, "str", "str"},
	}

	for _, tt := range tests {
		sPrettyJSON := &strings.Builder{}
		sJSON := &strings.Builder{}
		sTab := &strings.Builder{}
		sText := &strings.Builder{}

		prettyJSON := gfmt.NewJSON(sPrettyJSON, gfmt.WithPretty[gfmt.JSON](), gfmt.WithStyle[gfmt.JSON](styles.Get("native")))
		json := gfmt.NewJSON(sJSON)
		tab := gfmt.NewTab(sTab)
		text := gfmt.NewText(sText)

		us := []unit{
			{sPrettyJSON, prettyJSON},
			{sJSON, json},
			{sTab, tab},
			{sText, text},
		}

		for _, u := range us {
			t.Run(tt.kind.String()+"_"+reflect.TypeOf(u.w).Elem().Name(), func(t *testing.T) {
				want := tt.out

				postProc := func(s string) string { return s }
				if _, okText := u.w.(*gfmt.Text); okText {
					postProc = unJSON
				} else if _, okTab := u.w.(*gfmt.Tab); okTab {
					postProc = normalizeTable
				}

				_, err := u.w.Write(tt.in)
				require.NoError(t, err)
				require.Equal(t, want, u.b.String())

				if f, ok := u.w.(*gfmt.JSON); ok && f.Style != nil {
					// pretty JSON is too hard to verify, so we skip further tests
					return
				}

				if _, ok := u.w.(*gfmt.JSON); ok && tt.in == want {
					// JSON quotes strings, so the expected output needs to be quoted
					want = "\"" + want + "\""
				}

				_, isJSON := u.w.(*gfmt.JSON)

				// array
				u.b.Reset()
				_, err = u.w.Write([2]any{tt.in, tt.in})
				require.NoError(t, err)
				if isJSON && reflect.ValueOf(tt.in).CanFloat() {
					// TODO: JSON does not support arbitrary precision floats, so we skip this test
				} else if vf := reflect.ValueOf(tt.in); vf.CanFloat() {
					// f := strconv.FormatFloat(vf.Float(), 'f', -1, vf.Type().Bits())
					// require.Equal(t, postProc(fmt.Sprintf("[%s,%s]", f, f)), u.b.String())
				} else {
					require.Equal(t, postProc(fmt.Sprintf("[%s,%s]", want, want)), u.b.String())
				}

				// slice
				u.b.Reset()
				_, err = u.w.Write([]any{tt.in, tt.in})
				require.NoError(t, err)
				if isJSON && reflect.ValueOf(tt.in).CanFloat() {
					// TODO: JSON does not support arbitrary precision floats, so we skip this test
				} else if vf := reflect.ValueOf(tt.in); vf.CanFloat() && !isJSON {
					f := strconv.FormatFloat(vf.Float(), 'f', -1, vf.Type().Bits())
					require.Equal(t, postProc(fmt.Sprintf("[%s,%s]", f, f)), u.b.String())
				} else {
					require.Equal(t, postProc(fmt.Sprintf("[%s,%s]", want, want)), u.b.String())
				}

				// map
				if !isJSON {
					// JSON does not support arbitrary maps
					u.b.Reset()
					_, err = u.w.Write(map[any]any{tt.in: tt.in})
					require.NoError(t, err)
					require.Equal(t, postProc(fmt.Sprintf("%s:%s", want, want)), u.b.String())
				}
			})
		}
	}
}

// unJSON removes JSON-specific formatting such as [] and replaces comma with new line.
func unJSON(s string) string {
	return strings.Trim(strings.ReplaceAll(s, ",", "\n"), "[]")
}

// normalizeTable makes tabular output comparable by removing specific formatting.
func normalizeTable(s string) string {
	b := strings.Builder{}
	for _, p := range strings.Split(unJSON(s), ":") {
		b.WriteString(fmt.Sprintf("%-3s ", p))
	}
	if strings.Contains(s, ":") {
		// input was an array
		return b.String() + "\n"
	}
	// input was a single value
	return strings.TrimSpace(b.String())
}
