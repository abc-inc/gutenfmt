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

package formatter_test

import (
	"testing"
	"text/template"

	"github.com/abc-inc/gutenfmt/formatter"
	"github.com/stretchr/testify/require"
)

func TestFromTemplate(t *testing.T) {
	text := "mailto:{{.Mail}}\nDear {{.Name}}"
	f := formatter.FromTemplate(template.Must(template.New("letter").Parse(text)))
	s, _ := f.Format(map[string]string{"Name": "Jane Doe", "Mail": "jane.doe@local"})
	require.Equal(t, "mailto:jane.doe@local\nDear Jane Doe", s)
}
