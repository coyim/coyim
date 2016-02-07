package data

import "encoding/xml"

// BobData is a data element from http://xmpp.org/extensions/xep-0231.html.
type BobData struct {
	XMLName  xml.Name `xml:"urn:xmpp:bob data"`
	CID      string   `xml:"cid,attr"`
	MIMEType string   `xml:"type,attr"`
	Base64   string   `xml:",chardata"`
}
