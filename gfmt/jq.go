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
	"reflect"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

type Arg struct {
	Key    string
	Val    string
	String bool
}

func NewArg(kv string, strLiteral bool) (*Arg, error) {
	if k, v, ok := strings.Cut(kv, "="); ok {
		return &Arg{Key: k, Val: v, String: strLiteral}, nil
	}
	return nil, fmt.Errorf("invalid argument: %s", kv)
}

type JQ struct {
	writer Writer
	Expr   string
	Args   []Arg
}

func NewJQ(w Writer, expr string, args ...Arg) *JQ {
	return &JQ{
		writer: w,
		Expr:   expr,
		Args:   args,
	}
}

func (w JQ) Write(i any) (int, error) {
	b := bytes.Buffer{}
	if err := evalJQ(i, &b, w.Expr, w.Args...); err != nil {
		return 0, err
	}
	// If the output is NO json, e.g., a literal string or null, write it as is.
	if !json.Valid(b.Bytes()) || b.String() == "null\n" {
		return w.writer.Write(strings.TrimSuffix(b.String(), "\n"))
	}

	// Otherwise, create a new data structure and let the other writer handle it.
	var v any
	if err := json.Unmarshal(b.Bytes(), &v); err != nil {
		return 0, err
	}
	return w.writer.Write(v)
}

// evalJQ evaluates a jq expression against an input and write it to an output.
// Any top-level scalar values produced by the jq expression are written out as JSON scalars.
func evalJQ(v any, w io.Writer, expr string, args ...Arg) error {
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

	vars, vals, err := parseArgs(args...)
	if err != nil {
		return err
	}

	code, err := gojq.Compile(query,
		gojq.WithEnvironLoader(os.Environ),
		gojq.WithVariables(vars),
		gojq.WithFunction("raw", 0, 0, rawFunc),
	)
	if err != nil {
		return err
	}

	// gojq panic upon errors involving third-party types like lists of a custom struct.
	// As a workaround, we serialize the input to JSON and then deserialize it back.
	var buf []byte
	if buf, err = json.Marshal(v); err != nil {
		return err
	}
	if err = json.Unmarshal(buf, &v); err != nil {
		return err
	}

	iter := code.Run(v, vals...)
	for {
		val, hasNext := iter.Next()
		if !hasNext {
			break
		}
		if vErr, isErr := val.(error); isErr {
			var e *gojq.HaltError
			if errors.As(vErr, &e) && e.Value() == nil {
				break
			}
			return vErr
		}

		var j []byte
		if j, err = json.Marshal(val); err != nil {
			return err
		}
		if _, err = fmt.Fprintln(w, string(j)); err != nil {
			return err
		}
	}

	return nil
}

func parseArgs(args ...Arg) (ks []string, vs []any, err error) {
	for _, a := range args {
		if a.String {
			ks = append(ks, "$"+strings.TrimPrefix(a.Key, "$"))
			vs = append(vs, a.Val)
		} else {
			var j any
			if j, err = decode(a.Val); err != nil {
				return nil, nil, err
			}
			ks = append(ks, "$"+strings.TrimPrefix(a.Key, "$"))
			vs = append(vs, j)
		}
	}
	return
}

func jsonScalarToString(input interface{}) (string, bool) { //nolint:unused
	switch tt := input.(type) {
	case string:
		return `"` + tt + `"`, true
	case float64:
		if math.Trunc(tt) == tt {
			return strconv.FormatFloat(tt, 'f', 0, 64), true
		}
		return strconv.FormatFloat(tt, 'f', 2, 64), true
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

func rawFunc(a any, _ []any) any {
	if t := reflect.TypeOf(a); t != nil && t.Kind() == reflect.String {
		return strings.Trim(a.(string), `"`)
	} else if a == nil {
		return "null"
	}
	return a
}

func decode(s string) (v any, err error) {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	err = dec.Decode(&v)
	return
}
