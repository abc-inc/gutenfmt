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

package meta_test

import (
	"reflect"
	"testing"

	"github.com/abc-inc/gutenfmt/meta"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name     string `json:"username" yaml:"Username"`
	Mail     string `json:"email" yaml:"E-Mail"`
	Password string `json:"-" yaml:"-"`
	KeyPair  `json:"keys"`
}

type KeyPair struct {
	Pub  []byte `json:"pub"`
	priv []byte //nolint:structcheck,unused
}

func TestTagResolver_lookup(t *testing.T) {
	fs := meta.TagResolver{TagName: "yaml"}.Lookup(reflect.TypeOf(User{}))
	assert.Equal(t, 2, len(fs))
	assert.Equal(t, "Username", fs[0].Name)
	assert.Equal(t, "Name", fs[0].Field)
	assert.Equal(t, "E-Mail", fs[1].Name)
	assert.Equal(t, "Mail", fs[1].Field)
}
