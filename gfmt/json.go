package gfmt

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/abc-inc/gutenfmt/renderer"
)

type JSON struct {
	w        io.Writer
	Renderer *renderer.CompRenderer
}

func NewJSON(w io.Writer) *JSON {
	return &JSON{w, renderer.NewComp()}
}

func (f JSON) Write(i interface{}) (int, error) {
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

	j, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}
	return f.w.Write(j)
}
