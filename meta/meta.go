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
	"reflect"
)

// Field represents a single field found in a struct.
type Field struct {
	Field string
	Name  string
}

// Resolver returns a list of fields that should be recognized for the given type.
type Resolver func(typ reflect.Type) []Field

// Resolve holds the default Resolver.
var Resolve Resolver = TagResolver{"json"}.Lookup
