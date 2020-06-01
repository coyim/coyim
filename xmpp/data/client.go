package data

import "encoding/xml"

//Any provides a convenient way to debug any child element
type Any struct {
	XMLName xml.Name
	Body    string `xml:",innerxml"`
}

// Extensions implements generic XEPs.
type Extensions []*Extension

// Extension represents any XML node not included in the Stanza definition
type Extension Any

// StanzaError implements RFC 3920, section 9.3.
//TODO RFC 6120 obsoletes RFC 3920, section "8.3.2. Syntax"
type StanzaError struct {
	By   string `xml:"by,attr"`
	Code string `xml:"code,attr"` // NOTE: I could not find this in the spec

	// cancel -- do not retry (the error is unrecoverable)
	// continue -- proceed (the condition was only a warning)
	// modify -- retry after changing the data sent
	// auth -- retry after providing credentials
	// wait -- retry after waiting (the error is temporary)
	Type string `xml:"type,attr"`
	Text string `xml:"urn:ietf:params:xml:ns:xmpp-stanzas text"`

	//TODO: Replace references by a function call
	Condition struct {
		XMLName xml.Name
		Body    string `xml:",innerxml"`
	} `xml:",any"`

	ApplicationCondition *Any `xml:",any,ommitempty"`
}

// ClientMessage implements RFC 3921  B.1  jabber:client
type ClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"` // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject    *string     `xml:"subject"`
	Body       string      `xml:"body"`
	Thread     string      `xml:"thread"`
	Delay      *Delay      `xml:"delay,omitempty"`
	Encryption *Encryption `xml:"encryption,omitempty"`

	Error *StanzaError `xml:"error"`

	Extensions `xml:",any,omitempty"`
}

// ClientCaps contains information about client capabilities
type ClientCaps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Ext     string   `xml:"ext,attr"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
}
