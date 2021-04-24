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

// Package formatter provides various implementations to convert arbitrary values
// to different string representations such as JSON.
package formatter

import (
	"errors"
	"reflect"
)

// ErrUnsupported is the error resulting if a Formatter does not support the type.
var ErrUnsupported = errors.New("unsupported type")

// Formatter converts the given parameter to its string representation.
type Formatter interface {
	// Format returns a suitable string representation of i.
	//
	// Format must not modify the given parameter, even temporarily.
	Format(i interface{}) (string, error)
}

// Func is an adapter to allow the use of ordinary functions as Formatters.
// If f is a function with the appropriate signature,
// Func(f) is a Formatter that calls f.
type Func func(i interface{}) (string, error)

// Format returns a string by applying f to i.
func (f Func) Format(i interface{}) (string, error) {
	return f(i)
}

// Noop always returns an empty string and no error.
func Noop(_ interface{}) (s string, err error) {
	return
}

// NoopFormatter returns a simple Formatter that always returns an empty string and no error.
func NoopFormatter() Formatter { return Func(Noop) }

// CompFormatter combines multiple Formatters, each handling a different type.
type CompFormatter struct {
	byType map[string]Formatter
}

// NewComp creates and initializes a new CompFormatter.
func NewComp() *CompFormatter {
	return &CompFormatter{make(map[string]Formatter)}
}

// Format converts the given parameter to its string representation.
// If none of the registered Formatters can handle the given value, an error is returned.
func (cf CompFormatter) Format(i interface{}) (string, error) {
	if f, ok := cf.byType[typeName(reflect.TypeOf(i))]; ok {
		return f.Format(i)
	}
	return "", ErrUnsupported
}

// SetFormatter registers the Formatter for the given type.
// If a Formatter already exists for the type, it is replaced.
func (cf *CompFormatter) SetFormatter(n string, f Formatter) {
	cf.byType[n] = f
}

// SetFormatterFunc registers the Func for the given type.
// If a Formatter already exists for the type, it is replaced.
func (cf *CompFormatter) SetFormatterFunc(n string, f Func) {
	cf.byType[n] = f
}

// typeName returns the type's name.
// If the type cannot be determined e.g., []interface{}, a string representation
// is returned instead.
//
// Note that a string representation is not necessarily unique among types.
func typeName(typ reflect.Type) string {
	n := typ.Name()
	if n == "" {
		n = typ.String()
	}
	return n
}
