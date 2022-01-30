package msg

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

	"github.com/aler9/gomavlib/pkg/x25"
)

type fieldType int

const (
	typeDouble fieldType = iota + 1
	typeUint64
	typeInt64
	typeFloat
	typeUint32
	typeInt32
	typeUint16
	typeInt16
	typeUint8
	typeInt8
	typeChar
)

var fieldTypeFromGo = map[string]fieldType{
	"float64": typeDouble,
	"uint64":  typeUint64,
	"int64":   typeInt64,
	"float32": typeFloat,
	"uint32":  typeUint32,
	"int32":   typeInt32,
	"uint16":  typeUint16,
	"int16":   typeInt16,
	"uint8":   typeUint8,
	"int8":    typeInt8,
	"string":  typeChar,
}

var fieldTypeString = map[fieldType]string{
	typeDouble: "double",
	typeUint64: "uint64_t",
	typeInt64:  "int64_t",
	typeFloat:  "float",
	typeUint32: "uint32_t",
	typeInt32:  "int32_t",
	typeUint16: "uint16_t",
	typeInt16:  "int16_t",
	typeUint8:  "uint8_t",
	typeInt8:   "int8_t",
	typeChar:   "char",
}

var fieldTypeSizes = map[fieldType]byte{
	typeDouble: 8,
	typeUint64: 8,
	typeInt64:  8,
	typeFloat:  4,
	typeUint32: 4,
	typeInt32:  4,
	typeUint16: 2,
	typeInt16:  2,
	typeUint8:  1,
	typeInt8:   1,
	typeChar:   1,
}

func fieldGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

func msgGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToUpper(in[1:])
}

type decEncoderField struct {
	isEnum      bool
	ftype       fieldType
	name        string
	arrayLength byte
	index       int
	isExtension bool
}

// DecEncoder is an object that allows to decode and encode a Message.
type DecEncoder struct {
	fields       []*decEncoderField
	sizeNormal   byte
	sizeExtended byte
	elemType     reflect.Type
	crcExtra     byte
}

