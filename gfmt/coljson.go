package gfmt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/abc-inc/gutenfmt/renderer"
	"github.com/alecthomas/chroma/formatters"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers/j"
	"github.com/alecthomas/chroma/styles"
)

type ColJSON struct {
	w        *countingWriter
	Renderer *renderer.CompRenderer
	Prefix   string
	Indent   string
}

func NewColJSON(w io.Writer) *ColJSON {
	return &ColJSON{wrapCountingWriter(w), renderer.NewComp(), "", "  "}
}

func (f ColJSON) Write(i interface{}) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := f.Renderer.Render(i); err == nil {
		return io.WriteString(f.w, s)
	} else if err != renderer.ErrUnsupported {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return f.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !isCompositeType(typ) {
		return fmt.Fprint(f.w, renderer.StrVal(i))
	}

	j, err := json.MarshalIndent(i, f.Prefix, f.Indent)
	if err != nil {
		return 0, err
	}

	f.w.cnt = 0
	if err = f.highlight(f.w, string(j), "native"); err != nil {
		return 0, err
	}
	return int(f.w.cnt), nil
}

func (f ColJSON) highlight(w io.Writer, source, style string) error {
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	l := chroma.Coalesce(j.JSON)
	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return formatters.TTY8.Format(w, s, it)
}
