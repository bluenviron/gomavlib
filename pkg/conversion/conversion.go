// Package conversion contains functions to convert definitions from XML to Go.
package conversion

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

var (
	reMsgName     = regexp.MustCompile("^[A-Z0-9_]+$")
	reTypeIsArray = regexp.MustCompile(`^(.+?)\[([0-9]+)\]$`)
)

var tplDialect = template.Must(template.New("").Parse(
	`// Package {{ .PkgName }} contains the {{ .PkgName }} dialect.
//
//autogenerated:yes
package {{ .PkgName }}

import (
	"github.com/bluenviron/gomavlib/v3/pkg/message"
	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
)

// Dialect contains the dialect definition.
var Dialect = dial

// dial is not exposed directly in order not to display it in godoc.
var dial = &dialect.Dialect{
	Version: {{.Version}},
	Messages: []message.Message{
{{- range .Defs }}
	// {{ .Name }}
{{- range .Messages }}
		&Message{{ .Name }}{},
{{- end }}
{{- end }}
	},
}
`))

var tplEnum = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package {{ .PkgName }}

{{- if .Link }}

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/{{ .Enum.DefName }}"
)

{{- range .Enum.Description }}
// {{ . }}
{{- end }}
type {{ .Enum.Name }} = {{ .Enum.DefName }}.{{ .Enum.Name }}

const (
{{- $en := .Enum }}
{{- range .Enum.Values }}
{{- range .Description }}
		// {{ . }}
{{- end }}
		{{ .Name }} {{ $en.Name }} = {{ $en.DefName }}.{{ .Name }}
{{- end }}
)

{{- else }}

import (
	"strconv"
{{- if .Enum.Bitmask }}
	"strings"
{{- end }}
	"fmt"
)

{{- range .Enum.Description }}
// {{ . }}
{{- end }}
type {{ .Enum.Name }} uint64

const (
{{- $pn := .Enum.Name }}
{{- range .Enum.Values }}
{{- range .Description }}
	// {{ . }}
{{- end }}
	{{ .Name }} {{ $pn }} = {{ .Value }}
{{- end }}
)

{{- if .Enum.Bitmask }}
var values_{{ .Enum.Name }} = []{{ .Enum.Name }}{
{{- range .Enum.Values }}
	{{ .Name }},
{{- end }}
}
{{- end }}

var value_to_label_{{ .Enum.Name }} = map[{{ .Enum.Name }}]string{
{{- range .Enum.Values }}
	{{ .Name }}: "{{ .Name }}",
{{- end }}
}

