package data

import "encoding/xml"

// Encryption represents XEP-0380: Explicit Message Encryption
type Encryption struct {
	XMLName   xml.Name `xml:"urn:xmpp:eme:0 encryption"`
	Namespace string   `xml:"namespace,attr"`
	Name      string   `xml:"name,attr,omitempty"`
}
