package data

import "encoding/xml"

// Stanza represents a message from the XMPP server.
type Stanza struct {
	Name  xml.Name
	Value interface{}
}
