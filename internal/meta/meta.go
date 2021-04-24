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
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// IsContainerType returns true if a type is kind of a "container".
//
// Note that "container" is not an official classification.
// In this context, it represents a subset of composite types, which can hold
// a certain amount of elements that can be accessed in arbitrary order,
// namely array, struct, slice and map.
func IsContainerType(k reflect.Kind) bool {
	return k == reflect.Struct || k == reflect.Slice || k == reflect.Map || k == reflect.Array
}

// TypeName returns the type's name.
// If the type cannot be determined e.g., []interface{}, a string representation
// is returned instead.
//
// Note that a string representation is not necessarily unique among types.
func TypeName(typ reflect.Type) string {
	n := typ.Name()
	if n == "" {
		n = typ.String()
	}
	return n
}

// TODO: delete StrVal or strFormat
// TODO: rename to scalar??? GH?
// TODO: use always sprint()?

func StrVal(i interface{}) string {
	return strFormat(reflect.ValueOf(i))
}

func strFormat(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Uintptr:
		return fmt.Sprint(reflect.Indirect(v).Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprint(v.Complex())
	case reflect.Array:
		return strings.Trim(fmt.Sprint(v.Interface()), "[]")
	case reflect.Chan:
		return ""
	case reflect.Func:
		return funcName(v)
	case reflect.Interface:
		return fmt.Sprint(v.Interface())
	case reflect.Map:
		return fmt.Sprint(v.Interface())
	case reflect.Ptr:
		return strFormat(reflect.Indirect(v))
	case reflect.Slice:
		return strings.Trim(fmt.Sprint(v.Interface()), "[]")
	case reflect.String:
		return v.String()
	case reflect.Struct:
		return fmt.Sprint(v.Interface())
	// case reflect.UnsafePointer:
	default:
		panic(fmt.Sprintf("Cannot convert %s (%v)", reflect.TypeOf(v).Name(), v.Type()))
	}
}

func funcName(f reflect.Value) string {
	return runtime.FuncForPC(f.Pointer()).Name()
}
