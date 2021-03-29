package parser

import (
	"bytes"
	"encoding/xml"

	"github.com/semaphoreci/test-results/pkg/logger"
)

// XMLElement ...
type XMLElement struct {
	XMLName    xml.Name
	Attributes map[string]string `xml:"-"`
	Children   []XMLElement      `xml:",any"`
	Contents   []byte            `xml:",chardata"`
}

// NewXMLElement ...
func NewXMLElement() XMLElement {
	return XMLElement{}
}

// Attr ...
func (me *XMLElement) Attr(attr string) string {
	return me.Attributes[attr]
}

// Tag ...
func (me *XMLElement) Tag() string {
	return me.XMLName.Local
}

// UnmarshalXML ...
func (me *XMLElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	logger.Trace("Decoding element: %s", start.Name.Local)
	type alias XMLElement
	if err := d.DecodeElement((*alias)(me), &start); err != nil {
		logger.Error("Decoding element failed: %v", err)
		return err
	}

	me.Attributes = parseAttributes(start.Attr)
	return nil
}

// Parse ...
func (me *XMLElement) Parse(reader *bytes.Reader) error {
	logger.Debug("Parsing element started")
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&me); err != nil {
		logger.Error("Parsing element failed")
		return err
	}
	return nil
}

func parseAttributes(attrs []xml.Attr) map[string]string {
	attributes := make(map[string]string)

	for _, attr := range attrs {
		switch attr.Name {
		default:
			attributes[attr.Name.Local] = attr.Value
		}
	}

	return attributes
}
