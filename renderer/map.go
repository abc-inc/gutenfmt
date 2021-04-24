package renderer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/abc-inc/gutenfmt/internal/meta"
)

func FromMap(sep, delim string) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		return FromMapKeys(sep, delim, reflect.ValueOf(i).MapKeys()...).Render(i)
	})
}

func FromMapKeys(sep, delim string, ks ...reflect.Value) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		m := reflect.ValueOf(i)
		b := &strings.Builder{}
		for _, k := range ks {
			v := meta.StrVal(m.MapIndex(k).Interface())
			n := meta.StrVal(k.Interface())
			if _, err := fmt.Fprintf(b, "%s%s%s%s", n, sep, v, delim); err != nil {
				return "", err
			}
		}
		return b.String(), nil
	})
}

func FromMapSlice(sep, delim string, typ reflect.Type) Renderer {
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
		if v.Len() == 0 {
			return "", nil
		}

		b := &strings.Builder{}

		var ks []interface{}
		e := v.Index(0)
		for idx, k := range e.MapKeys() {
			if idx > 0 {
				b.WriteString(sep)
			}
			n := meta.StrVal(k.Interface())
			if !contains(ks, n) {
				ks = append(ks, n)
				b.WriteString(n)
			}
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			m := v.Index(i)
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := m.MapIndex(reflect.ValueOf(k.(string)))
				if v.IsValid() {
					b.WriteString(meta.StrVal(v))
				}
			}
		}

		return b.String(), nil
	})
}

func FromMapSliceKeys(sep, delim string, ks ...reflect.Value) Renderer {
	if len(ks) == 0 {
		return NoopRenderer()
	}

	return RendererFunc(func(i interface{}) (string, error) {
		v := reflect.ValueOf(i)
		b := &strings.Builder{}

		b.WriteString(meta.StrVal(ks[0].Interface()))
		for idx := 1; idx < len(ks); idx++ {
			b.WriteString(sep)
			b.WriteString(meta.StrVal(ks[idx].Interface()))
		}

		for i := 0; i < v.Len(); i++ {
			b.WriteString(delim)
			for idx, k := range ks {
				if idx > 0 {
					b.WriteString(sep)
				}
				v := v.Index(i).MapIndex(k)
				if v.IsValid() {
					b.WriteString(meta.StrVal(v))
				}
			}
		}

		return b.String(), nil
	})
}
