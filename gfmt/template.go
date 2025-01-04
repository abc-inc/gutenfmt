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
	"io"
	"text/template"
)

type Tmpl struct {
	writer io.Writer
	tmpl   *template.Template
}

func NewTemplatePattern(output io.Writer, p string) Writer {
	return NewTemplate(output, template.Must(template.New("output").Parse(p)))
}

func NewTemplate(w io.Writer, tmpl *template.Template) Writer {
	return &Tmpl{
		writer: w,
		tmpl:   tmpl,
	}
}

func (w Tmpl) Write(a any) (int, error) {
	return 0, w.tmpl.Execute(w.writer, a)
}
