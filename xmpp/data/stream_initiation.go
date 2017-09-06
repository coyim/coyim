package data

import "encoding/xml"

// SI is a data element from http://xmpp.org/extensions/xep-0095.html.
type SI struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/si si"`
	ID       string   `xml:"id,attr,omitempty"`
	MIMEType string   `xml:"mime-type,attr,omitempty"`
	Profile  string   `xml:"profile,attr,omitempty"`
	Any      *Any     `xml:",any,omitempty"`
	File     *File    `xml:",omitempty"`
	Feature  FeatureNegotation
}

// FeatureNegotation is a data element from http://xmpp.org/extensions/xep-0020.html.
type FeatureNegotation struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/feature-neg feature"`
	Form    Form
}
