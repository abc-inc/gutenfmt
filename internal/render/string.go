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

package render

import (
	"fmt"
	"reflect"
	"runtime"
)

// ToString returns a human-readable string representation of i.
//
// The following types are treated specially:
//
// - arrays, slices: surrounding [] are removed
//
// - functions: best effort, since functions could be anonymous
//
// - pointers: dereferenced before determining their string representation
//
// - others: formatted using the default formats like fmt.Sprint
func ToString(i interface{}) string {
	if i == nil {
		return ""
	}
	typ := reflect.TypeOf(i)
	switch typ.Kind() { //nolint:exhaustive
	case reflect.Array, reflect.Slice:
		s := fmt.Sprint(i)
		return s[1 : len(s)-1]
	case reflect.Chan:
		return typ.String()
	case reflect.Func:
		return funcName(reflect.ValueOf(i))
	case reflect.Ptr:
		return ToString(reflect.Indirect(reflect.ValueOf(i)).Interface())
	case reflect.String:
		return i.(string)
	default:
		return fmt.Sprint(i)
	}
}

// funcName returns the name of the function f points to.
func funcName(f reflect.Value) string {
	return runtime.FuncForPC(f.Pointer()).Name()
}
