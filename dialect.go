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

	"github.com/aler9/gomavlib/x25"
)

type dialectFieldType int

const (
	typeDouble dialectFieldType = iota + 1
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

var dialectFieldTypeFromGo = map[string]dialectFieldType{
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

var dialectFieldTypeString = map[dialectFieldType]string{
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

var dialectFieldTypeSizes = map[dialectFieldType]byte{
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

func dialectFieldGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

func dialectMsgGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToUpper(in[1:])
}

// Dialect contains available messages and the configuration needed to encode and
// decode them.
type Dialect struct {
	version  uint
	messages map[uint32]*dialectMessage
}

// NewDialect allocates a Dialect.
func NewDialect(version uint, messages []Message) (*Dialect, error) {
	d := &Dialect{
		version:  version,
		messages: make(map[uint32]*dialectMessage),
	}

	for _, msg := range messages {
		mp, err := newDialectMessage(msg)
		if err != nil {
			return nil, fmt.Errorf("message %T: %s", msg, err)
		}
		d.messages[msg.GetId()] = mp
	}

	return d, nil
}

// MustDialect is like NewDialect but panics in case of error.
func MustDialect(version uint, messages []Message) *Dialect {
	d, err := NewDialect(version, messages)
	if err != nil {
		panic(err)
	}
	return d
}

type dialectMessageField struct {
	isEnum      bool
	ftype       dialectFieldType
	name        string
	arrayLength byte
	index       int
	isExtension bool
}

type dialectMessage struct {
	elemType     reflect.Type
	fields       []*dialectMessageField
	sizeNormal   byte
	sizeExtended byte
	crcExtra     byte
}

func newDialectMessage(msg Message) (*dialectMessage, error) {
	mp := &dialectMessage{}
	mp.elemType = reflect.TypeOf(msg).Elem()

	mp.fields = make([]*dialectMessageField, mp.elemType.NumField())

	// get name
	if strings.HasPrefix(mp.elemType.Name(), "Message") == false {
		return nil, fmt.Errorf("message struct name must begin with 'Message'")
	}
	msgName := dialectMsgGoToDef(mp.elemType.Name()[len("Message"):])

	// collect message fields
	for i := 0; i < mp.elemType.NumField(); i++ {
		field := mp.elemType.Field(i)
		arrayLength := byte(0)
		goType := field.Type

		// array
		if goType.Kind() == reflect.Array {
			arrayLength = byte(goType.Len())
			goType = goType.Elem()
		}

		isEnum := false
		var dialectType dialectFieldType

		// enum
		if field.Tag.Get("mavenum") != "" {
			isEnum = true

			if goType.Kind() != reflect.Int {
				return nil, fmt.Errorf("an enum must be an int")
			}

			tagEnum := field.Tag.Get("mavenum")

			if len(tagEnum) == 0 {
				return nil, fmt.Errorf("enum but tag not specified")
			}

			dialectType = dialectFieldTypeFromGo[tagEnum]
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
			dialectType = dialectFieldTypeFromGo[goType.Name()]
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
			size = dialectFieldTypeSizes[dialectType] * arrayLength
		} else {
			size = dialectFieldTypeSizes[dialectType]
		}

		mp.fields[i] = &dialectMessageField{
			isEnum: isEnum,
			ftype:  dialectType,
			name: func() string {
				if mavname := field.Tag.Get("mavname"); mavname != "" {
					return mavname
				}
				return dialectFieldGoToDef(field.Name)
			}(),
			arrayLength: arrayLength,
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
			if w1, w2 := dialectFieldTypeSizes[mp.fields[i].ftype], dialectFieldTypeSizes[mp.fields[j].ftype]; w1 != w2 {
				return w1 > w2
			}
		}
		// sort by original index
		return mp.fields[i].index < mp.fields[j].index
	})

	// generate CRC extra
	// https://mavlink.io/en/guide/serialization.html#crc_extra
	mp.crcExtra = func() byte {
		h := x25.New()
		h.Write([]byte(msgName + " "))

		for _, f := range mp.fields {
			// skip extensions
			if f.isExtension == true {
				continue
			}

			h.Write([]byte(dialectFieldTypeString[f.ftype] + " "))
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

func (mp *dialectMessage) decode(buf []byte, isFrameV2 bool) (Message, error) {
	msg := reflect.New(mp.elemType)

	if isFrameV2 == true {
		// in V2 buffer length can be > message or < message
		// in this latter case it must be filled with zeros to support empty-byte de-truncation
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
				n := decodeValue(target.Index(i), buf, f)
				buf = buf[n:]
			}

		default:
			n := decodeValue(target, buf, f)
			buf = buf[n:]
		}
	}

	return msg.Interface().(Message), nil
}

func (mp *dialectMessage) encode(msg Message, isFrameV2 bool) ([]byte, error) {
	var buf []byte

	if isFrameV2 == true {
		buf = make([]byte, mp.sizeExtended)
	} else {
		buf = make([]byte, mp.sizeNormal)
	}

	start := buf

	// encode field by field
	for _, f := range mp.fields {
		// skip extensions in V1 frames
		if isFrameV2 == false && f.isExtension == true {
			continue
		}

		target := reflect.ValueOf(msg).Elem().Field(f.index)

		switch target.Kind() {
		case reflect.Array:
			length := target.Len()
			for i := 0; i < length; i++ {
				n := encodeValue(buf, target.Index(i), f)
				buf = buf[n:]
			}

		default:
			n := encodeValue(buf, target, f)
			buf = buf[n:]
		}
	}

	buf = start

	// empty-byte truncation
	// even with truncation, message length must be at least 1 byte
	// https://github.com/mavlink/c_library_v2/blob/master/mavlink_helpers.h#L103
	if isFrameV2 == true {
		end := len(buf)
		for end > 1 && buf[end-1] == 0x00 {
			end--
		}
		buf = buf[:end]
	}

	return buf, nil
}

func decodeValue(target reflect.Value, buf []byte, f *dialectMessageField) int {
	if f.isEnum == true {
		switch f.ftype {
		case typeUint8:
			target.SetInt(int64(buf[0]))
			return 1

		case typeInt8:
			target.SetInt(int64(buf[0]))
			return 1

		case typeUint16:
			target.SetInt(int64(binary.LittleEndian.Uint16(buf)))
			return 2

		case typeUint32:
			target.SetInt(int64(binary.LittleEndian.Uint32(buf)))
			return 4

		case typeInt32:
			target.SetInt(int64(binary.LittleEndian.Uint32(buf)))
			return 4

		case typeUint64:
			target.SetInt(int64(binary.LittleEndian.Uint64(buf)))
			return 8

		default:
			panic("unexpected type")
		}
	}

	switch tt := target.Addr().Interface().(type) {
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
}

func encodeValue(buf []byte, target reflect.Value, f *dialectMessageField) int {
	if f.isEnum == true {
		switch f.ftype {
		case typeUint8:
			buf[0] = byte(target.Int())
			return 1

		case typeInt8:
			buf[0] = byte(target.Int())
			return 1

		case typeUint16:
			binary.LittleEndian.PutUint16(buf, uint16(target.Int()))
			return 2

		case typeUint32:
			binary.LittleEndian.PutUint32(buf, uint32(target.Int()))
			return 4

		case typeInt32:
			binary.LittleEndian.PutUint32(buf, uint32(target.Int()))
			return 4

		case typeUint64:
			binary.LittleEndian.PutUint64(buf, uint64(target.Int()))
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
