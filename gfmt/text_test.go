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
	"strings"
	"testing"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func TestText_Write(t *testing.T) {
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
		{"struct_stringer", NewUser("John", "Doe"), "John Doe <john.doe@local>"},
		{"struct", NewOrg("Enterprise"), "Name:Enterprise\nTeams:"},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")}, "^\ntrue\n-42\na\nf l <f.l@local>$"},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, "^(:::\na a:1)|(a a:1\n:::)$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewText(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestText_WriteAllTypes(t *testing.T) {
	b := &strings.Builder{}
	w := NewText(b)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(allTypes).Name(), func(i interface{}) (string, error) {
		return i.(AllTypes).DefName, nil
	})
	_, err := w.Write(allTypes)
	NoError(t, err)
	Equal(t, "DefName", b.String())
}

func TestText_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewText(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Equal(t, "John Doe <john.doe@local>", b.String())
}

func TestText_WriteStructSlice(t *testing.T) {
	type data struct {
		A interface{}
		B string
		C interface{}
	}

	b := &strings.Builder{}
	_, err := NewText(b).Write([]data{{A: 'a', B: "b", C: true}, {A: "d", B: "e", C: "f"}})
	NoError(t, err)
	Equal(t, "A:B:C\n97:b:true\nd:e:f", b.String())
}

func TestText_WriteStructPtrSlice(t *testing.T) {
	type data struct {
		A, B string
	}

	b := &strings.Builder{}
	_, err := NewText(b).Write([]*data{{A: "1", B: "2"}, {A: "3", B: "4"}})
	NoError(t, err)
	Equal(t, "A:B\n1:2\n3:4", b.String())
}

func TestText_WriteMap(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewText(b).Write(map[string]interface{}{"a": 'a', "b": "b", "c": true})
	NoError(t, err)
	Regexp(t, "([a-c]:(97|b|true)\n){2}([a-c]:(97|b|true))", b.String())
}

func TestText_WriteMapSlice(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewText(b).Write([]map[string]interface{}{{"a": 'a', "b": true}, {"a": "c", "b": "d"}})
	NoError(t, err)
	Regexp(t, "^(a:b\n97:true\nc:d)|(b:a\ntrue:97\nd:c)", b.String())
}