var label_to_value_{{ .Enum.Name }} = map[string]{{ .Enum.Name }}{
{{- range .Enum.Values }}
	"{{ .Name }}": {{ .Name }},
{{- end }}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e {{ .Enum.Name }}) MarshalText() ([]byte, error) {
{{- if .Enum.Bitmask }}
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for _, val := range values_{{ .Enum.Name }} {
		if e&val == val {
			names = append(names, value_to_label_{{ .Enum.Name }}[val])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
{{- else }}
	if name, ok := value_to_label_{{ .Enum.Name }}[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
{{- end }}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *{{ .Enum.Name }}) UnmarshalText(text []byte) error {
{{- if .Enum.Bitmask }}
	labels := strings.Split(string(text), " | ")
	var mask {{ .Enum.Name }}
	for _, label := range labels {
		if value, ok := label_to_value_{{ .Enum.Name }}[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= {{ .Enum.Name }}(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	*e = mask
{{- else }}
	if value, ok := label_to_value_{{ .Enum.Name }}[string(text)]; ok {
	   *e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
	   *e = {{ .Enum.Name }}(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
{{- end }}
	return nil
}

// String implements the fmt.Stringer interface.
func (e {{ .Enum.Name }}) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
{{- end }}
`))

var tplMessage = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive,misspell,govet,lll
package {{ .PkgName }}

{{- if .Link }}

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/{{ .Msg.DefName }}"
)

{{- range .Msg.Description }}
// {{ . }}
{{- end }}
type Message{{ .Msg.Name }} = {{ .Msg.DefName }}.Message{{ .Msg.Name }}

{{- else }}

{{- range .Msg.Description }}
// {{ . }}
{{- end }}
type Message{{ .Msg.Name }} struct {
{{- range .Msg.Fields }}
{{- range .Description }}
	// {{ . }}
{{- end }}
	{{ .Line }}
{{- end }}
}

// GetID implements the message.Message interface.
func (*Message{{ .Msg.Name }}) GetID() uint32 {
	return {{ .Msg.ID }}
}

{{- end }}
`))

var dialectTypeToGo = map[string]string{
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
}

func defAddrToName(pa string) string {
	var b string
	u, err := url.ParseRequestURI(pa)
	if err == nil {
		b = path.Base(u.Path)
	} else {
		b = path.Base(pa)
	}

	b = strings.TrimSuffix(b, path.Ext(b))
	return strings.ToLower(strings.ReplaceAll(b, "_", ""))
}

func dialectNameGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

func dialectNameDefToGo(in string) string {
	re := regexp.MustCompile("_[a-z]")
	in = strings.ToLower(in)
	in = re.ReplaceAllStringFunc(in, func(match string) string {
		return strings.ToUpper(match[1:2])
	})
	return strings.ToUpper(in[:1]) + in[1:]
}

func parseDescription(in string) []string {
	var lines []string

	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

func uintPow(base, exp uint64) uint64 {
	result := uint64(1)
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

type outEnumValue struct {
	Value       uint64
	Name        string
	Description []string
}

type outEnum struct {
	DefName     string
	Name        string
	Description []string
	Values      []*outEnumValue
	Bitmask     bool
}

type outField struct {
	Description []string
	Line        string
}

type outMessage struct {
	DefName     string
	OrigName    string
	Name        string
	Description []string
	ID          int
	Fields      []*outField
}

type outDefinition struct {
	Name     string
	Enums    []*outEnum
	Messages []*outMessage
}

func processDefinition(
	version *string,
	processedDefs map[string]struct{},
	isRemote bool,
	defAddr string,
) ([]*outDefinition, error) {
	// skip already processed
	if _, ok := processedDefs[defAddr]; ok {
		return nil, nil
	}
	processedDefs[defAddr] = struct{}{}

	fmt.Fprintf(os.Stderr, "processing definition %s\n", defAddr)

	content, err := getDefinition(isRemote, defAddr)
	if err != nil {
		return nil, err
	}

	def, err := definitionDecode(content)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %w", err)
	}

	addrPath, _ := filepath.Split(defAddr)

	var outDefs []*outDefinition

	// includes
	for _, subDefAddr := range def.Includes {
		// prepend url to remote address
		if isRemote {
			subDefAddr = addrPath + subDefAddr
		}
		var subDefs []*outDefinition
		subDefs, err = processDefinition(version, processedDefs, isRemote, subDefAddr)
		if err != nil {
			return nil, err
		}
		outDefs = append(outDefs, subDefs...)
	}

	// version (process it after includes, in order to allow overriding it)
	if def.Version != "" {
		*version = def.Version
	}

	outDef := &outDefinition{
		Name: defAddrToName(defAddr),
	}

	// enums
	for _, enum := range def.Enums {
		oute := &outEnum{
			DefName:     outDef.Name,
			Name:        enum.Name,
			Description: parseDescription(enum.Description),
			Bitmask:     enum.Bitmask,
		}

		for _, entry := range enum.Entries {
			var v uint64

			switch {
			case strings.HasPrefix(entry.Value, "0b"):
				var tmp uint64
				tmp, err = strconv.ParseUint(entry.Value[2:], 2, 64)
				if err != nil {
					return nil, err
				}
				v = tmp

			case strings.HasPrefix(entry.Value, "0x"):
				var tmp uint64
				tmp, err = strconv.ParseUint(entry.Value[2:], 16, 64)
				if err != nil {
					return nil, err
				}
				v = tmp

			case strings.Contains(entry.Value, "**"):
				parts := strings.SplitN(entry.Value, "**", 2)

				var x uint64
				x, err = strconv.ParseUint(parts[0], 10, 64)
				if err != nil {
					return nil, err
				}

				var y uint64
				y, err = strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return nil, err
				}

				v = uintPow(x, y)

			default:
				var tmp uint64
				tmp, err = strconv.ParseUint(entry.Value, 10, 64)
				if err != nil {
					return nil, err
				}
				v = tmp
			}

			oute.Values = append(oute.Values, &outEnumValue{
				Value:       v,
				Name:        entry.Name,
				Description: parseDescription(entry.Description),
			})
		}

		outDef.Enums = append(outDef.Enums, oute)
	}

	// messages
	for _, msg := range def.Messages {
		var outMsg *outMessage
		outMsg, err = processMessage(outDef.Name, msg)
		if err != nil {
			return nil, err
		}
		outDef.Messages = append(outDef.Messages, outMsg)
	}

	outDefs = append(outDefs, outDef)
	return outDefs, nil
}

func getDefinition(isRemote bool, defAddr string) ([]byte, error) {
	if isRemote {
		byt, err := download(defAddr)
		if err != nil {
			return nil, fmt.Errorf("unable to download: %w", err)
		}
		return byt, nil
	}

	byt, err := os.ReadFile(defAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to open: %w", err)
	}
	return byt, nil
}

func download(addr string) ([]byte, error) {
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad return code: %v", res.StatusCode)
	}

	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func processMessage(defName string, msgDef *definitionMessage) (*outMessage, error) {
	if m := reMsgName.FindStringSubmatch(msgDef.Name); m == nil {
		return nil, fmt.Errorf("unsupported message name: %s", msgDef.Name)
	}

	outMsg := &outMessage{
		DefName:     defName,
		OrigName:    msgDef.Name,
		Name:        dialectNameDefToGo(msgDef.Name),
		Description: parseDescription(msgDef.Description),
		ID:          msgDef.ID,
	}

	for _, f := range msgDef.Fields {
		outField, err := processField(f)
		if err != nil {
			return nil, err
		}
		outMsg.Fields = append(outMsg.Fields, outField)
	}

	return outMsg, nil
}

func processField(fieldDef *dialectField) (*outField, error) {
	outF := &outField{
		Description: parseDescription(fieldDef.Description),
	}
	tags := make(map[string]string)

	newname := dialectNameDefToGo(fieldDef.Name)

	// name conversion is not univoque: add tag
	if dialectNameGoToDef(newname) != fieldDef.Name {
		tags["mavname"] = fieldDef.Name
	}

	outF.Line += newname

	typ := fieldDef.Type
	arrayLen := ""

	if typ == "uint8_t_mavlink_version" {
		typ = "uint8_t"
	}

	// string or array
	if matches := reTypeIsArray.FindStringSubmatch(typ); matches != nil {
		// string
		if matches[1] == "char" {
			tags["mavlen"] = matches[2]
			typ = "char"
			// array
		} else {
			arrayLen = matches[2]
			typ = matches[1]
		}
	}

	// extension
	if fieldDef.Extension {
		tags["mavext"] = "true"
	}

	typ = dialectTypeToGo[typ]
	if typ == "" {
		return nil, fmt.Errorf("unknown type: %s", typ)
	}

	outF.Line += " "
	if arrayLen != "" {
		outF.Line += "[" + arrayLen + "]"
	}
	if fieldDef.Enum != "" {
		outF.Line += fieldDef.Enum
		tags["mavenum"] = typ
	} else {
		outF.Line += typ
	}

	if len(tags) > 0 {
		var tmp []string
		for k, v := range tags {
			tmp = append(tmp, fmt.Sprintf("%s:\"%s\"", k, v))
		}
		sort.Strings(tmp)
		outF.Line += " `" + strings.Join(tmp, " ") + "`"
	}
	return outF, nil
}

func writeDialect(
	dir string,
	defName string,
	version string,
	outDefs []*outDefinition,
	enums map[string]*outEnum,
) error {
	var buf bytes.Buffer
	err := tplDialect.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Version": func() int {
			ret, _ := strconv.Atoi(version)
			return ret
		}(),
		"Defs":  outDefs,
		"Enums": enums,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "dialect.go"), buf.Bytes(), 0o644)
}

func writeEnum(
	dir string,
	defName string,
	enum *outEnum,
	link bool,
) error {
	var buf bytes.Buffer
	err := tplEnum.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Enum":    enum,
		"Link":    link && defName != enum.DefName,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "enum_"+strings.ToLower(enum.Name)+".go"), buf.Bytes(), 0o644)
}

func writeMessage(
	dir string,
	defName string,
	msg *outMessage,
	link bool,
) error {
	var buf bytes.Buffer
	err := tplMessage.Execute(&buf, map[string]interface{}{
		"PkgName": defName,
		"Msg":     msg,
		"Link":    link && defName != msg.DefName,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "message_"+strings.ToLower(msg.OrigName)+".go"), buf.Bytes(), 0o644)
}

// Convert converts a XML definition into a Golang definition.
func Convert(path string, link bool) error {
	version := ""
	processedDefs := make(map[string]struct{})
	_, err := url.ParseRequestURI(path)
	isRemote := (err == nil)
	defName := defAddrToName(path)

	_, err = os.Stat(defName)
	if !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", defName)
	}

	os.Mkdir(defName, 0o755)

	// parse all definitions recursively
	outDefs, err := processDefinition(&version, processedDefs, isRemote, path)
	if err != nil {
		return err
	}

	// merge enums together
	enums := make(map[string]*outEnum)
	for _, def := range outDefs {
		for _, defEnum := range def.Enums {
			if _, ok := enums[defEnum.Name]; !ok {
				enums[defEnum.Name] = defEnum
			} else {
				enums[defEnum.Name].DefName = defName
				enums[defEnum.Name].Values = append(enums[defEnum.Name].Values, defEnum.Values...)
			}
		}
	}

	err = writeDialect(defName, defName, version, outDefs, enums)
	if err != nil {
		return err
	}

	for _, enum := range enums {
		err = writeEnum(defName, defName, enum, link)
		if err != nil {
			return err
		}
	}

	for _, def := range outDefs {
		for _, msg := range def.Messages {
			err = writeMessage(defName, defName, msg, link)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
