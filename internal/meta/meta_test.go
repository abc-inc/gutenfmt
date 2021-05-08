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

package meta

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestTypeName(t *testing.T) {
	type User struct {
		Name     string `json:"username"`
		Mail     string `json:"email"`
		Password string `json:"-"`
	}

	s0 := struct{}{}
	s1 := struct{ string }{""}
	s2 := struct{ s string }{""}

	Equal(t, "int", TypeName(reflect.TypeOf(math.MaxInt16)))
	Equal(t, "int", TypeName(reflect.TypeOf(math.MaxInt32)))
	Equal(t, "int", TypeName(reflect.TypeOf(math.MaxInt64)))
	Equal(t, "string", TypeName(reflect.TypeOf("")))

	Equal(t, "User", TypeName(reflect.TypeOf(User{})))
	Equal(t, "Decoder", TypeName(reflect.TypeOf(json.Decoder{})))

	Equal(t, "struct {}", TypeName(reflect.TypeOf(s0)))
	Equal(t, "struct { string }", TypeName(reflect.TypeOf(s1)))
	Equal(t, "struct { s string }", TypeName(reflect.TypeOf(s2)))

	Equal(t, "map[interface {}]bool", TypeName(reflect.TypeOf(map[interface{}]bool{})))
}

func TestToString(t *testing.T) {
	Equal(t, "", ToString(nil))
	Equal(t, "x", ToString("x"))
	Equal(t, "true", ToString(true))

	Equal(t, "", ToString([]int{}))
	Equal(t, "", ToString([0]int{}))
	Equal(t, "1", ToString([]int{1}))
	Equal(t, "1 2", ToString(&[]int{1, 2}))

	Equal(t, "chan int", ToString(make(chan int)))

	Regexp(t, `/internal/meta\.TestToString$`, ToString(TestToString))
}
