// Copyright 2025 The gutenfmt authors
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

package gfmt

import (
	"reflect"

	"github.com/jmespath/go-jmespath"
)

type JMESPath struct {
	writer Writer
	Expr   *jmespath.JMESPath
	Indent string
	Color  bool
}

func NewJMESPath(w Writer, expr string) *JMESPath {
	return &JMESPath{
		writer: w,
		Expr:   jmespath.MustCompile(expr),
		Indent: "",
		Color:  false,
	}
}

func (w JMESPath) Write(i any) (int, error) {
	typ := reflect.TypeOf(i)
	k := typ.Kind()

	if k == reflect.Pointer {
		return w.Write(reflect.ValueOf(i).Elem().Interface())
	}
	if isPrimitive(typ) {
		return w.writer.Write(i)
	}
	if k == reflect.Array || k == reflect.Slice {
		if isPrimitive(typ.Elem()) {
			return w.writer.Write(i)
		}
	}

	if k == reflect.Map || k == reflect.Struct || k == reflect.Array || k == reflect.Slice {
		v, err := w.Expr.Search(i)
		if err != nil {
			return 0, err
		}
		return w.writer.Write(v)
	}
	panic("unsupported type: " + typ.String())
}

func isPrimitive(typ reflect.Type) bool {
	switch typ.Kind() { //nolint:exhaustive
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}
