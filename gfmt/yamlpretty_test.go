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

func TestPrettyYAML_Write(t *testing.T) {
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
		{"int_slice", []int{1, 2, 3}, `- \x1b\[\d+m1\x1b\[0m.*\n\x1b\[0m- \x1b\[\d+m`},
		{"struct", NewUser("John", "Doe"), `E-Mail\x1b\[0m:\x1b\[1m\x1b\[30m \x1b\[0mjohn.doe@local`},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")}, `-.+true.+\n.+-.+-\x1b\[\d+m42`},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, `':'.+:.+':'.+`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewPrettyYAML(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestPrettyYAML_WriteJSONTypes(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewPrettyYAML(b).Write(jsonTypes)
	NoError(t, err)
	Contains(t, b.String(), "mptr\x1b[0m:")
}

func TestPrettyYAML_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewPrettyYAML(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Regexp(t, `\x1b\[\d+m.*\x1b\[\d+mUsername\x1b\[0m:.*\x1b\[0mJohn Doe\x1b`, b.String())
}
