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
	"time"

	. "github.com/abc-inc/gutenfmt/renderer"
	. "github.com/stretchr/testify/require"
)

func TestFromStruct(t *testing.T) {
	u := NewUser("Jane", "Doe")
	r := FromStruct(": ", ", ", reflect.TypeOf(u))
	s, _ := r.Render(u)
	Equal(t, "username: Jane Doe, email: jane.doe@local", s)
}

func TestFromStruct_Anonymous(t *testing.T) {
	type Model struct {
		Company string `json:"Manufacturer,omitempty"`
		Model   string `json:"Model,omitempty"`
		rating  int
	}

	m := &Model{"Company", "Awesome", 5}

	// The embedded struct pointer will be handled properly.
	// Although, it will be rendered with its default string representation.
	car := struct {
		*Model      `json:"Type,omitempty"`
		Mileage     int `json:"Miles"`
		serviceDate time.Time
	}{m, 1337, time.Now()}

	r := FromStruct(": ", ", ", reflect.TypeOf(car))
	s, err := r.Render(car)
	NoError(t, err)
	Equal(t, "Type: {Company Awesome 5}, Miles: 1337", s)
}

func TestFromStructSlice(t *testing.T) {
	u := NewUser("Jane", "Doe")
	r := FromStructSlice(" | ", "\n", reflect.TypeOf(u))
	s, _ := r.Render([]*User{u})
	Equal(t, "username | email\nJane Doe | jane.doe@local", s)
}
