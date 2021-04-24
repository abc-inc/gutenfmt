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
	"strings"
	"testing"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func TestNewAutoJSON(t *testing.T) {
	f := NewAutoJSON(&strings.Builder{})
	IsType(t, &JSON{}, f)
}

func TestJSON_Write(t *testing.T) {
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
		{"int_slice", []int{1, 2, 3}, `\[1,2,3\]`},
		{"struct", NewUser("John", "Doe"), `"email":"john.doe@local"`},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")}, `,true,-42,"a",{"username":"f l",`},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, `{":":":"."a a":1}|{"a a":1,":":":"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewJSON(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestJSON_WriteJSONTypes(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewJSON(b).Write(jsonTypes)
	NoError(t, err)
	Contains(t, b.String(), "\"Ptr\":{\"username\":\"")
}

func TestJSON_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	n, err := NewJSON(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Equal(t, 48, n)
	Contains(t, b.String(), "\"username\":\"John Doe\",\"email\":\"john.doe@local\"")
}
