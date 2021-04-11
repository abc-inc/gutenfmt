package gfmt_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	. "github.com/abc-inc/gutenfmt/gfmt"
)

func ExampleJSON_Write_struct() {
	o := NewJSON(os.Stdout)
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	// The custom renderer demonstrates how to combine literal strings and JSON on a struct that was not intended for serialization.
	o.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		ms, _ := json.Marshal(i.(Team).Members())
		return `{"team":"` + i.(Team).Name() + `","members":` + string(ms) + `}`, nil
	})

	_, _ = o.Write(t)
	_, _ = o.Write("\n\n")
	_, _ = o.Write(u)
	// Output:
	// {"team":"SUPPORT","members":[{"username":"John Doe","email":"john.doe@local"},{"username":"John Doe","email":"john.doe@local"}]}
	//
	// {"username":"John Doe","email":"john.doe@local"}
}

func ExampleJSON_Write_structSlice() {
	o := NewJSON(os.Stdout)
	u := *NewUser("John", "Doe")

	_, _ = o.Write([]User{u, u})
	// Output:
	// [{"username":"John Doe","email":"john.doe@local"},{"username":"John Doe","email":"john.doe@local"}]
}

func ExampleTab_Write_struct() {
	b := &strings.Builder{}
	o := NewTab(b)
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	o.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return fmt.Sprintf("name\t%s\t", i.(Team).Name()), nil
	})

	_, _ = o.Write(u)
	_, _ = o.Write("\n")
	_, _ = o.Write(t)

	// Since the Output cannot contain trailing spaces, it gets stripped from the table in this example.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username   John Doe
	// email      john.doe@local
	// name   SUPPORT
}

func ExampleTab_Write_structSlice() {
	b := &strings.Builder{}
	u := *NewUser("John", "Doe")
	_, _ = NewTab(b).Write([]User{u, u})

	// Since the Output cannot contain trailing spaces, it gets stripped from the table in this example.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username   email
	// John Doe   john.doe@local
	// John Doe   john.doe@local
}

func ExampleText_Write_struct() {
	o := NewText(os.Stdout)
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	o.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = o.Write(t)
	_, _ = o.Write("\n\n")
	_, _ = o.Write(u)
	// Output:
	// SUPPORT
	//
	// John Doe <john.doe@local>
}

func ExampleText_Write_structSlice() {
	o := NewText(os.Stdout)
	u := *NewUser("John", "Doe")

	_, _ = o.Write([]User{u, u})
	// Output:
	// John Doe <john.doe@local>
	// John Doe <john.doe@local>
}
