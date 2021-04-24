package renderer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

// FromStruct creates a new Renderer for a struct type.
//
// It looks for the "json" key in the fields' tag and uses the name, if defined.
// As a special case, if the field tag is "-", the field is omitted.
// Options like "omitempty" are ignored.
func FromStruct(sep, delim string, typ reflect.Type) Renderer {
	fns, pns := jsonMetadata(typ)
	if len(pns) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}
		for idx, fn := range fns {
			s := fmt.Sprintf("%s%s%s%s", pns[idx], sep, meta.StrVal(v.FieldByName(fn).Interface()), delim)
			if _, err := b.WriteString(s); err != nil {
				return "", err
			}
		}
		return b.String()[:b.Len()-len(delim)], nil
	})
}

// FromStructSlice creates a new Renderer for a struct slice.
//
// The fields are determined by the slice's element type.
// Like FromStruct, it looks for the "json" key in the struct field's tag.
func FromStructSlice(sep, delim string, typ reflect.Type) Renderer {
	_, pns := jsonMetadata(typ.Elem())
	if len(pns) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		b.WriteString(strings.Join(pns, sep))
		b.WriteString(delim)

		v := reflect.ValueOf(i)
		for idx := 0; idx < v.Len(); idx++ {
			e := v.Index(idx)
			for p := range pns {
				b.WriteString(meta.StrVal(e.Field(p).Interface()))
				b.WriteString(sep)
			}
			b.WriteString(delim)
		}
		return b.String(), nil
	})
}
