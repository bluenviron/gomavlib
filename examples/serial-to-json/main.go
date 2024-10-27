package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint
// 2) encode incoming messages into JSON
// 3) print messages in the console

// float values "NaN", "+Inf", "-Inf" cannot be converted into JSON numbers.
// convert them to string before performing the JSON encoding.
func convertIncompatibleFloats(obj any) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	convert := func(f float64) interface{} {
		if math.IsNaN(f) {
			return "NaN"
		} else if math.IsInf(f, 1) {
			return "+Inf"
		} else if math.IsInf(f, -1) {
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
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return v.Interface()

	case reflect.Array:
		nl := v.Len()
		out := make([]interface{}, nl)
		for i := 0; i < nl; i++ {
			out[i] = convertIncompatibleFloats(v.Index(i).Interface())
		}
		return out

	case reflect.Struct:
		t := v.Type()
		nf := t.NumField()
		out := make(map[string]interface{}, nf)
		for i := 0; i < nf; i++ {
			out[t.Field(i).Name] = convertIncompatibleFloats(v.Field(i).Interface())
		}
		return out

	default:
		panic(fmt.Errorf("unsupported type: %s", v.Kind()))
	}
}

func main() {
	// create a node which communicates with a serial endpoint
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			// encode incoming messages
			enc, err := json.Marshal(struct {
				Type    string
				Content interface{}
			}{
				Type:    fmt.Sprintf("%T", frm.Message()),
				Content: convertIncompatibleFloats(frm.Message()),
			})
			if err != nil {
				panic(err)
			}

			// print messages in the console
			log.Printf("%s\n", enc)
		}
	}
}
