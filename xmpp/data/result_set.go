package data

import "encoding/xml"

// ResultSet is a data element from http://xmpp.org/extensions/xep-0059.html.
type ResultSet struct {
	XMLName xml.Name    `xml:"http://jabber.org/protocol/rsm set"`
	Max     int         `xml:"max,omitempty"`
	Count   int         `xml:"count,omitempty"`
	Index   int         `xml:"index,omitempty"`
	After   string      `xml:"after,omitempty"`
	Before  string      `xml:"before,omitempty"`
	Last    string      `xml:"last,omitempty"`
	First   ResultFirst `xml:"first,omitempty"`
}

// ResultFirst is a data element from http://xmpp.org/extensions/xep-0059.html.
type ResultFirst struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/rsm first"`
	Index   string   `xml:"index,attr,omitempty"`
	Data    string   `xml:",cdata"`
}
