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

package gfmt_test

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/abc-inc/gutenfmt/formatter"
	. "github.com/abc-inc/gutenfmt/gfmt"
)

var u = *NewUser("John", "Doe")
var t = *NewTeam("Support", u, u)

func ExampleJSON_Write_struct() {
	w := NewJSON(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		tm := i.(Team)
		return `{"team":"` + tm.Name() + `","members":` + strconv.Itoa(len(tm.Members())) + `}`, nil
	})

	_, _ = w.Write(u)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(t)
	// Output:
	// {"username":"John Doe","email":"john.doe@local"}
	//
	// {"team":"SUPPORT","members":2}
}

func ExampleJSON_Write_structSlice() {
	w := NewJSON(os.Stdout)

	_, _ = w.Write([]User{u, u})
	// Output:
	// [{"username":"John Doe","email":"john.doe@local"},{"username":"John Doe","email":"john.doe@local"}]
}

func ExampleTab_Write_struct() {
	b := &strings.Builder{}
	w := NewTab(b)
	f := formatter.AsTab(formatter.Func(func(i interface{}) (string, error) {
		return fmt.Sprintf("name\t%s\t", i.(Team).Name()), nil
	}))
	w.Formatter.SetFormatter(reflect.TypeOf(t).Name(), f)

	_, _ = w.Write(u)
	_, _ = w.Write("\n")
	_, _ = w.Write(t)

	// Since the Output cannot contain trailing spaces, it gets stripped from the table in this Example.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username John Doe
	// email    john.doe@local
	// name SUPPORT
}

func ExampleTab_Write_structSlice() {
	b := &strings.Builder{}
	w := NewTab(b)
	typ := reflect.TypeOf([]Team{}).String()
	w.Formatter.SetFormatter(typ, formatter.AsTab(formatter.Func(func(i interface{}) (string, error) {
		b := strings.Builder{}
		b.WriteString("name\tmembers\t\n")

		ts := i.([]Team)
		for _, t := range ts {
			b.WriteString(fmt.Sprintf("%s\t%d\t\n", t.Name(), len(t.Members())))
		}
		return b.String(), nil
	})))

	_, _ = w.Write([]User{u, *NewUser("Rudolf", "Lingens")})
	_, _ = w.Write("\n")
	_, _ = w.Write([]Team{t, t})

	// Since the Example Output cannot contain trailing spaces, they get stripped.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username       email
	// John Doe       john.doe@local
	// Rudolf Lingens rudolf.lingens@local
	// name    members
	// SUPPORT 2
	// SUPPORT 2
}

func ExampleText_Write_struct() {
	w := NewText(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write(t)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(u)
	// Output:
	// SUPPORT
	//
	// John Doe <john.doe@local>
}

func ExampleText_Write_structSlice() {
	w := NewText(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf([]Team{}).String(), func(i interface{}) (string, error) {
		b := strings.Builder{}
		for _, t := range i.([]Team) {
			b.WriteString(t.Name() + "\n")
		}
		return b.String(), nil
	})

	_, _ = w.Write([]User{u, u})
	_, _ = w.Write("\n\n")
	_, _ = w.Write([]Team{t, t})
	// Output:
	// username:email
	// John Doe:john.doe@local
	// John Doe:john.doe@local
	//
	// SUPPORT
	// SUPPORT
}

func ExampleYAML_Write_struct() {
	w := NewYAML(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write(t)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(u)
	// Output:
	// SUPPORT
	//
	// Username: John Doe
	// E-Mail: john.doe@local
}

func ExampleYAML_Write_structSlice() {
	w := NewYAML(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write([]User{u, u})
	// Output:
	// - Username: John Doe
	//   E-Mail: john.doe@local
	// - Username: John Doe
	//   E-Mail: john.doe@local
}
