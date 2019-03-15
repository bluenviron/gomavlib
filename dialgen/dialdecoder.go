package main

import (
	"encoding/xml"
	"strconv"
)

// DialectEnumValue is the part of a dialect definition that contains an enum value.
type DialectEnumValue struct {
	Value       string `xml:"value,attr"`
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
}

// DialectEnum is the part of a dialect definition that contains an enum.
type DialectEnum struct {
	Name        string              `xml:"name,attr"`
	Description string              `xml:"description"`
	Values      []*DialectEnumValue `xml:"entry"`
}

// DialectField is the part of a dialect definition that contains a field.
type DialectField struct {
	Extension   bool   `xml:"-"`
	Type        string `xml:"type,attr"`
	Name        string `xml:"name,attr"`
	Enum        string `xml:"enum,attr"`
	Description string `xml:",innerxml"`
}

// DialectMessage is the part of a dialect definition that contains a message.
type DialectMessage struct {
	Id          int
	Name        string
	Description string
	Fields      []*DialectField
}

// UnmarshalXML implements xml.Unmarshaler
// we must unmarshal manually due to extension fields
func (m *DialectMessage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// unmarshal attributes
	for _, a := range start.Attr {
		switch a.Name.Local {
		case "id":
			v, _ := strconv.Atoi(a.Value)
			m.Id = v
		case "name":
			m.Name = a.Value
		}
	}

	inExtensions := false
	for {
		t, _ := d.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "description":
				err := d.DecodeElement(&m.Description, &se)
				if err != nil {
					return err
				}

			case "extensions":
				inExtensions = true

			case "field":
				field := &DialectField{Extension: inExtensions}
				err := d.DecodeElement(&field, &se)
				if err != nil {
					return err
				}
				m.Fields = append(m.Fields, field)
			}
		}
	}
	return nil
}

// DialectDefinition is a Mavlink dialect definition.
type DialectDefinition struct {
	Version  int               `xml:"version"`
	Dialect  int               `xml:"dialect"`
	Includes []string          `xml:"include"`
	Enums    []*DialectEnum    `xml:"enums>enum"`
	Messages []*DialectMessage `xml:"messages>message"`
}

func dialectDecode(content []byte) (*DialectDefinition, error) {
	def := &DialectDefinition{}
	err := xml.Unmarshal(content, def)
	return def, err
}
