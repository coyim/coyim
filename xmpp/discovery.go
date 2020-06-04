// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"sort"

	"github.com/coyim/coyim/xmpp/data"
)

// HasSupportTo uses XEP-0030 to checks if an entity supports a feature.
// The entity is identified by its JID and the feature by its XML namespace.
// It returns true if the feaure is reported to be supported and false
// otherwise (including if any error happened).
func (c *conn) HasSupportTo(entity string, feature string) bool {
	if res, ok := c.DiscoveryFeatures(entity); ok {
		return stringArrayContains(res, feature)
	}

	return false
}

func stringArrayContains(r []string, a string) bool {
	for _, f := range r {
		if f == a {
			return true
		}
	}

	return false
}

func (c *conn) sendDiscoveryInfo(to string) (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	return c.SendIQ(to, "get", &data.DiscoveryInfoQuery{})
}

func parseDiscoveryInfoReply(iq *data.ClientIQ) (*data.DiscoveryInfoQuery, error) {
	reply := &data.DiscoveryInfoQuery{}
	err := xml.Unmarshal(iq.Query, reply)
	return reply, err
}

func (c *conn) sendDiscoveryItems(to string) (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	return c.SendIQ(to, "get", &data.DiscoveryItemsQuery{})
}

func parseDiscoveryItemsReply(iq *data.ClientIQ) (*data.DiscoveryItemsQuery, error) {
	reply := &data.DiscoveryItemsQuery{}
	err := xml.Unmarshal(iq.Query, reply)
	return reply, err
}

// TODO: at some point we need to cache these features somewhere

// QueryServiceInformation sends a service discovery information ("disco#info") query.
// See XEP-0030, Section "3. Discovering Information About a Jabber Entity"
// This method blocks until conn#Next() receives the response to the IQ.
func (c *conn) QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error) {
	reply, _, err := c.sendDiscoveryInfo(entity)
	if err != nil {
		return nil, err
	}

	stanza, ok := <-reply
	if !ok {
		return nil, errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: failed to parse response")
	}

	return parseDiscoveryInfoReply(iq)
}

// QueryServiceItems sends a Service Discovery items ("disco#items") query.
// See XEP-0030, Section "4. Discovering the Items Associated with a Jabber Entity"
// This method blocks until conn#Next() receives the response to the IQ.
func (c *conn) QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error) {
	reply, _, err := c.sendDiscoveryItems(entity)
	if err != nil {
		return nil, err
	}

	stanza, ok := <-reply
	if !ok {
		return nil, errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: failed to parse response")
	}

	return parseDiscoveryItemsReply(iq)
}

func (c *conn) DiscoveryFeatures(entity string) ([]string, bool) {
	discoveryReply, err := c.QueryServiceInformation(entity)
	if err != nil {
		return nil, false
	}

	var result []string
	for _, f := range discoveryReply.Features {
		result = append(result, f.Var)
	}

	return result, true
}

func (c *conn) DiscoveryFeaturesAndIdentities(entity string) ([]data.DiscoveryIdentity, []string, bool) {
	discoveryReply, err := c.QueryServiceInformation(entity)
	if err != nil {
		return nil, nil, false
	}

	var result []string
	for _, f := range discoveryReply.Features {
		result = append(result, f.Var)
	}

	return discoveryReply.Identities, result, true
}

//DiscoveryReply returns a minimum reply to a http://jabber.org/protocol/disco#info query
func DiscoveryReply(name string) data.DiscoveryInfoQuery {
	return data.DiscoveryInfoQuery{
		Identities: []data.DiscoveryIdentity{
			{
				Category: "client",
				Type:     "pc",

				//NOTE: this is optional as per XEP-0030
				Name: name,
			},
		},
		//TODO: extract constants that document which XEPs are supported
		Features: []data.DiscoveryFeature{
			{Var: "http://jabber.org/protocol/disco#info"},                    //XEP-0030
			{Var: "urn:xmpp:bob"},                                             //XEP-0231
			{Var: "urn:xmpp:ping"},                                            //XEP-0199
			{Var: "http://jabber.org/protocol/caps"},                          //XEP-0115
			{Var: "jabber:iq:version"},                                        //XEP-0092
			{Var: "vcard-temp"},                                               //XEP-0054
			{Var: "jabber:x:data"},                                            //XEP-004
			{Var: "http://jabber.org/protocol/si"},                            //XEP-0096
			{Var: "http://jabber.org/protocol/si/profile/file-transfer"},      //XEP-0096
			{Var: "http://jabber.org/protocol/si/profile/directory-transfer"}, //XEP-xxxx: SI Directory Transfer
			//			{Var: "http://jabber.org/protocol/si/profile/encrypted-data-transfer"}, //XEP-xxxx: SI Encrypted Data Transfer
			{Var: "http://jabber.org/protocol/bytestreams"}, //XEP-0047
			{Var: "urn:xmpp:eme:0"},                         //XEP-0380
		},
	}
}

// VerificationString returns a SHA-1 verification string as defined in XEP-0115.
// See http://xmpp.org/extensions/xep-0115.html#ver
func VerificationString(r *data.DiscoveryInfoQuery) (string, error) {
	h := sha1.New()

	seen := make(map[string]bool)
	identitySorter := &xep0115Sorter{}
	for i := range r.Identities {
		identitySorter.add(&r.Identities[i])
	}
	sort.Sort(identitySorter)
	for _, id := range identitySorter.s {
		id := id.(*data.DiscoveryIdentity)
		c := id.Category + "/" + id.Type + "/" + id.Lang + "/" + id.Name + "<"
		if seen[c] {
			return "", errors.New("duplicate discovery identity")
		}
		seen[c] = true
		_, _ = io.WriteString(h, c)
	}

	seen = make(map[string]bool)
	featureSorter := &xep0115Sorter{}
	for i := range r.Features {
		featureSorter.add(&r.Features[i])
	}
	sort.Sort(featureSorter)
	for _, f := range featureSorter.s {
		f := f.(*data.DiscoveryFeature)
		if seen[f.Var] {
			return "", errors.New("duplicate discovery feature")
		}
		seen[f.Var] = true
		_, _ = io.WriteString(h, f.Var+"<")
	}

	seen = make(map[string]bool)
	for _, f := range r.Forms {
		if len(f.Fields) == 0 {
			continue
		}
		fieldSorter := &xep0115Sorter{}
		for i := range f.Fields {
			fieldSorter.add(&f.Fields[i])
		}
		sort.Sort(fieldSorter)
		formTypeField := fieldSorter.s[0].(*data.FormFieldX)
		if formTypeField.Var != "FORM_TYPE" {
			continue
		}
		if seen[formTypeField.Type] {
			return "", errors.New("multiple forms of the same type")
		}
		seen[formTypeField.Type] = true
		if len(formTypeField.Values) != 1 {
			return "", errors.New("form does not have a single FORM_TYPE value")
		}
		if formTypeField.Type != "hidden" {
			continue
		}
		_, _ = io.WriteString(h, formTypeField.Values[0]+"<")
		for _, field := range fieldSorter.s[1:] {
			field := field.(*data.FormFieldX)
			_, _ = io.WriteString(h, field.Var+"<")
			values := append([]string{}, field.Values...)
			sort.Strings(values)
			for _, v := range values {
				_, _ = io.WriteString(h, v+"<")
			}
		}
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
