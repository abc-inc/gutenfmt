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

package formatter_test

import (
	"reflect"
	"testing"

	. "github.com/abc-inc/gutenfmt/formatter"
	. "github.com/stretchr/testify/require"
)

var m = map[string]bool{"y": true, "n": false}

func TestFromMap(t *testing.T) {
	f := FromMap("\t", "\t\n")
	s, _ := f.Format(m)
	Regexp(t, "(y\ttrue\t\nn\tfalse)|(n\tfalse\t\ny\ttrue)", s)
}

func TestFromMapKeys_duplicate(t *testing.T) {
	f := FromMapKeys("\t", "\t\n", reflect.ValueOf("y"), reflect.ValueOf("y"))
	s, _ := f.Format(m)
	Equal(t, "y\ttrue\t\ny\ttrue\t\n", s)

	f = FromMapSliceKeys("\t", "\t\n", reflect.ValueOf("y"), reflect.ValueOf("y"))
	s, _ = f.Format([]map[string]bool{m, m})
	Equal(t, "y\ty\t\ntrue\ttrue\t\ntrue\ttrue", s)
}

func TestFromMapKeys_invalidKey(t *testing.T) {
	f := FromMapKeys("\t", "\t\n", reflect.ValueOf("t"))
	s, _ := f.Format(m)
	Equal(t, "t\t\t\n", s)

	f = FromMapSliceKeys("\t", "\t\n", reflect.ValueOf("t"))
	s, _ = f.Format([]map[string]bool{m, m})
	Equal(t, "t\t\n\t\n", s)
}

func TestFromMapSlice(t *testing.T) {
	f := FromMapSlice("\t", "\t\n")
	s, _ := f.Format([]map[string]bool{m})
	Regexp(t, "(y\tn\t\ntrue\tfalse)|(n\ty\t\nfalse\ttrue)", s)
}
