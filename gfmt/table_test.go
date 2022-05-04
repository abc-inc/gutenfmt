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
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/abc-inc/gutenfmt/formatter"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func TestTab_Write(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"nil", nil, ""},
		{"bool", false, "false"},
		{"int", -42, "-42"},
		{"string", "∮∯∰", "∮∯∰"},
		{"empty_array", [0]string{}, "^$"},
		{"int_slice", []int{1, 2, 3}, "1\n2\n3"},
		{"struct", NewUser("John", "Doe"), "username John Doe \nemail    john.doe@local"},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")}, "^\ntrue\n-42\na\nf l <f.l@local>$"},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, "(a a 1   \n:   :   )|(:   :   \na a 1   )"},
		{"map_slice", []map[string]interface{}{{"a": 1, "b": 2}, {"b": 3, "a": 4}},
			"(a   b   \n1   2   \n4   3)|(b   a   \n2   1   \n3   4)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewTab(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestTab_WriteAllTypes(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewTab(b).Write(allTypes)
	NoError(t, err)

	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")

	Equal(t, `DefName        DefName
OmitEmpty      OmitEmpty
custom         CustOmitEmpty
EmptyOmitEmpty`+`
Bool           true
Int            -4
Int8           -8
Int16          -16
Int32          -32
Int64          -64
Uint           4
Uint8          8
Uint16         16
Uint32         32
Uint64         64
Uintptr        128
Float32        3.4028235e+38
Float64        1.7976931348623157e+308
Complex64      (0-2.71i)
Complex128     (0-3.14i)
Array          a b`+`
Chan           chan<- int`+`
Func           github.com/abc-inc/gutenfmt/gfmt_test.NewUser
Interface`+`
Map            map[]
Ptr            f l <f.l@local>
Slice          a b`+`
String`+`
Struct         f l <f.l@local>
StructSlice    af al <af.al@local> bf bl <bf.bl@local>`,
		s)
}

func TestTab_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewTab(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Equal(t, "username John Doe \nemail    john.doe@local", b.String())
}

func TestTab_WriteStructSlice(t *testing.T) {
	type data struct {
		A interface{}
		B string
		C interface{}
	}

	b := &strings.Builder{}
	_, err := NewTab(b).Write([]data{{A: 'a', B: "b", C: true}, {A: "d", B: "e", C: "f"}})
	NoError(t, err)
	Equal(t, "A   B   C    \n97  b   true \nd   e   f", b.String())
}

func TestTab_WriteStructPtrSlice(t *testing.T) {
	type data struct {
		A, B string
	}

	b := &strings.Builder{}
	_, err := NewTab(b).Write([]*data{{A: "1", B: "2"}, {A: "3", B: "4"}})
	NoError(t, err)
	Equal(t, "A   B   \n1   2   \n3   4", b.String())
}

func TestTab_WriteMap(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewTab(b).Write(map[string]interface{}{"a": 'a', "b": "b", "c": true})
	NoError(t, err)
	Regexp(t, "([a-c]   (97|b|true)[ ]{1,4}\n){3}", b.String())
}

func TestTab_WriteMapSliceCustom(t *testing.T) {
	mss := []map[string]string{{"a": "w", "b": "x"}, {"a": "y", "b": "z"}}
	msi := []map[string]int{{"c": 1, "d": 2}, {"c": 3, "d": 4}}

	b := &strings.Builder{}
	w := NewTab(b)
	w.Formatter.SetFormatter(reflect.TypeOf(mss).String(),
		formatter.AsTab(formatter.FromMapSliceKeys("\t", "\t\n", reflect.ValueOf("a"))))

	_, err := w.Write(mss)
	NoError(t, err)
	Equal(t, "a   \nw   \ny", b.String())

	b.Reset()
	_, err = w.Write(msi)
	NoError(t, err)
	Regexp(t, "(c   d   \n1   2   \n3   4)|(d   c   \n2   1   \n4   3)", b.String())
}
