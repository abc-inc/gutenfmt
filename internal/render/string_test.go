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

package render_test

import (
	"testing"

	. "github.com/abc-inc/gutenfmt/internal/render"
	. "github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	Equal(t, "", ToString(nil))
	Equal(t, "x", ToString("x"))
	Equal(t, "true", ToString(true))

	Equal(t, "", ToString([]int{}))
	Equal(t, "", ToString([0]int{}))
	Equal(t, "1", ToString([]int{1}))
	Equal(t, "1 2", ToString(&[]int{1, 2}))

	Equal(t, "chan int", ToString(make(chan int)))

	Regexp(t, `/internal/render_test\.TestToString$`, ToString(TestToString))
}
