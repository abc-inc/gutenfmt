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

package renderer_test

import (
	"reflect"
	"testing"

	. "github.com/abc-inc/gutenfmt/renderer"
	. "github.com/stretchr/testify/require"
)

func TestAsTab(t *testing.T) {
	ks := []reflect.Value{reflect.ValueOf("a"), reflect.ValueOf("long_key")}
	m := map[string]int{"a": 1, "long_key": 2}

	r := FromMapKeys("\t", "\t\n", ks...)
	s, _ := r.Render(m)
	Equal(t, "a\t1\t\nlong_key\t2\t\n", s)

	s, _ = AsTab(r).Render(m)
	Equal(t, "a        1   \nlong_key 2   \n", s)
}
