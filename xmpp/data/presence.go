package data

import "encoding/xml"

// ClientPresence contains XMPP information about a presence update
type ClientPresence struct {
	XMLName xml.Name `xml:"jabber:client presence"`
	From    string   `xml:"from,attr,omitempty"`
	ID      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"` // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
	Lang    string   `xml:"lang,attr,omitempty"`

	Show     string      `xml:"show,omitempty"`   // away, chat, dnd, xa
	Status   string      `xml:"status,omitempty"` // sb []clientText
	Priority string      `xml:"priority,omitempty"`
	Caps     *ClientCaps `xml:"c"`

	Error *StanzaError `xml:"jabber:client error"`
	Delay *Delay       `xml:"delay,omitempty"`

	Chat *struct {
		Item struct {
			JID         string `xml:"jid,attr,omitempty"`
			Affiliation string `xml:"affiliation,attr,omitempty"`
			Role        string `xml:"role,attr,omitempty"`
		} `xml:"item,omitempty"`
		Status struct {
			Code int `xml:"code,attr,omitempty"`
		} `xml:"status,omitempty"`
	} `xml:"http://jabber.org/protocol/muc#user x,omitempty"`

	Extra string `xml:",innerxml"`
}
