// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"
)

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// Scan XML token stream to finc next Element (start or end)
func nextElement(p *xml.Decoder) (xml.Token, error) {
	for {
		t, err := p.Token()
		if err != nil {
			return xml.StartElement{}, err
		}

		switch elem := t.(type) {
		case xml.StartElement, xml.EndElement:
			return t, nil
		case xml.CharData:
			// https://xmpp.org/rfcs/rfc6120.html#xml-whitespace
			// rfc6120, section 1.4: "whitespace" is used to refer to any character
			// or characters matching [...] SP, HTAB, CR, or LF.
			switch string(elem) {
			case " ", "\t", "\r", "\n": //TODO: consider more than one whitespace
				log.Println("xmpp: received whitespace ping")
			}
		case xml.ProcInst:
			if !(elem.Target == "xml" && strings.HasPrefix(string(elem.Inst), "version=")) {
				log.Printf("xmpp: received unhandled ProcInst element: target=%s inst=%s\n", elem.Target, string(elem.Inst))
			}
		default:
			log.Printf("xmpp: received unhandled element: %#v\n", elem)
		}
	}
}

// Scan XML token stream to find next StartElement.
func nextStart(p *xml.Decoder) (xml.StartElement, error) {
	for {
		t, err := nextElement(p)
		if err != nil {
			return xml.StartElement{}, err
		}

		if start, ok := t.(xml.StartElement); ok {
			return start, nil
		}
	}
}

// Scan XML token stream for next element and save into val.
// If val == nil, allocate new element based on proto map.
// Either way, return val.
func next(c interfaces.Conn) (xml.Name, interface{}, error) {
	elem, err := nextElement(c.In())
	if err != nil {
		return xml.Name{}, nil, err
	}

	c.Lock().Lock()
	defer c.Lock().Unlock()

	switch el := elem.(type) {
	case xml.StartElement:
		return decodeStartElement(c, el)
	case xml.EndElement:
		return decodeEndElement(el)
	}

	return xml.Name{}, nil, fmt.Errorf("unexpected element %s", elem)
}

func decodeStartElement(c interfaces.Conn, se xml.StartElement) (xml.Name, interface{}, error) {

	// Put it in an interface and allocate one.
	var nv interface{}
	if t, e := c.CustomStorage()[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else if t, e := defaultStorage[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else {
		return xml.Name{}, nil, errors.New("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	// Unmarshal into that storage.
	if err := c.In().DecodeElement(nv, &se); err != nil {
		return xml.Name{}, nil, err
	}

	return se.Name, nv, nil
}

func decodeEndElement(ce xml.EndElement) (xml.Name, interface{}, error) {
	switch ce.Name {
	case xml.Name{NsStream, "stream"}:
		return ce.Name, &data.StreamClose{}, nil
	}

	return ce.Name, nil, nil
}

var defaultStorage = map[xml.Name]reflect.Type{
	xml.Name{Space: NsStream, Local: "features"}: reflect.TypeOf(data.StreamFeatures{}),
	xml.Name{Space: NsStream, Local: "error"}:    reflect.TypeOf(data.StreamError{}),
	xml.Name{Space: NsTLS, Local: "starttls"}:    reflect.TypeOf(data.StartTLS{}),
	xml.Name{Space: NsTLS, Local: "proceed"}:     reflect.TypeOf(data.ProceedTLS{}),
	xml.Name{Space: NsTLS, Local: "failure"}:     reflect.TypeOf(data.FailureTLS{}),
	xml.Name{Space: NsSASL, Local: "mechanisms"}: reflect.TypeOf(data.SaslMechanisms{}),
	xml.Name{Space: NsSASL, Local: "challenge"}:  reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "response"}:   reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "abort"}:      reflect.TypeOf(data.SaslAbort{}),
	xml.Name{Space: NsSASL, Local: "success"}:    reflect.TypeOf(data.SaslSuccess{}),
	xml.Name{Space: NsSASL, Local: "failure"}:    reflect.TypeOf(data.SaslFailure{}),
	xml.Name{Space: NsBind, Local: "bind"}:       reflect.TypeOf(data.BindBind{}),
	xml.Name{Space: NsClient, Local: "message"}:  reflect.TypeOf(data.ClientMessage{}),
	xml.Name{Space: NsClient, Local: "presence"}: reflect.TypeOf(data.ClientPresence{}),
	xml.Name{Space: NsClient, Local: "iq"}:       reflect.TypeOf(data.ClientIQ{}),
	xml.Name{Space: NsClient, Local: "error"}:    reflect.TypeOf(data.ClientError{}),
}
