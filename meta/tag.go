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
	"strings"
)

// TagResolver reads the tags of struct fields to extract metadata.
//
// TagName can be any of, but not limited to
//
// - json: used by the encoding/json package, detailed at json.Marshal()
//
// - xml: used by the encoding/xml package, detailed at xml.Marshal()
//
// - yaml: used by the gopkg.in/yaml.v3 package, detailed at yaml.Marshal()
//
// - db: used by the github.com/jmoiron/sqlx package
//
// - orm: used by the github.com/beego/beego/orm package
//
// - gorm: used by the github.com/go-gorm/gorm package
type TagResolver struct {
	// TagName is the name of the tag to lookup.
	TagName string
}

// Lookup processes tags with a certain key in the fields' tag and uses the name, if defined.
// As a special case, if the field tag is "-", the field is omitted.
// Options like "omitempty" are ignored.
func (r TagResolver) Lookup(typ reflect.Type) (fs []field) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for idx := 0; idx < typ.NumField(); idx++ {
		sf := typ.Field(idx)
		if n := r.fieldName(sf); n != "" {
			fs = append(fs, field{sf.Name, n})
		}
	}
	return
}

// fieldName returns the name as set by the tag.
//
// As a special case, if the field tag is "-", an empty string is returned.
// Likewise, unexported fields result in an empty string, regardless of it's field tag.
//
// The name may be empty in order to specify options without overriding the default field name.
func (r TagResolver) fieldName(sf reflect.StructField) string {
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
	tag := sf.Tag.Get(r.TagName)
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

// parseTag splits a struct field's tag into its name and comma-separated options.
func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
