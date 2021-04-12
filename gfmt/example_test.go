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

var u = *NewUser("John", "Doe")
var t = *NewTeam("Support", u, u)

func ExampleJSON_Write_struct() {
	f := NewJSON(os.Stdout)
	f.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		t := i.(Team)
		return `{"team":"` + t.Name() + `","members":` + strconv.Itoa(len(t.Members())) + `}`, nil
	})

	_, _ = f.Write(u)
	_, _ = f.Write("\n\n")
	_, _ = f.Write(t)
	// Output:
	// {"username":"John Doe","email":"john.doe@local"}
	//
	// {"team":"SUPPORT","members":2}
}

func ExampleJSON_Write_structSlice() {
	f := NewJSON(os.Stdout)

	_, _ = f.Write([]User{u, u})
	// Output:
	// [{"username":"John Doe","email":"john.doe@local"},{"username":"John Doe","email":"john.doe@local"}]
}

func ExampleTab_Write_struct() {
	b := &strings.Builder{}
	f := NewTab(b)
	f.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return fmt.Sprintf("name\t%s\t", i.(Team).Name()), nil
	})

	_, _ = f.Write(u)
	_, _ = f.Write("\n")
	_, _ = f.Write(t)

	// Since the Output cannot contain trailing spaces, it gets stripped from the table in this Example.
	s := regexp.MustCompile(`\s+\n`).ReplaceAllString(b.String(), "\n")
	fmt.Println(s)
	// Output:
	// username   John Doe
	// email      john.doe@local
	// name   SUPPORT
}

func ExampleTab_Write_structSlice() {
	b := &strings.Builder{}
	f := NewTab(b)
	typ := renderer.TypeName(reflect.TypeOf([]Team{}))
	f.Renderer.SetRendererFunc(typ, func(i interface{}) (string, error) {
		b := strings.Builder{}
		b.WriteString("name\tmembers\t\n")

		ts := i.([]Team)
		for _, t := range ts {
			b.WriteString(fmt.Sprintf("%s\t%d\t\n", t.Name(), len(t.Members())))
		}
		return b.String(), nil
	})

	_, _ = f.Write([]User{u, u})
	_, _ = f.Write("\n")
	_, _ = f.Write([]Team{t, t})

	// Since the Example Output cannot contain trailing spaces, they get stripped.
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
	f := NewText(os.Stdout)
	f.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = f.Write(t)
	_, _ = f.Write("\n\n")
	_, _ = f.Write(u)
	// Output:
	// SUPPORT
	//
	// John Doe <john.doe@local>
}

func ExampleText_Write_structSlice() {
	f := NewText(os.Stdout)
	f.Renderer.SetRendererFunc(reflect.TypeOf(t).Name(), func(i interface{}) (string, error) {
		return i.(Team).Name(), nil
	})

	_, _ = f.Write([]User{u, u})
	_, _ = f.Write("\n\n")
	_, _ = f.Write([]Team{t, t})
	// Output:
	// John Doe <john.doe@local>
	// John Doe <john.doe@local>
	//
	// SUPPORT
	// SUPPORT
}
