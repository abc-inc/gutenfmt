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
	"strings"
	"testing"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func TestYAML_Write(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"nil", nil, ""},
		{"bool", false, "false"},
		{"int", -42, "-42"},
		{"string", "∮∯∰", "∮∯∰"},
		{"empty_array", [0]string{}, `^\[\]$`},
		{"int_slice", []int{1, 2, 3}, `- 1\n- 2\n- 3`},
		{"struct", NewUser("John", "Doe"), `^Username: John Doe\nE-Mail: john.doe@local$`},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")},
			`- true\n- -42\n- a\n- Username: f l\n  E-Mail: f.l@local`},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, `':': ':'\na a: 1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewYAML(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestYAML_WriteJSONTypes(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewYAML(b).Write(jsonTypes)
	NoError(t, err)
	Contains(t, b.String(), "structslice:\n  - Username: af al\n    E-Mail: af.al@local\n  - Username: bf bl\n")
}

func TestYAML_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	n, err := NewYAML(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Equal(t, 41, n)
	Equal(t, b.String(), "Username: John Doe\nE-Mail: john.doe@local")
}
