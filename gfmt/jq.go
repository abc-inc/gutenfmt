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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

type JQ struct {
	writer Writer
	Expr   string
}

func NewJQ(w Writer, expr string) *JQ {
	return &JQ{
		writer: w,
		Expr:   expr,
	}
}

func (w JQ) Write(i any) (int, error) {
	b := bytes.Buffer{}
	if err := evalJQ(i, &b, w.Expr); err != nil {
		return 0, err
	}
	var v any
	if err := json.Unmarshal(b.Bytes(), &v); err != nil {
		return 0, err
	}
	return w.writer.Write(v)
}

// evalJQ evaluates a jq expression against an input and write it to an output.
// Any top-level scalar values produced by the jq expression are written out as JSON scalars.
func evalJQ(val any, output io.Writer, expr string) error {
	query, err := gojq.Parse(expr)
	if err != nil {
		var e *gojq.ParseError
		if errors.As(err, &e) {
			str, line, column := getLineColumn(expr, e.Offset-len(e.Token))
			return fmt.Errorf(
				"failed to parse jq expression (line %d, column %d)\n    %s\n    %*c  %w",
				line, column, str, column, '^', err,
			)
		}
		return err
	}

	code, err := gojq.Compile(query, gojq.WithEnvironLoader(os.Environ))
	if err != nil {
		return err
	}

	iter := code.Run(val)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			var e *gojq.HaltError
			if errors.As(err, &e) && e.Value() == nil {
				break
			}
			return err
		}
		if text, ok := jsonScalarToString(v); ok {
			if _, err := fmt.Fprintln(output, text); err != nil {
				return err
			}
		} else {
			j, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintln(output, string(j)); err != nil {
				return err
			}
		}
	}

	return nil
}

func jsonScalarToString(input interface{}) (string, bool) {
	switch tt := input.(type) {
	case string:
		return `"` + tt + `"`, true
	case float64:
		if math.Trunc(tt) == tt {
			return strconv.FormatFloat(tt, 'f', 0, 64), true
		} else {
			return strconv.FormatFloat(tt, 'f', 2, 64), true
		}
	case nil:
		return "", true
	case bool:
		return fmt.Sprintf("%v", tt), true
	default:
		return "", false
	}
}

func getLineColumn(expr string, offset int) (string, int, int) {
	for line := 1; ; line++ {
		index := strings.Index(expr, "\n")
		if index < 0 {
			return expr, line, offset + 1
		}
		if index >= offset {
			return expr[:index], line, offset + 1
		}
		expr = expr[index+1:]
		offset -= index + 1
	}
}
