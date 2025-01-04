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

package formatter

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompFormatter_Format(t *testing.T) {
	f := NewComp()
	_, err := f.Format("x")
	require.Error(t, err)

	f.SetFormatterFunc("string", func(i any) (string, error) {
		return strings.ToUpper(i.(string)), nil
	})
	s, err := f.Format("x")
	require.NoError(t, err)
	require.Equal(t, "X", s)

	f.SetFormatter("string", NoopFormatter())
	s, err = f.Format("x")
	require.NoError(t, err)
	require.Equal(t, "", s)
}

func TestTypeName(t *testing.T) {
	type User struct {
		Name     string `json:"username"`
		Mail     string `json:"email"`
		Password string `json:"-"`
	}

	s0 := struct{}{}
	s1 := struct{ string }{""}
	s2 := struct{ s string }{""}

	require.Equal(t, "int", typeName(reflect.TypeOf(math.MaxInt16)))
	require.Equal(t, "int", typeName(reflect.TypeOf(math.MaxInt32)))
	require.Equal(t, "int", typeName(reflect.TypeOf(math.MaxInt64)))
	require.Equal(t, "string", typeName(reflect.TypeOf("")))

	require.Equal(t, "User", typeName(reflect.TypeOf(User{})))
	require.Equal(t, "Decoder", typeName(reflect.TypeOf(json.Decoder{})))

	require.Equal(t, "struct {}", typeName(reflect.TypeOf(s0)))
	require.Equal(t, "struct { string }", typeName(reflect.TypeOf(s1)))
	require.Equal(t, "struct { s string }", typeName(reflect.TypeOf(s2)))

	require.Equal(t, "map[interface {}]bool", typeName(reflect.TypeOf(map[any]bool{})))
}
