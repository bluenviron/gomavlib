package gomavlib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var dialectTypeSizes = map[string]byte{
	"double":   8,
	"uint64_t": 8,
	"int64_t":  8,
	"float":    4,
	"uint32_t": 4,
	"int32_t":  4,
	"uint16_t": 2,
	"int16_t":  2,
	"uint8_t":  1,
	"int8_t":   1,
	"char":     1,
}

// DialectTypeDefToGo converts a xml type into its go equivalent.
func DialectTypeDefToGo(typ string) string {
	return map[string]string{
		"double":   "float64",
		"uint64_t": "uint64",
		"int64_t":  "int64",
		"float":    "float32",
		"uint32_t": "uint32",
		"int32_t":  "int32",
		"uint16_t": "uint16",
		"int16_t":  "int16",
		"uint8_t":  "uint8",
		"int8_t":   "int8",
		"char":     "string",
	}[typ]
}

// DialectTypeGoToDef converts a go type into its xml equivalent.
func DialectTypeGoToDef(typ string) string {
	return map[string]string{
		"float64": "double",
		"uint64":  "uint64_t",
		"int64":   "int64_t",
		"float32": "float",
		"uint32":  "uint32_t",
		"int32":   "int32_t",
		"uint16":  "uint16_t",
		"int16":   "int16_t",
		"uint8":   "uint8_t",
		"int8":    "int8_t",
		"string":  "char",
	}[typ]
}

// DialectFieldGoToDef converts a go field into its xml equivalent.
func DialectFieldGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

// DialectFieldDefToGo converts a xml field into its go equivalent.
func DialectFieldDefToGo(in string) string {
	return DialectMsgDefToGo(in)
}

// DialectMsgGoToDef converts a go field into its xml equivalent.
func DialectMsgGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToUpper(in[1:])
}

// DialectMsgDefToGo converts a xml field into its go equivalent.
func DialectMsgDefToGo(in string) string {
	re := regexp.MustCompile("_[a-z]")
	in = strings.ToLower(in)
	in = re.ReplaceAllStringFunc(in, func(match string) string {
		return strings.ToUpper(match[1:2])
	})
	return strings.ToUpper(in[:1]) + in[1:]
}

type parserMessageField struct {
	ftype       string
	name        string
	arrayLength byte
	index       int
	isExtension bool
}

type parserMessage struct {
	elemType     reflect.Type
	fields       []parserMessageField
	sizeNormal   byte
	sizeExtended byte
	crcExtra     byte
}

func newParserMessage(msg Message) (*parserMessage, error) {
	mp := &parserMessage{}

	mp.elemType = reflect.TypeOf(msg).Elem()
	mp.fields = make([]parserMessageField, mp.elemType.NumField())

	// get name
	if strings.HasPrefix(mp.elemType.Name(), "Message") == false {
		return nil, fmt.Errorf("message struct name must begin with 'Message'")
	}
	msgName := DialectMsgGoToDef(mp.elemType.Name()[len("Message"):])

	// collect message fields
	for i := 0; i < mp.elemType.NumField(); i++ {
		field := mp.elemType.Field(i)
		fieldType := field.Type
		fieldArrayLength := byte(0)

		// array
		if fieldType.Kind() == reflect.Array {
			fieldArrayLength = byte(fieldType.Len())
			fieldType = fieldType.Elem()
		}

		// validate type
		defType := DialectTypeGoToDef(fieldType.Name())
		if defType == "" {
			fmt.Printf("invalid field type: %v\n", fieldType)
			return nil, fmt.Errorf("invalid field type: %v", fieldType)
		}

		// string
		if fieldType.Kind() == reflect.String {
			slen, err := strconv.Atoi(field.Tag.Get("mavlen"))
			if err != nil {
				return nil, err
			}
			fieldArrayLength = byte(slen)
		}

		// extension
		isExtension := (field.Tag.Get("mavext") == "true")

		// size
		var size byte
		if fieldArrayLength > 0 {
			size = dialectTypeSizes[defType] * fieldArrayLength
		} else {
			size = dialectTypeSizes[defType]
		}

		mp.fields[i] = parserMessageField{
			ftype: defType,
			name: func() string {
				if mavname := field.Tag.Get("mavname"); mavname != "" {
					return mavname
				}
				return DialectFieldGoToDef(field.Name)
			}(),
			arrayLength: fieldArrayLength,
			index:       i,
			isExtension: isExtension,
		}

		mp.sizeExtended += size
		if isExtension == false {
			mp.sizeNormal += size
		}
	}

	// reorder fields as described in
	// https://mavlink.io/en/guide/serialization.html#field_reordering
	sort.Slice(mp.fields, func(i, j int) bool {
		// sort by weight if not extension
		if mp.fields[i].isExtension == false && mp.fields[j].isExtension == false {
			if w1, w2 := dialectTypeSizes[mp.fields[i].ftype], dialectTypeSizes[mp.fields[j].ftype]; w1 != w2 {
				return w1 > w2
			}
		}
		// sort by original index
		return mp.fields[i].index < mp.fields[j].index
	})

	// generate CRC extra
	// https://mavlink.io/en/guide/serialization.html#crc_extra
	mp.crcExtra = func() byte {
		h := NewX25()
		h.Write([]byte(msgName + " "))

		for _, f := range mp.fields {
			// skip extensions
			if f.isExtension == true {
				continue
			}

			h.Write([]byte(f.ftype + " "))
			h.Write([]byte(f.name + " "))

			if f.arrayLength > 0 {
				h.Write([]byte{f.arrayLength})
			}
		}
		sum := h.Sum16()
		return byte((sum & 0xFF) ^ (sum >> 8))
	}()

	return mp, nil
}

