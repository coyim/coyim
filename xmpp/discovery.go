// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"sort"

	"github.com/twstrike/coyim/xmpp/data"
)

// VerificationString returns a SHA-1 verification string as defined in XEP-0115.
// See http://xmpp.org/extensions/xep-0115.html#ver
func VerificationString(r *data.DiscoveryReply) (string, error) {
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
		io.WriteString(h, c)
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
		io.WriteString(h, f.Var+"<")
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
		io.WriteString(h, formTypeField.Values[0]+"<")
		for _, field := range fieldSorter.s[1:] {
			field := field.(*data.FormFieldX)
			io.WriteString(h, field.Var+"<")
			values := append([]string{}, field.Values...)
			sort.Strings(values)
			for _, v := range values {
				io.WriteString(h, v+"<")
			}
		}
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
