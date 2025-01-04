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
	"fmt"
	"strings"
	"testing"

	"github.com/abc-inc/gutenfmt/gfmt"
	"github.com/stretchr/testify/require"
)

func TestJMESPathWriter_Write(t *testing.T) {
	b := &strings.Builder{}
	w := gfmt.NewJMESPath(gfmt.WrapIOWriter(b), "b")
	n, err := w.Write(&map[string]any{"a": 1, "b": 2})
	require.Equal(t, 1, n)
	require.NoError(t, err)
	require.Equal(t, "2", b.String())
}

func TestJMESPathWriter_WritePrimitive(t *testing.T) {
	tests := []struct {
		input    any
		expr     string
		expected string
	}{
		{
			input:    42,
			expr:     "*",
			expected: "42",
		},
		{
			input:    "string",
			expr:     "*",
			expected: "string",
		},
		{
			input:    []byte("bytes"),
			expr:     "*",
			expected: fmt.Sprint([]byte("bytes")),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.input), func(t *testing.T) {
			b := &strings.Builder{}
			w := gfmt.NewJMESPath(gfmt.WrapIOWriter(b), tt.expr)
			_, err := w.Write(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, b.String())
		})
	}
}
