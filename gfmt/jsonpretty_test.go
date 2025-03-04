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

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/stretchr/testify/require"
)

func TestPrettyJSON_Write(t *testing.T) {
	tests := []struct {
		name string
		arg  any
		want string
	}{
		{"nil", nil, ""},
		{"bool", false, "false"},
		{"int", -42, "-42"},
		{"string", "∮∯∰", "∮∯∰"},
		{"empty_array", [0]string{}, `^\x1b\[\d+m\[\]\x1b\[\d+m$`},
		{"int_slice", []int{1, 2, 3}, `\x1b\[\d+m\[\x1b\[0m.*\n.*1.*,.*\n.*2.*,.*\n.*3.*\n\x1b\[\d+m\]`},
		{"struct", NewUser("John", "Doe"), `"email".+:.* .+"john.doe@local"`},
		{"mixed_array", []any{[0]string{}, true, -42, "a", NewUser("f", "l")}, `.*true.+,.*\n.*-42.+,.*(\n.+)+"f l"`},
		{"map", map[string]any{"a a": 1, ":": ":"}, `":".+:.+":".+`},
	}

	for _, tt := range tests {
		b := strings.Builder{}
		w := gfmt.NewJSON(&b, gfmt.WithPretty[gfmt.JSON](), gfmt.WithStyle[gfmt.JSON](styles.Get("native")))
		t.Run(tt.name, func(t *testing.T) {
			b.Reset()
			_, err := w.Write(tt.arg)
			require.NoError(t, err)
			require.Regexp(t, tt.want, b.String())
		})
	}
}

func TestPrettyJSON_WriteJSONTypes(t *testing.T) {
	b := strings.Builder{}
	_, err := gfmt.NewJSON(&b, gfmt.WithStyle[gfmt.JSON](styles.Get("native"))).Write(jsonTypes)
	require.NoError(t, err)
	require.Contains(t, b.String(), "\"Ptr\"\x1b[0m\x1b[37m:")
}

func TestPrettyJSON_WriteStruct(t *testing.T) {
	b := strings.Builder{}
	_, err := gfmt.NewJSON(&b, gfmt.WithStyle[gfmt.JSON](styles.Get("native"))).Write(NewUser("John", "Doe"))
	require.NoError(t, err)
	require.Regexp(t, `\x1b\[\d+m.*\x1b\[\d+m"John Doe"`, b.String())
}
