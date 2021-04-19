package renderer

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"
)

var ErrUnsupported = errors.New("unsupported type")

type Renderer interface {
	Render(i interface{}) (string, error)
}

type RendererFunc func(i interface{}) (string, error)

func (r RendererFunc) Render(i interface{}) (string, error) {
	return r(i)
}

var NoopRenderer = RendererFunc(func(i interface{}) (s string, err error) {
	return
})

type CompRenderer struct {
	byType map[string]Renderer
}

func NewComp() *CompRenderer {
	return &CompRenderer{make(map[string]Renderer)}
}

func (cr CompRenderer) Render(i interface{}) (string, error) {
	r, ok := cr.byType[TypeName(reflect.TypeOf(i))]
	if !ok {
		return "", ErrUnsupported
	}
	return r.Render(i)
}

func (cr *CompRenderer) SetRenderer(n string, r Renderer) {
	cr.byType[n] = r
}

func (cr *CompRenderer) SetRendererFunc(n string, r RendererFunc) {
	cr.byType[n] = r
}

func FromStruct(typ reflect.Type) Renderer {
	fns, pns := jsonMetadata(typ)
	if len(pns) == 0 {
		return NoopRenderer
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}
		for i, fn := range fns {
			if _, err := b.WriteString(fmt.Sprintf("%s\t%s\t\n", pns[i], StrVal(v.FieldByName(fn).Interface()))); err != nil {
				return "", err
			}
		}
		return b.String()[:b.Len()-1], nil
	})
}

func FromStructSlice(typ reflect.Type, delim string) Renderer {
	_, pns := jsonMetadata(typ.Elem())
	if len(pns) == 0 {
		return NoopRenderer
	}

	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		b.WriteString(strings.Join(pns, delim))
		b.WriteString(delim + "\n")

		v := reflect.ValueOf(i)
		for idx := 0; idx < v.Len(); idx++ {
			e := v.Index(idx)
			for p := range pns {
				b.WriteString(StrVal(e.Field(p).Interface()))
				b.WriteString(delim)
			}
			b.WriteByte('\n')
		}
		return b.String(), nil
	})
}

func FromTemplate(tmpl template.Template) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		if err := tmpl.Execute(b, i); err != nil {
			return "", err
		}
		return b.String(), nil
	})
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func FromMap() Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		iter := v.MapRange()

		b := &strings.Builder{}
		for iter.Next() {
			k := iter.Key().Interface()
			v := iter.Value().Interface()
			if _, err := fmt.Fprintf(b, "%v\t%v\t\n", StrVal(k), StrVal(v)); err != nil {
				return "", err
			}
		}
		return b.String(), nil
	})
}

func FromMapSlice(typ reflect.Type) Renderer {
	contains := func(es []interface{}, s interface{}) bool {
		for _, e := range es {
			if e == s {
				return true
			}
		}
		return false
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		var ks []interface{}
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			for _, k := range e.MapKeys() {
				n := StrVal(k.Interface())
				if !contains(ks, n) {
					ks = append(ks, n)
					b.WriteString(n)
					b.WriteByte('\t')
				}
			}
		}
		b.WriteByte('\n')

		for i := 0; i < v.Len(); i++ {
			m := v.Index(i)
			for _, n := range ks {
				v := m.MapIndex(reflect.ValueOf(n.(string)))
				if v.IsValid() {
					b.WriteString(StrVal(v))
				}
				b.WriteString("\t")
			}
			b.WriteString("\n")
		}

		return b.String()[:b.Len()-1], nil
	})
}

func FromMapSliceKeys(sep string, ks ...reflect.Value) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		for _, k := range ks {
			b.WriteString(StrVal(k.Interface()))
			b.WriteString(sep)
		}
		b.WriteByte('\n')

		for i := 0; i < v.Len(); i++ {
			for _, k := range ks {
				v := v.Index(i).MapIndex(k)
				if v.IsValid() {
					b.WriteString(StrVal(v))
				}
				b.WriteString(sep)
			}
			b.WriteByte('\n')
		}

		return b.String()[:b.Len()-1], nil
	})
}

func StrVal(i interface{}) string {
	return strFormat(reflect.ValueOf(i))
}

func AsTab(r Renderer) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		s, err := r.Render(i)
		if err != nil {
			return "", err
		}

		b := &strings.Builder{}
		tw := tabwriter.NewWriter(b, 1, 4, 3, ' ', 0)
		tw.Write([]byte(s))
		if err = tw.Flush(); err != nil {
			return "", err
		}
		return b.String(), nil
	})
}
