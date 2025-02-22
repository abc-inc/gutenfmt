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

	"github.com/abc-inc/gutenfmt/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	assert.Equal(t, "", render.ToString(nil))
	assert.Equal(t, "x", render.ToString("x"))
	assert.Equal(t, "true", render.ToString(true))

	assert.Equal(t, "", render.ToString([]int{}))
	assert.Equal(t, "", render.ToString([0]int{}))
	assert.Equal(t, "1", render.ToString([]int{1}))
	assert.Equal(t, "1 2", render.ToString(&[]int{1, 2}))

	assert.Equal(t, "chan int", render.ToString(make(chan int)))

	assert.Regexp(t, `/internal/render_test\.TestToString$`, render.ToString(TestToString))
}