// NewDecEncoder allocates a DecEncoder.
func NewDecEncoder(msg Message) (*DecEncoder, error) {
	mde := &DecEncoder{}
	mde.elemType = reflect.TypeOf(msg).Elem()

	mde.fields = make([]*decEncoderField, mde.elemType.NumField())

	// get name
	if !strings.HasPrefix(mde.elemType.Name(), "Message") {
		return nil, fmt.Errorf("message struct name must begin with 'Message'")
	}
	msgName := msgGoToDef(mde.elemType.Name()[len("Message"):])

	// collect message fields
	for i := 0; i < mde.elemType.NumField(); i++ {
		field := mde.elemType.Field(i)
		arrayLength := byte(0)
		goType := field.Type

		// array
		if goType.Kind() == reflect.Array {
			arrayLength = byte(goType.Len())
			goType = goType.Elem()
		}

		isEnum := false
		var dialectType fieldType

		// enum
		if field.Tag.Get("mavenum") != "" {
			isEnum = true

			if goType.Kind() != reflect.Uint32 {
				return nil, fmt.Errorf("an enum must be an uint32")
			}

			tagEnum := field.Tag.Get("mavenum")

			if len(tagEnum) == 0 {
				return nil, fmt.Errorf("enum but tag not specified")
			}

			dialectType = fieldTypeFromGo[tagEnum]
			if dialectType == 0 {
				return nil, fmt.Errorf("invalid go type: %v", tagEnum)
			}

			switch dialectType {
			case typeUint8:
			case typeInt8:
			case typeUint16:
			case typeUint32:
			case typeInt32:
			case typeUint64:
				break

			default:
				return nil, fmt.Errorf("type '%v' cannot be used as enum", tagEnum)
			}
		} else {
			dialectType = fieldTypeFromGo[goType.Name()]
			if dialectType == 0 {
				return nil, fmt.Errorf("invalid go type: %v", goType.Name())
			}

			// string or char
			if goType.Kind() == reflect.String {
				tagLen := field.Tag.Get("mavlen")

				// char
				if len(tagLen) == 0 {
					arrayLength = 1

					// string
				} else {
					slen, err := strconv.Atoi(tagLen)
					if err != nil {
						return nil, fmt.Errorf("string has invalid length: %v", tagLen)
					}
					arrayLength = byte(slen)
				}
			}
		}

		// extension
		isExtension := (field.Tag.Get("mavext") == "true")

		// size
		var size byte
		if arrayLength > 0 {
			size = fieldTypeSizes[dialectType] * arrayLength
		} else {
			size = fieldTypeSizes[dialectType]
		}

		mde.fields[i] = &decEncoderField{
			isEnum: isEnum,
			ftype:  dialectType,
			name: func() string {
				if mavname := field.Tag.Get("mavname"); mavname != "" {
					return mavname
				}
				return fieldGoToDef(field.Name)
			}(),
			arrayLength: arrayLength,
			index:       i,
			isExtension: isExtension,
		}

		mde.sizeExtended += size
		if !isExtension {
			mde.sizeNormal += size
		}
	}

	// reorder fields as described in
	// https://mavlink.io/en/guide/serialization.html#field_reordering
	sort.Slice(mde.fields, func(i, j int) bool {
		// sort by weight if not extension
		if !mde.fields[i].isExtension && !mde.fields[j].isExtension {
			if w1, w2 := fieldTypeSizes[mde.fields[i].ftype], fieldTypeSizes[mde.fields[j].ftype]; w1 != w2 {
				return w1 > w2
			}
		}
		// sort by original index
		return mde.fields[i].index < mde.fields[j].index
	})

	// generate CRC extra
	// https://mavlink.io/en/guide/serialization.html#crc_extra
	mde.crcExtra = func() byte {
		h := x25.New()
		h.Write([]byte(msgName + " "))

		for _, f := range mde.fields {
			// skip extensions
			if f.isExtension {
				continue
			}

			h.Write([]byte(fieldTypeString[f.ftype] + " "))
			h.Write([]byte(f.name + " "))

			if f.arrayLength > 0 {
				h.Write([]byte{f.arrayLength})
			}
		}
		sum := h.Sum16()
		return byte((sum & 0xFF) ^ (sum >> 8))
	}()

	return mde, nil
}

// CRCExtra returns the message CRC extra.
func (mde *DecEncoder) CRCExtra() byte {
	return mde.crcExtra
}

// Decode decodes a Message.
func (mde *DecEncoder) Decode(buf []byte, isV2 bool) (Message, error) {
	msg := reflect.New(mde.elemType)

	if isV2 {
		// in V2 buffer length can be > message or < message
		// in this latter case it must be filled with zeros to support empty-byte de-truncation
		// and extension fields
		if len(buf) < int(mde.sizeExtended) {
			buf = append(buf, bytes.Repeat([]byte{0x00}, int(mde.sizeExtended)-len(buf))...)
		}
	} else {
		// in V1 buffer must fit message perfectly
		if len(buf) != int(mde.sizeNormal) {
			return nil, fmt.Errorf("wrong size: expected %d, got %d", mde.sizeNormal, len(buf))
		}
	}

	// decode field by field
	for _, f := range mde.fields {
		// skip extensions in V1 frames
		if !isV2 && f.isExtension {
			continue
		}

		target := msg.Elem().Field(f.index)

		switch target.Kind() {
		case reflect.Array:
			length := target.Len()
			for i := 0; i < length; i++ {
				n := valueDecode(target.Index(i), buf, f)
				buf = buf[n:]
			}

		default:
			n := valueDecode(target, buf, f)
			buf = buf[n:]
		}
	}

	return msg.Interface().(Message), nil
}

