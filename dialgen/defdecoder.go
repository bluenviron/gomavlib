package main

import (
	"encoding/xml"
	"strconv"
)

type definitionEnumValue struct {
	Value       string `xml:"value,attr"`
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`
}

type definitionEnum struct {
	Name        string                 `xml:"name,attr"`
	Description string                 `xml:"description"`
	Values      []*definitionEnumValue `xml:"entry"`
}

type dialectField struct {
	Extension   bool   `xml:"-"`
	Type        string `xml:"type,attr"`
	Name        string `xml:"name,attr"`
	Enum        string `xml:"enum,attr"`
	Description string `xml:",innerxml"`
}

type definitionMessage struct {
	Id          int
	Name        string
	Description string
	Fields      []*dialectField
}

// UnmarshalXML implements xml.Unmarshaler
// we must unmarshal manually due to extension fields
func (m *definitionMessage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				field := &dialectField{Extension: inExtensions}
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

type definition struct {
	Version  string               `xml:"version"`
	Dialect  int                  `xml:"dialect"`
	Includes []string             `xml:"include"`
	Enums    []*definitionEnum    `xml:"enums>enum"`
	Messages []*definitionMessage `xml:"messages>message"`
}

func definitionDecode(content []byte) (*definition, error) {
	def := &definition{}
	err := xml.Unmarshal(content, def)
	return def, err
}
