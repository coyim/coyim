package data

import "encoding/xml"

// BindBind represents a bind
// RFC 6120, section 7
type BindBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource"`
	Jid      string   `xml:"jid"`
}
