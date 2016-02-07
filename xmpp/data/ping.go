package data

import "encoding/xml"

// PingRequest represents a Ping IQ as defined by XEP-0199
type PingRequest struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
