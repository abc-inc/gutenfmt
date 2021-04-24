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

func jsonMetadata(typ reflect.Type) ([]string, []string) {
	var fns []string
	var pns []string
	for idx := 0; idx < typ.NumField(); idx++ {
		f := typ.Field(idx)
		if n := jsonPropName(f); n != "" {
			fns = append(fns, f.Name)
			pns = append(pns, n)
		}
	}
	return fns, pns
}

func jsonPropName(f reflect.StructField) string {
	if n, ok := f.Tag.Lookup("json"); strings.HasPrefix(n, "-") || f.PkgPath != "" {
		return ""
	} else if n = strings.SplitN(n, ",", 2)[0]; ok && n != "" {
		return n
	}
	return f.Name
}
