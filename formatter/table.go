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

package formatter

import (
	"strings"
	"text/tabwriter"
)

// AsTab returns a new Formatter that translates tabbed columns into properly aligned text.
func AsTab(f Formatter) Formatter {
	return Func(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		tw := tabwriter.NewWriter(b, 4, 4, 1, ' ', 0)
		_, err := FormatTab(tw, f, i)
		return b.String(), err
	})
}

// FormatTab formats the given input and translates tabbed columns into properly aligned text.
func FormatTab(tw *tabwriter.Writer, f Formatter, i interface{}) (n int, err error) {
	s, err := f.Format(i)
	if err != nil {
		return 0, err
	}
	if n, err = tw.Write([]byte(s)); err != nil {
		return n, err
	}
	return n, tw.Flush()
}