// Encode encodes a message.
func (mde *DecEncoder) Encode(msg Message, isV2 bool) ([]byte, error) {
	var buf []byte

	if isV2 {
		buf = make([]byte, mde.sizeExtended)
	} else {
		buf = make([]byte, mde.sizeNormal)
	}

	start := buf

	// encode field by field
	for _, f := range mde.fields {
		// skip extensions in V1 frames
		if !isV2 && f.isExtension {
			continue
		}

		target := reflect.ValueOf(msg).Elem().Field(f.index)

		switch target.Kind() {
		case reflect.Array:
			length := target.Len()
			for i := 0; i < length; i++ {
				n := valueEncode(buf, target.Index(i), f)
				buf = buf[n:]
			}

		default:
			n := valueEncode(buf, target, f)
			buf = buf[n:]
		}
	}

	buf = start

	// empty-byte truncation
	// even with truncation, message length must be at least 1 byte
	// https://github.com/mavlink/c_library_v2/blob/master/mavlink_helpers.h#L103
	if isV2 {
		end := len(buf)
		for end > 1 && buf[end-1] == 0x00 {
			end--
		}
		buf = buf[:end]
	}

	return buf, nil
}

func valueDecode(target reflect.Value, buf []byte, f *decEncoderField) int {
	if f.isEnum {
		switch f.ftype {
		case typeUint8:
			target.SetUint(uint64(buf[0]))
			return 1

		case typeInt8:
			target.SetUint(uint64(buf[0]))
			return 1

		case typeUint16:
			target.SetUint(uint64(binary.LittleEndian.Uint16(buf)))
			return 2

		case typeUint32:
			target.SetUint(uint64(binary.LittleEndian.Uint32(buf)))
			return 4

		case typeInt32:
			target.SetUint(uint64(binary.LittleEndian.Uint32(buf)))
			return 4

		case typeUint64:
			target.SetUint(binary.LittleEndian.Uint64(buf))
			return 8

		default:
			panic("unexpected type")
		}
	}

	switch tt := target.Addr().Interface().(type) {
	case *string:
		// find string end or NULL character
		end := 0
		for end < int(f.arrayLength) && buf[end] != 0 {
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
}

func valueEncode(buf []byte, target reflect.Value, f *decEncoderField) int {
	if f.isEnum {
		switch f.ftype {
		case typeUint8:
			buf[0] = byte(target.Uint())
			return 1

		case typeInt8:
			buf[0] = byte(target.Uint())
			return 1

		case typeUint16:
			binary.LittleEndian.PutUint16(buf, uint16(target.Uint()))
			return 2

		case typeUint32:
			binary.LittleEndian.PutUint32(buf, uint32(target.Uint()))
			return 4

		case typeInt32:
			binary.LittleEndian.PutUint32(buf, uint32(target.Uint()))
			return 4

		case typeUint64:
			binary.LittleEndian.PutUint64(buf, target.Uint())
			return 8

		default:
			panic("unexpected type")
		}
	}

	switch tt := target.Addr().Interface().(type) {
	case *string:
		copy(buf[:f.arrayLength], *tt)
		return int(f.arrayLength) // return length including zeros

	case *int8:
		buf[0] = uint8(*tt)
		return 1

	case *uint8:
		buf[0] = *tt
		return 1

	case *int16:
		binary.LittleEndian.PutUint16(buf, uint16(*tt))
		return 2

	case *uint16:
		binary.LittleEndian.PutUint16(buf, *tt)
		return 2

	case *int32:
		binary.LittleEndian.PutUint32(buf, uint32(*tt))
		return 4

	case *uint32:
		binary.LittleEndian.PutUint32(buf, *tt)
		return 4

	case *int64:
		binary.LittleEndian.PutUint64(buf, uint64(*tt))
		return 8

	case *uint64:
		binary.LittleEndian.PutUint64(buf, *tt)
		return 8

	case *float32:
		binary.LittleEndian.PutUint32(buf, math.Float32bits(*tt))
		return 4

	case *float64:
		binary.LittleEndian.PutUint64(buf, math.Float64bits(*tt))
		return 8

	default:
		panic("unexpected type")
	}
}
