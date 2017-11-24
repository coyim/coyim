package data

import "encoding/xml"

// Hash is a data element from http://xmpp.org/extensions/xep-0300.html.
type Hash struct {
	XMLName   xml.Name `xml:"urn:xmpp:hashes:2 hash"`
	Algorithm string   `xml:"algo,attr"`
	Base64    string   `xml:",chardata"`
}

// HashUsed is a data element from http://xmpp.org/extensions/xep-0300.html.
type HashUsed struct {
	XMLName   xml.Name `xml:"urn:xmpp:hashes:2 hash-used"`
	Algorithm string   `xml:"algo,attr"`
}
