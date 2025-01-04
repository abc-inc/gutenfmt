// Copyright 2025 The gutenfmt authors
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

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/alecthomas/chroma/styles"
	"github.com/stretchr/testify/require"
)

func TestJQWriter_Write(t *testing.T) {
	exp := heredoc.Doc(`
	c:
	  - c
	  - c
	d:
	  e: 4`)

	b := strings.Builder{}
	w := gfmt.NewJQ(gfmt.NewYAML(&b), ".b")
	_, err := w.Write(map[string]any{"a": 1, "b": map[string]any{"c": []string{"c", "c"}, "d": map[string]any{"e": 4}}})
	require.NoError(t, err)
	require.Equal(t, exp, b.String())
}

func TestJQWriter_WritePretty(t *testing.T) {
	b := strings.Builder{}
	_, _ = gfmt.NewJQ(gfmt.NewJSON(&b, gfmt.WithPretty(), gfmt.WithStyle(styles.Fallback)), ".").Write([]string{"A"})
	require.Regexp(t, `^\[\n  `, b.String())

	b = strings.Builder{}
	_, _ = gfmt.NewJQ(gfmt.NewJSON(&b, gfmt.WithPretty(), gfmt.WithStyle(styles.Fallback)), ".nestedNumber").Write(map[string]any{"nestedNumber": "5"})
	require.Regexp(t, `5`, b.String())

	b = strings.Builder{}
	_, _ = gfmt.NewJQ(gfmt.NewJSON(&b, gfmt.WithPretty(), gfmt.WithStyle(styles.Fallback)), ".").Write(`a5`)
	require.Regexp(t, `a5`, b.String())

	b = strings.Builder{}
	_, _ = gfmt.NewJQ(gfmt.WrapIOWriter(&b), ".nestedString").Write(map[string]any{"nestedString": "a5"})
	require.Regexp(t, `a5`, b.String())
}

func TestJQWriter_WriteAny(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expr     string
		expected string
	}{
		{name: "integer", input: 42, expr: ".", expected: "42"},
		{name: "string", input: `hello`, expr: ".", expected: "hello"},
		{name: "array", input: []any{1, 2, 3}, expr: ".[1]", expected: "2"},
		{name: "nested", input: map[string]any{"a": map[string]any{"b": 2}}, expr: ".a.b", expected: "2"},
		{name: "invalid", input: map[string]any{"a": 1}, expr: ".x", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			w := gfmt.NewJQ(gfmt.WrapIOWriter(b), tt.expr)
			_, err := w.Write(tt.input)
			if tt.expected == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, b.String())
			}
		})
	}
}
