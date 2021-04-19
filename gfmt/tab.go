package gfmt

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"github.com/abc-inc/gutenfmt/renderer"
)

type Tab struct {
	w        *countingWriter
	Renderer *renderer.CompRenderer
}

func NewTab(w io.Writer) *Tab {
	return &Tab{wrapCountingWriter(w), renderer.NewComp()}
}

func (f Tab) Write(i interface{}) (int, error) {
	if i == nil {
		return 0, nil
	}

	if s, err := f.Renderer.Render(i); err == nil {
		return f.w.WriteString(s)
	} else if err != renderer.ErrUnsupported {
		return 0, err
	}

	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		return f.Write(reflect.Indirect(reflect.ValueOf(i)).Interface())
	} else if !isCompositeType(typ) {
		return fmt.Fprint(f.w, i)
	}

	f.w.cnt = 0
	tw := tabwriter.NewWriter(f.w, 1, 4, 3, ' ', 0)

	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		_, err := f.writeSlice(tw, reflect.ValueOf(i))
		return int(f.w.cnt), err
	case reflect.Map:
		err := f.writeMap(tw, reflect.ValueOf(i).MapRange())
		return int(f.w.cnt), err
	default:
		_, err := f.writeStructKeyVal(tw, i)
		return int(f.w.cnt), err
	}
}

func (f Tab) writeSlice(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	if v.Len() == 0 {
		return 0, nil
	}

	if v.Type().Elem().Kind() == reflect.Struct {
		n, err := f.writeStructTab(tw, v)
		return n, err
	} else if v.Type().Elem().Kind() == reflect.Map {
		n, err := f.writeMapTab(tw, v.Interface())
		return n, err
	}

	n, err := f.w.WriteString(renderer.StrVal(v.Index(0).Interface()))
	if err != nil {
		return n, err
	}

	cnt := n
	for idx := 1; idx < v.Len(); idx++ {
		n, err = f.w.WriteString("\n" + renderer.StrVal(v.Index(idx).Interface()))
		cnt += n
		if err != nil {
			return cnt, err
		}
	}
	return cnt, tw.Flush()
}

func (f Tab) writeMap(tw *tabwriter.Writer, iter *reflect.MapIter) error {
	for iter.Next() {
		k := renderer.StrVal(iter.Key().Interface())
		v := renderer.StrVal(iter.Value().Interface())
		if _, err := fmt.Fprintf(tw, "%v\t%v\t\n", k, v); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (f Tab) writeMapTab(tw *tabwriter.Writer, i interface{}) (int, error) {
	r := renderer.FromMapSlice(reflect.TypeOf(i))
	s, err := r.Render(i)
	if err != nil {
		return 0, err
	}

	n, err := tw.Write([]byte(s))
	tw.Flush()
	return n, err
}

func (f Tab) writeStructKeyVal(tw *tabwriter.Writer, i interface{}) (int, error) {
	r := renderer.FromStruct(reflect.TypeOf(i))
	s, err := r.Render(i)
	if err != nil {
		return 0, err
	}

	n, err := tw.Write([]byte(s))
	tw.Flush()
	return n, err
}

func (f Tab) writeStructTab(tw *tabwriter.Writer, v reflect.Value) (int, error) {
	r := renderer.FromStructSlice(v.Type(), "\t")
	s, err := r.Render(v.Interface())
	if err != nil {
		return 0, err
	}

	n, err := tw.Write([]byte(s))
	tw.Flush()
	return n, err
}
