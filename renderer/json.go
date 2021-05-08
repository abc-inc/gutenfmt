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

package renderer

import (
	"reflect"
	"strings"
)

// jsonMetadata returns struct fields and the JSON property names.
func jsonMetadata(typ reflect.Type) ([]string, []string) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var fns []string
	var pns []string
	for idx := 0; idx < typ.NumField(); idx++ {
		sf := typ.Field(idx)
		if n := jsonPropName(sf); n != "" {
			fns = append(fns, sf.Name)
			pns = append(pns, n)
		}
	}
	return fns, pns
}

// jsonPropName returns the JSON property name.
//
// As a special case, if the field tag is "-", an empty string is returned.
// Likewise, unexported fields result in an empty string, regardless of it's field tag.
//
// The name may be empty in order to specify options without overriding the default field name.
func jsonPropName(sf reflect.StructField) string {
	isUnexported := sf.PkgPath != ""
	if sf.Anonymous {
		t := sf.Type
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if isUnexported && t.Kind() != reflect.Struct {
			// Ignore embedded fields of unexported non-struct types.
			return ""
		}
		// Do not ignore embedded fields of unexported struct types
		// since they may have exported fields.
	} else if isUnexported {
		// Ignore unexported non-embedded fields.
		return ""
	}
	tag := sf.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	name := parseTag(tag)
	ft := sf.Type
	if ft.Name() == "" && ft.Kind() == reflect.Ptr {
		// Follow pointer.
		ft = ft.Elem()
	}

	if name == "" && (!sf.Anonymous || ft.Kind() != reflect.Struct) {
		name = sf.Name
	}

	return name
}

// parseTag splits a struct field's json tag into its name and comma-separated options.
func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
