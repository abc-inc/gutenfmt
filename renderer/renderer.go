package renderer

import (
	"errors"
	"reflect"
	"strings"
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

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}
		for i, fn := range fns {
			b.WriteString(pns[i])
			b.WriteByte('\t')
			b.WriteString(StrVal(v.FieldByName(fn).Interface()))
			b.WriteString("\t\n")
		}
		return b.String()[:b.Len()-1], nil
	})
}

func FromStructSlice(typ reflect.Type) Renderer {
	_, pns:=jsonMetadata(typ.Elem())
	if len(pns)== 0 {
		return NoopRenderer
	}

	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		for _, pn := range pns {
			b.WriteString(pn)
			b.WriteByte('\t')
		}
		b.WriteByte('\n')

		v := reflect.ValueOf(i)
		for idx := 0; idx < v.Len(); idx++ {
			e := v.Index(idx)
			for p := range pns {
				b.WriteString(StrVal(e.Field(p).Interface()))
				b.WriteByte('\t')
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

func FromMapSliceKeys(ks ...reflect.Value) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		for _, k := range ks {
			b.WriteString(StrVal(k.Interface()))
			b.WriteByte('\t')
		}
		b.WriteByte('\n')

		for i := 0; i < v.Len(); i++ {
			for _, k := range ks {
				v := v.Index(i).MapIndex(k)
				if v.IsValid() {
					b.WriteString(StrVal(v))
				}
				b.WriteByte('\t')
			}
			b.WriteByte('\n')
		}

		return b.String()[:b.Len()-1], nil
	})
}

func StrVal(i interface{}) string {
	return strFormat(reflect.ValueOf(i))
}
