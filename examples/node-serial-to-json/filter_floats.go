package main

import (
	"fmt"
	"math"
	"reflect"
)

// float values "NaN", "+Inf", "-Inf" cannot be converted into JSON numbers.
// convert them to string before performing the JSON encoding.
func filterFloats(obj any) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	convert := func(f float64) interface{} {
		switch {
		case math.IsNaN(f):
			return "NaN"
		case math.IsInf(f, 1):
			return "+Inf"
		case math.IsInf(f, -1):
			return "-Inf"
		}
		return f
	}

	switch v.Kind() {
	case reflect.Float32:
		return convert(float64(v.Interface().(float32)))

	case reflect.Float64:
		return convert(v.Interface().(float64))

	case reflect.Bool, reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64, reflect.String:
		return v.Interface()

	case reflect.Array:
		nl := v.Len()
		out := make([]interface{}, nl)
		for i := 0; i < nl; i++ {
			out[i] = filterFloats(v.Index(i).Interface())
		}
		return out

	case reflect.Struct:
		t := v.Type()
		nf := t.NumField()
		out := make(map[string]interface{}, nf)
		for i := 0; i < nf; i++ {
			out[t.Field(i).Name] = filterFloats(v.Field(i).Interface())
		}
		return out

	default:
		panic(fmt.Errorf("unsupported type: %s", v.Kind()))
	}
}
