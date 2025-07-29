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

package render

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
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
func ToString(i any) string {
	if i == nil {
		return ""
	}

	v := reflect.ValueOf(i)
	switch {
	case v.CanInt():
		return strconv.FormatInt(v.Int(), 10)
	case v.CanUint():
		return strconv.FormatUint(v.Uint(), 10)
	default:
		// check the type's kind
	}

	typ := v.Type()
	switch typ.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Array, reflect.Slice:
		s := fmt.Sprint(i)
		return s[1 : len(s)-1]
	case reflect.Chan:
		return typ.String()
	case reflect.Complex64:
		return strconv.FormatComplex(complex128(i.(complex64)), 'f', -1, 64)
	case reflect.Complex128:
		return strconv.FormatComplex(i.(complex128), 'f', -1, 128)
	case reflect.Float32:
		// float32 requires special handling, because of IEEE 754.
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case reflect.Func:
		return funcName(reflect.ValueOf(i))
	case reflect.Ptr:
		if v.IsValid() {
			e := reflect.Indirect(v)
			if e.IsValid() {
				return ToString(e.Interface())
			}
		}
		return ""
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