func (mp *parserMessage) decode(buf []byte, isFrameV2 bool) (Message, error) {
	msg := reflect.New(mp.elemType)

	if isFrameV2 == true {
		// in V2 buffer can be > message or < message
		// in this case it must be filled with zeros to support empty-byte de-truncation
		// and extension fields
		if len(buf) < int(mp.sizeExtended) {
			buf = append(buf, bytes.Repeat([]byte{0x00}, int(mp.sizeExtended)-len(buf))...)
		}

	} else {
		// in V1 buffer must fit message perfectly
		if len(buf) != int(mp.sizeNormal) {
			return nil, fmt.Errorf("unexpected size (%d vs %d)", len(buf), mp.sizeNormal)
		}
	}

	// decode field by field
	for _, f := range mp.fields {
		// skip extensions in V1 frames
		if isFrameV2 == false && f.isExtension == true {
			continue
		}

		target := msg.Elem().Field(f.index)

		switch target.Kind() {
		case reflect.Array:
			length := target.Len()
			for i := 0; i < length; i++ {
				n := decodeValue(target.Index(i).Addr().Interface(), buf, nil)
				buf = buf[n:]
			}

		default:
			n := decodeValue(target.Addr().Interface(), buf, &f)
			buf = buf[n:]
		}
	}

	return msg.Interface().(Message), nil
}

func decodeValue(target interface{}, buf []byte, f *parserMessageField) int {
	switch tt := target.(type) {
	case *string:
		// find nil character or string end
		end := 0
		for buf[end] != 0 && end < int(f.arrayLength) {
			end++
		}
		*tt = string(buf[:end])
		return int(f.arrayLength) // return length including zeros

	case *int8:
		*tt = int8(buf[0])
		return 1

	case *uint8:
		*tt = buf[0]
		return 1

	case *int16:
		*tt = int16(binary.LittleEndian.Uint16(buf))
		return 2

	case *uint16:
		*tt = binary.LittleEndian.Uint16(buf)
		return 2

	case *int32:
		*tt = int32(binary.LittleEndian.Uint32(buf))
		return 4

	case *uint32:
		*tt = binary.LittleEndian.Uint32(buf)
		return 4

	case *int64:
		*tt = int64(binary.LittleEndian.Uint64(buf))
		return 8

	case *uint64:
		*tt = binary.LittleEndian.Uint64(buf)
		return 8

	case *float32:
		*tt = math.Float32frombits(binary.LittleEndian.Uint32(buf))
		return 4

	case *float64:
		*tt = math.Float64frombits(binary.LittleEndian.Uint64(buf))
		return 8

	default:
		panic("unexpected type")
	}
	return 0
}

func (mp *parserMessage) encode(msg Message, isFrameV2 bool) ([]byte, error) {
	bbuf := bytes.NewBuffer(nil)

	for _, f := range mp.fields {
		// skip extensions for v1 frames
		if isFrameV2 == false && f.isExtension == true {
			continue
		}

		target := reflect.ValueOf(msg).Elem().Field(f.index).Addr().Interface()

		switch tt := target.(type) {
		case *string:
			sbuf := make([]byte, f.arrayLength)
			copy(sbuf, *tt)
			_, err := bbuf.Write(sbuf)
			if err != nil {
				return nil, err
			}

		default:
			err := binary.Write(bbuf, binary.LittleEndian, target)
			if err != nil {
				return nil, err
			}
		}
	}

	buf := bbuf.Bytes()

	// empty-byte truncation
	if isFrameV2 == true {
		end := len(buf)
		for end > 0 && buf[end-1] == 0x00 {
			end--
		}
		buf = buf[:end]
	}

	return buf, nil
}
