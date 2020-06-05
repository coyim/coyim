package data

import (
	"encoding/xml"
	"fmt"
)

// ErrorReply reflects an XMPP error stanza. See
// http://xmpp.org/rfcs/rfc6120.html#stanzas-error-syntax
type ErrorReply struct {
	XMLName xml.Name    `xml:"error"`
	Type    string      `xml:"type,attr"`
	Code    int         `xml:"code,attr,omitempty"`
	Error   interface{} `xml:"error"`
	Error2  interface{} `xml:"error2,omitempty"`
	Text    string      `xml:"urn:ietf:params:xml:ns:xmpp-stanzas text,omitempty"`
}

// ErrorBadRequest reflects a bad-request stanza. See
// http://xmpp.org/rfcs/rfc6120.html#stanzas-error-conditions-bad-request
type ErrorBadRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas bad-request"`
}

// ErrorForbidden reflects a forbidden stanza.
type ErrorForbidden struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas forbidden"`
}

// ErrorNotAcceptable reflects a not acceptable stanza
type ErrorNotAcceptable struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-acceptable"`
}

// ErrorItemNotFound reflects an item not found stanza
type ErrorItemNotFound struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas item-not-found"`
}

// ErrorUnexpectedRequest reflects an unexpected request stanza
type ErrorUnexpectedRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas unexpected-request"`
}

// ErrorNoValidStreams reflects an error when no stream types offered were acceptable
// Ref: XEP-0095
type ErrorNoValidStreams struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si no-valid-streams"`
}

// StreamError represents an XMPP Stream Error as defined in RFC 6120, section 4.9
type StreamError struct {
	XMLName              xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Text                 string   `xml:"text,omitempty"`
	AppSpecificCondition *Any     `xml:",any,omitempty"`

	DefinedCondition StreamErrorCondition
}

func (s *StreamError) String() string {
	if len(s.Text) > 0 {
		return s.Text
	}

	if s.AppSpecificCondition != nil {
		return fmt.Sprintf("%s", s.AppSpecificCondition.XMLName)
	}

	return ""
}

// StreamErrorCondition represents a defined stream error condition
// as defined in RFC 6120, section 4.9.3
type StreamErrorCondition string

// MarshalXML implements xml.Marshaler interface
func (c StreamErrorCondition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Space: "urn:ietf:params:xml:ns:xmpp-streams", Local: string(c),
		},
	}

	_ = e.EncodeToken(t)
	_ = e.EncodeToken(t.End())

	return nil
}

// Stream error conditions as defined in RFC 6120, section 4.9.3
const (
	BadFormat              StreamErrorCondition = "bad-format"
	BadNamespacePrefix     StreamErrorCondition = "bad-namespace-prefix"
	Conflict               StreamErrorCondition = "conflict"
	ConnectionTimeout      StreamErrorCondition = "connection-timeout"
	HostGone               StreamErrorCondition = "host-gone"
	HostUnknown            StreamErrorCondition = "host-unknown"
	ImproperAddressing     StreamErrorCondition = "improper-addressing"
	InternalServerError    StreamErrorCondition = "internal-server-error"
	InvalidFrom            StreamErrorCondition = "invalid-from"
	InvalidNamespace       StreamErrorCondition = "invalid-namespace"
	InvalidXML             StreamErrorCondition = "invalid-xml"
	NotAuthorized          StreamErrorCondition = "not-authorized"
	NotWellFormed          StreamErrorCondition = "not-well-formed"
	PolicyViolation        StreamErrorCondition = "policy-violation"
	RemoteConnectionFailed StreamErrorCondition = "remote-connection-failed"
	Reset                  StreamErrorCondition = "reset"
	ResourceConstraint     StreamErrorCondition = "resource-constraint"
	RestrictedXML          StreamErrorCondition = "restricted-xml"
	SeeOtherHost           StreamErrorCondition = "see-other-host"
	SystemShutdown         StreamErrorCondition = "system-shutdown"
	UndefinedCondition     StreamErrorCondition = "undefined-condition"
	UnsupportedEncoding    StreamErrorCondition = "unsupported-encoding"
	UnsupportedFeature     StreamErrorCondition = "unsupported-feature"
	UnsupportedStanzaType  StreamErrorCondition = "unsupported-stanza-type"
	UnsupportedVersion     StreamErrorCondition = "unsupported-version"
)
