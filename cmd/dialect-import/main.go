package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"
)

var reMsgName = regexp.MustCompile("^[A-Z0-9_]+$")
var reTypeIsArray = regexp.MustCompile(`^(.+?)\[([0-9]+)\]$`)

var tplDialect = template.Must(template.New("").Parse(
	`//nolint:golint,misspell,govet
{{- if .Comment -}}
// {{ .Comment }}
{{- end }}
package {{ .PkgName }}

import (
{{- if .Enums }}
	"errors"
	"strconv"
{{- end }}

	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/aler9/gomavlib/pkg/dialect"
)

// Dialect contains the dialect object that can be passed to the library.
var Dialect = dial

// dialect is not exposed directly such that it is not displayed in godoc.
var dial = &dialect.Dialect{ {{.Version}}, []msg.Message{
{{- range .Defs }}
    // {{ .Name }}
{{- range .Messages }}
    &Message{{ .Name }}{},
{{- end }}
{{- end }}
} }

{{ range .Enums }}
// {{ .Description }}
type {{ .Name }} int

const (
{{- $pn := .Name }}
{{- range .Values }}
	// {{ .Description }}
	{{ .Name }} {{ $pn }} = {{ .Value }}
{{- end }}
)

// MarshalText implements the encoding.TextMarshaler interface.
func (e {{ .Name }}) MarshalText() ([]byte, error) {
	switch e {
{{- range .Values }}
	case {{ .Name }}:
		return []byte("{{ .Name }}"), nil
{{- end }}
	}
	return nil, errors.New("invalid value")
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *{{ .Name }}) UnmarshalText(text []byte) error {
	switch string(text) {
{{- range .Values }}
	case "{{ .Name }}":
		*e = {{ .Name }}
		return nil
{{- end }}
	}
	return errors.New("invalid value")
}

// String implements the fmt.Stringer interface.
func (e {{ .Name }}) String() string {
	byts, err := e.MarshalText()
	if err == nil {
		return string(byts)
	}
	return strconv.FormatInt(int64(e), 10)
}

{{ end }}

{{ range .Defs }}
// {{ .Name }}

{{ range .Messages }}
// {{ .Description }}
type Message{{ .Name }} struct {
{{- range .Fields }}
	// {{ .Description }}
    {{ .Line }}
{{- end }}
}

// GetID implements the msg.Message interface.
func (*Message{{ .Name }}) GetID() uint32 {
    return {{ .ID }}
}
{{ end }}
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

func dialectFieldGoToDef(in string) string {
	re := regexp.MustCompile("([A-Z])")
	in = re.ReplaceAllString(in, "_${1}")
	return strings.ToLower(in[1:])
}

func dialectFieldDefToGo(in string) string {
	return dialectMsgDefToGo(in)
}

func dialectMsgDefToGo(in string) string {
	re := regexp.MustCompile("_[a-z]")
	in = strings.ToLower(in)
	in = re.ReplaceAllStringFunc(in, func(match string) string {
		return strings.ToUpper(match[1:2])
	})
	return strings.ToUpper(in[:1]) + in[1:]
}

func filterDesc(in string) string {
	return strings.Replace(in, "\n", "", -1)
}

type outEnumValue struct {
	Value       string
	Name        string
	Description string
}

type outEnum struct {
	Name        string
	Description string
	Values      []*outEnumValue
}

type outField struct {
	Description string
	Line        string
}

type outMessage struct {
	Name        string
	Description string
	ID          int
	Fields      []*outField
}

type outDefinition struct {
	Name     string
	Enums    []*outEnum
	Messages []*outMessage
}

func definitionProcess(version *string, defsProcessed map[string]struct{}, isRemote bool, defAddr string) ([]*outDefinition, error) {
	// skip already processed
	if _, ok := defsProcessed[defAddr]; ok {
		return nil, nil
	}
	defsProcessed[defAddr] = struct{}{}

	fmt.Fprintf(os.Stderr, "definition %s\n", defAddr)

	content, err := definitionGet(isRemote, defAddr)
	if err != nil {
		return nil, err
	}

	def, err := definitionDecode(content)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %s", err)
	}

	addrPath, addrName := filepath.Split(defAddr)

	var outDefs []*outDefinition

	// version
	if def.Version != "" {
		if *version != "" && *version != def.Version {
			return nil, fmt.Errorf("version defined twice (%s and %s)", def.Version, *version)
		}
		*version = def.Version
	}

	// includes
	for _, inc := range def.Includes {
		// prepend url to remote address
		if isRemote {
			inc = addrPath + inc
		}
		subDefs, err := definitionProcess(version, defsProcessed, isRemote, inc)
		if err != nil {
			return nil, err
		}
		outDefs = append(outDefs, subDefs...)
	}

	outDef := &outDefinition{
		Name: addrName,
	}

	// enums
	for _, enum := range def.Enums {
		oute := &outEnum{
			Name:        enum.Name,
			Description: filterDesc(enum.Description),
		}
		for _, val := range enum.Values {
			oute.Values = append(oute.Values, &outEnumValue{
				Value:       val.Value,
				Name:        val.Name,
				Description: filterDesc(val.Description),
			})
		}
		outDef.Enums = append(outDef.Enums, oute)
	}

	// messages
	for _, msg := range def.Messages {
		outMsg, err := messageProcess(msg)
		if err != nil {
			return nil, err
		}
		outDef.Messages = append(outDef.Messages, outMsg)
	}

	outDefs = append(outDefs, outDef)
	return outDefs, nil
}

func definitionGet(isRemote bool, defAddr string) ([]byte, error) {
	if isRemote {
		byt, err := download(defAddr)
		if err != nil {
			return nil, fmt.Errorf("unable to download: %s", err)
		}
		return byt, nil
	}

	byt, err := ioutil.ReadFile(defAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to open: %s", err)
	}
	return byt, nil
}

func download(desturl string) ([]byte, error) {
	res, err := http.Get(desturl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad return code: %v", res.StatusCode)
	}

	byt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func messageProcess(msg *definitionMessage) (*outMessage, error) {
	if m := reMsgName.FindStringSubmatch(msg.Name); m == nil {
		return nil, fmt.Errorf("unsupported message name: %s", msg.Name)
	}

	outMsg := &outMessage{
		Name:        dialectMsgDefToGo(msg.Name),
		Description: filterDesc(msg.Description),
		ID:          msg.ID,
	}

	for _, f := range msg.Fields {
		outField, err := fieldProcess(f)
		if err != nil {
			return nil, err
		}
		outMsg.Fields = append(outMsg.Fields, outField)
	}

	return outMsg, nil
}

func fieldProcess(field *dialectField) (*outField, error) {
	outF := &outField{
		Description: filterDesc(field.Description),
	}
	tags := make(map[string]string)

	newname := dialectFieldDefToGo(field.Name)

	// name conversion is not univoque: add tag
	if dialectFieldGoToDef(newname) != field.Name {
		tags["mavname"] = field.Name
	}

	outF.Line += newname

	typ := field.Type
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
	if field.Extension {
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
	if field.Enum != "" {
		outF.Line += field.Enum
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

func run() error {
	kingpin.CommandLine.Help = "Convert Mavlink dialects from XML format into Go format."

	argPkgName := kingpin.Flag("package", "Package name").Default("main").String()
	argComment := kingpin.Flag("comment", "comment to add before the package name").Default("").String()
	argMainDef := kingpin.Arg("xml", "Path or url pointing to a XML Mavlink dialect").Required().String()

	kingpin.Parse()

	mainDef := *argMainDef
	comment := *argComment
	pkgName := *argPkgName

	version := ""
	defsProcessed := make(map[string]struct{})
	isRemote := func() bool {
		_, err := url.ParseRequestURI(mainDef)
		return err == nil
	}()

	// parse all definitions recursively
	outDefs, err := definitionProcess(&version, defsProcessed, isRemote, mainDef)
	if err != nil {
		return err
	}

	// merge enums together
	enums := make(map[string]*outEnum)
	for _, def := range outDefs {
		for _, defEnum := range def.Enums {
			if _, ok := enums[defEnum.Name]; !ok {
				enums[defEnum.Name] = &outEnum{
					Name:        defEnum.Name,
					Description: defEnum.Description,
				}
			}
			enum := enums[defEnum.Name]

			enum.Values = append(enum.Values, defEnum.Values...)
		}
	}

	// fill enum missing values
	for _, enum := range enums {
		nextVal := 0
		for _, v := range enum.Values {
			if v.Value != "" {
				nextVal, _ = strconv.Atoi(v.Value)
				nextVal++
			} else {
				v.Value = strconv.Itoa(nextVal)
				nextVal++
			}
		}
	}

	// dump
	return tplDialect.Execute(os.Stdout, map[string]interface{}{
		"PkgName": pkgName,
		"Comment": comment,
		"Version": func() int {
			ret, _ := strconv.Atoi(version)
			return ret
		}(),
		"Defs":  outDefs,
		"Enums": enums,
	})
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
}
