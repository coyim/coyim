package data

import "encoding/xml"

// StartTLS represents a TLS start
// RFC 6120, section 5
type StartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required xml.Name `xml:"required"`
}

// ProceedTLS represents a TLS proceed
type ProceedTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

// FailureTLS represents a TLS failure
type FailureTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}
