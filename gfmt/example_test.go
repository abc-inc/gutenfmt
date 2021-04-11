package gfmt_test

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	. "github.com/abc-inc/gutenfmt/gfmt"
	"github.com/abc-inc/gutenfmt/renderer"
)

func ExampleJSON_Write_struct() {
	o := NewJSON(os.Stdout)
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	// The custom renderer demonstrates how to combine literal strings and JSON on a struct that is not intended for serialization.
	o.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		t := i.(Team)
		return `{"team":"` + t.Name() + `","members":` + strconv.Itoa(len(t.Members())) + `}`, nil
	})

	_, _ = o.Write(u)
	_, _ = o.Write("\n\n")
	_, _ = o.Write(t)
	// Output:
	// {"username":"John Doe","email":"john.doe@local"}
	//
	// {"team":"SUPPORT","members":2}
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
	t := *NewTeam("Support", u, u)

	o := NewTab(b)
	typ := renderer.TypeName(reflect.TypeOf([]Team{}))
	o.Renderer.SetRendererFunc(typ, func(i interface{}) (string, error) {
		b := strings.Builder{}
		b.WriteString("name\tmembers\t\n")

		ts := i.([]Team)
		for _, t := range ts {
			b.WriteString(fmt.Sprintf("%s\t%d\t\n", t.Name(), len(t.Members())))
		}
		return b.String(), nil
	})

	_, _ = o.Write([]User{u, u})
	_, _ = o.Write("\n")
	_, _ = o.Write([]Team{t, t})

	// Since the Output cannot contain trailing spaces, it gets stripped from the table in this example.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username   email
	// John Doe   john.doe@local
	// John Doe   john.doe@local
	// name      members
	// SUPPORT   2
	// SUPPORT   2
}

func ExampleText_Write_struct() {
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	o := NewText(os.Stdout)
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
	u := *NewUser("John", "Doe")
	t := *NewTeam("Support", u, u)

	o := NewText(os.Stdout)
	o.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = o.Write([]User{u, u})
	_, _ = o.Write("\n\n")
	_, _ = o.Write([]Team{t, t})
	// Output:
	// John Doe <john.doe@local>
	// John Doe <john.doe@local>
	//
	// SUPPORT
	// SUPPORT
}
