package renderer

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func strFormat(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Uintptr:
		return fmt.Sprint(reflect.Indirect(v).Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprint(v.Complex())
	case reflect.Array:
		return strings.Trim(fmt.Sprint(v.Interface()), "[]")
	case reflect.Chan:
		return ""
	case reflect.Func:
		return funcName(v)
	case reflect.Interface:
		return fmt.Sprint(v.Interface())
	case reflect.Map:
		return fmt.Sprint(v.Interface())
	case reflect.Ptr:
		return strFormat(reflect.Indirect(v))
	case reflect.Slice:
		return strings.Trim(fmt.Sprint(v.Interface()), "[]")
	case reflect.String:
		return v.String()
	case reflect.Struct:
		return fmt.Sprint(v.Interface())
	// case reflect.UnsafePointer:
	default:
		panic("Cannot convert " + reflect.TypeOf(v).Name() + " (" + fmt.Sprintf("%v", v.Type()) + ")")
	}
}

func TypeName(typ reflect.Type) string {
	n := typ.Name()
	if n == "" {
		n = typ.String()
	}
	return n
}

func funcName(f reflect.Value) string {
	return runtime.FuncForPC(f.Pointer()).Name()
}

func jsonMetadata(typ reflect.Type) ([]string, []string) {
	var fns []string
	var pns []string
	for idx := 0; idx < typ.NumField(); idx++ {
		f := typ.Field(idx)
		if n := jsonPropName(f); n != "" {
			fns = append(fns, f.Name)
			pns = append(pns, n)
		}
	}
	return fns, pns
}

func jsonPropName(f reflect.StructField) string {
	if n, ok := f.Tag.Lookup("json"); strings.HasPrefix(n, "-") || f.PkgPath != "" {
		return ""
	} else if n = strings.SplitN(n, ",", 2)[0]; ok && n != "" {
		return n
	}
	return f.Name
}
