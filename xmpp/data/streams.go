package data

import "encoding/xml"

// StreamFeatures contain information about the features this stream supports
// RFC 3920  C.1  Streams name space
//TODO RFC 6120 obsoletes RFC 3920
type StreamFeatures struct {
	XMLName            xml.Name `xml:"http://etherx.jabber.org/streams features"`
	StartTLS           StartTLS
	Mechanisms         SaslMechanisms
	Bind               BindBind
	InBandRegistration *InBandRegistration

	// This is a hack for now to get around the fact that the new encoding/xml
	// doesn't unmarshal to XMLName elements.
	Session *string `xml:"session"`

	//TODO: Support additional features, like
	//https://xmpp.org/extensions/xep-0115.html
	//Roster versioning: rfc6121 section 2.6
	//and the features described here
	//https://xmpp.org/registrar/stream-features.html
	//	any []Any `xml:",any,omitempty"`
}

// StreamClose represents a request to close the stream
type StreamClose struct{}
