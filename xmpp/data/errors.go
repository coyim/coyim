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
	Error   interface{} `xml:"error"`
}

// ErrorBadRequest reflects a bad-request stanza. See
// http://xmpp.org/rfcs/rfc6120.html#stanzas-error-conditions-bad-request
type ErrorBadRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas bad-request"`
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
			"urn:ietf:params:xml:ns:xmpp-streams", string(c),
		},
	}

	e.EncodeToken(t)
	e.EncodeToken(t.End())

	return nil
}

// Stream error conditions as defined in RFC 6120, section 4.9.3
const (
	BadFormat              StreamErrorCondition = "bad-format"
	BadNamespacePrefix                          = "bad-namespace-prefix"
	Conflict                                    = "conflict"
	ConnectionTimeout                           = "connection-timeout"
	HostGone                                    = "host-gone"
	HostUnknown                                 = "host-unknown"
	ImproperAddressing                          = "improper-addressing"
	InternalServerError                         = "internal-server-error"
	InvalidFrom                                 = "invalid-from"
	InvalidNamespace                            = "invalid-namespace"
	InvalidXML                                  = "invalid-xml"
	NotAuthorized                               = "not-authorized"
	NotWellFormed                               = "not-well-formed"
	PolicyViolation                             = "policy-violation"
	RemoteConnectionFailed                      = "remote-connection-failed"
	Reset                                       = "reset"
	ResourceConstraint                          = "resource-constraint"
	RestrictedXML                               = "restricted-xml"
	SeeOtherHost                                = "see-other-host"
	SystemShutdown                              = "system-shutdown"
	UndefinedCondition                          = "undefined-condition"
	UnsupportedEncoding                         = "unsupported-encoding"
	UnsupportedFeature                          = "unsupported-feature"
	UnsupportedStanzaType                       = "unsupported-stanza-type"
	UnsupportedVersion                          = "unsupported-version"
)
