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
	"github.com/abc-inc/gutenfmt/gfmt"
)

var user = *NewUser("John", "Doe")
var team = *NewTeam("Support", user, user)

func ExampleJSON_Write_struct() {
	w := gfmt.NewJSON(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(team).Name(), func(i any) (string, error) {
		tm := i.(Team)
		return `{"team":"` + tm.Name() + `","members":` + strconv.Itoa(len(tm.Members())) + `}`, nil
	})

	_, _ = w.Write(user)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(team)
	// Output:
	// {"username":"John Doe","email":"john.doe@local"}
	//
	// {"team":"SUPPORT","members":2}
}

func ExampleJSON_Write_structSlice() {
	w := gfmt.NewJSON(os.Stdout)

	_, _ = w.Write([]User{user, user})
	// Output:
	// [{"username":"John Doe","email":"john.doe@local"},{"username":"John Doe","email":"john.doe@local"}]
}

func ExampleTab_Write_struct() {
	b := &strings.Builder{}
	w := gfmt.NewTab(b)
	f := formatter.AsTab(formatter.Func(func(i any) (string, error) {
		return fmt.Sprintf("name\t%s\t", i.(Team).Name()), nil
	}))
	w.Formatter.SetFormatter(reflect.TypeOf(team).Name(), f)

	_, _ = w.Write(user)
	_, _ = w.Write("\n")
	_, _ = w.Write(team)

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
	w := gfmt.NewTab(b)
	typ := reflect.TypeOf([]Team{}).String()
	w.Formatter.SetFormatter(typ, formatter.AsTab(formatter.Func(func(i any) (string, error) {
		buf := strings.Builder{}
		buf.WriteString("name\tmembers\t\n")

		ts := i.([]Team)
		for _, t := range ts {
			buf.WriteString(fmt.Sprintf("%s\t%d\t\n", t.Name(), len(t.Members())))
		}
		return buf.String(), nil
	})))

	_, _ = w.Write([]User{user, *NewUser("Rudolf", "Lingens")})
	_, _ = w.Write("\n")
	_, _ = w.Write([]Team{team, team})

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
	w := gfmt.NewText(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(team).Name(), func(i any) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write(team)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(user)
	// Output:
	// SUPPORT
	//
	// John Doe <john.doe@local>
}

func ExampleText_Write_structSlice() {
	w := gfmt.NewText(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf([]Team{}).String(), func(i any) (string, error) {
		b := strings.Builder{}
		for _, t := range i.([]Team) {
			b.WriteString(t.Name() + "\n")
		}
		return b.String(), nil
	})

	_, _ = w.Write([]User{user, user})
	_, _ = w.Write("\n\n")
	_, _ = w.Write([]Team{team, team})
	// Output:
	// username:email
	// John Doe:john.doe@local
	// John Doe:john.doe@local
	//
	// SUPPORT
	// SUPPORT
}

func ExampleYAML_Write_struct() {
	w := gfmt.NewYAML(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(team).Name(), func(i any) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write(team)
	_, _ = w.Write("\n\n")
	_, _ = w.Write(user)
	// Output:
	// SUPPORT
	//
	// Username: John Doe
	// E-Mail: john.doe@local
}

func ExampleYAML_Write_structSlice() {
	w := gfmt.NewYAML(os.Stdout)
	w.Formatter.SetFormatterFunc(reflect.TypeOf(team).Name(), func(i any) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = w.Write([]User{user, user})
	// Output:
	// - Username: John Doe
	//   E-Mail: john.doe@local
	// - Username: John Doe
	//   E-Mail: john.doe@local
}
